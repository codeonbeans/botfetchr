package tgbot

import (
	"botvideosaver/config"
	"botvideosaver/internal/client/browserpool"
	"botvideosaver/internal/client/instagram"
	"botvideosaver/internal/client/vk"
	"botvideosaver/internal/logger"
	"botvideosaver/internal/storage"
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"golang.org/x/net/proxy"
)

type VideoSaver interface {
	IsValidURL(url string) bool
	GetVideoURL(url string) (string, error)
	GetVideoID(url string) (string, error)
}

type VideoSaverFactory func(ua string, browser *rod.Browser) (VideoSaver, error)

var videoSaverFactory = map[SaverType]VideoSaverFactory{
	SaverTypeInstagram: func(ua string, browser *rod.Browser) (VideoSaver, error) {
		return instagram.NewClient(ua, browser)
	},
	SaverTypeVK: func(ua string, browser *rod.Browser) (VideoSaver, error) {
		return vk.NewClient(ua, browser)
	},
}

type DefaultBot struct {
	bot         *bot.Bot
	storage     *storage.Storage
	browserPool browserpool.Client
}

func New(store *storage.Storage) (*DefaultBot, error) {
	logger.Log.Sugar().Info("Initializing bot...")

	var httpClient *http.Client

	// Check if proxy is enabled in the configuration
	if config.GetConfig().TelegramBot.Proxy.Enabled {
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
			return nil, fmt.Errorf("failed to create SOCKS5 dialer: %w", err)
		}

		httpClient = &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return dialer.Dial(network, addr)
				},
			},
		}
	} else {
		httpClient = http.DefaultClient
	}

	var (
		err        error
		defaultBot = &DefaultBot{}
	)

	// Assign storage
	defaultBot.storage = store

	opts := []bot.Option{
		bot.WithDefaultHandler(func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			// add recover panic here
			defer func() {
				if r := recover(); r != nil {
					logger.Log.Sugar().Errorf("Recovered from panic: %v", r)
					debug.PrintStack()
				}
			}()

			defaultBot.Handler(ctx, update)
		}),
		bot.WithHTTPClient(time.Minute, httpClient),
		bot.WithDebugHandler(logger.Log.Sugar().Debugf),
	}

	if config.GetConfig().TelegramBot.LogDebug {
		opts = append(opts, bot.WithDebug())
	}

	// Assign bot client
	defaultBot.bot, err = bot.New(config.GetConfig().TelegramBot.Token, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot client: %w", err)
	}

	// Assign browser pool
	defaultBot.browserPool, err = browserpool.NewClient(browserpool.Config{
		Proxies:       config.GetConfig().BrowserPool.Proxies,
		PoolSize:      config.GetConfig().BrowserPool.PoolSize,
		TaskQueueSize: config.GetConfig().BrowserPool.TaskQueueSize,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create browser pool: %w", err)
	}

	return defaultBot, nil
}

func (b *DefaultBot) GetVideoSaver(url string, ua string, browser *rod.Browser) (VideoSaver, error) {
	for saverType, factory := range videoSaverFactory {
		client, err := factory(ua, browser)
		if err != nil {
			return nil, fmt.Errorf("failed to create client for type %s: %w", saverType, err)
		}

		if client.IsValidURL(url) {
			return client, nil
		}
	}

	return nil, fmt.Errorf("no valid video saver found for URL: %s", url)
}

func (b *DefaultBot) Start(ctx context.Context) {
	logger.Log.Sugar().Info("Starting Telegram bot...")
	b.bot.Start(ctx)
}
