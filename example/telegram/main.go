package main

import (
	"botvideosaver/config"
	"botvideosaver/internal/client/instagram"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"golang.org/x/net/proxy"
)

// Send any text message to the bot after the bot has been started

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	address := fmt.Sprintf(
		"%s:%d",
		config.GetConfig().TelegramBot.Proxy.Address,
		config.GetConfig().TelegramBot.Proxy.Port,
	)
	username := config.GetConfig().TelegramBot.Proxy.Username
	password := config.GetConfig().TelegramBot.Proxy.Password

	dialer, err := proxy.SOCKS5("tcp", address, &proxy.Auth{
		User:     username,
		Password: password,
	}, proxy.Direct)
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
		},
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
		bot.WithHTTPClient(time.Minute, client),
	}

	b, err := bot.New(config.GetConfig().TelegramBot.Token, opts...)
	if err != nil {
		panic(err)
	}

	fmt.Println("Bot started!")

	b.Start(ctx)

}

var insta = instagram.NewClient()

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println("Received update:", update)

	urls := strings.Split(update.Message.Text, "\n")
	for _, instaUrl := range urls {
		go func() {
			fmt.Println("Processing URL:", instaUrl)
			url, err := insta.GetVideoURL(instaUrl)
			if err != nil {
				fmt.Printf("Failed to get video URL: %v\n", err)
			}

			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("Failed to download video: %v\n", err)
				return
			}
			defer resp.Body.Close()

			videoID, err := insta.GetVideoID(instaUrl)
			if err != nil {
				fmt.Printf("Failed to get video ID: %v\n", err)
			}

			inputVideo := &models.InputMediaVideo{
				Media:           fmt.Sprintf("attach://%s.mp4", videoID),
				MediaAttachment: resp.Body,
			}

			if _, err := b.SendMediaGroup(ctx, &bot.SendMediaGroupParams{
				ChatID: update.Message.Chat.ID,
				Media:  []models.InputMedia{inputVideo},
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
			}); err != nil {
				fmt.Printf("Failed to send video: %v\n", err)
			}
		}()
	}
}
