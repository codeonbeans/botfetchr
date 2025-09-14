package tgbot

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/codeonbeans/botfetchr/config"
	"github.com/codeonbeans/botfetchr/internal/client/browserpool"
	"github.com/codeonbeans/botfetchr/internal/client/mediasaver/instagram"
	"github.com/codeonbeans/botfetchr/internal/client/mediasaver/vk"
	"github.com/codeonbeans/botfetchr/internal/logger"
	"github.com/codeonbeans/botfetchr/internal/storage"

	"github.com/corpix/uarand"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/go-rod/rod"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"golang.org/x/net/proxy"
)

type MediaSaver interface {
	GetUA() string
	GetFilename(ogUrl, directUrl string) string
	GetVideoURLs(ctx context.Context, browser *rod.Browser, url string) ([]string, error)
	IsValidURL(url string) bool

	SetUserAgent(ua string)
	SetQuality(quality string)
	SetTimeout(timeout time.Duration)
}

type MediaSaverFactory func() (MediaSaver, error)

var mediaSaverFactory = map[SaverType]MediaSaverFactory{
	SaverTypeInstagram: func() (MediaSaver, error) {
		client := instagram.NewClient()
		configMediaSaver(client)
		return client, nil
	},
	SaverTypeVK: func() (MediaSaver, error) {
		client := vk.NewClient()
		configMediaSaver(client)
		return client, nil
	},
}

type DefaultBot struct {
	*bot.Bot
	subscriptionMux sync.Mutex
	storage         *storage.Storage
	cacheManager    *marshaler.Marshaler
	browserPool     browserpool.Client
}

func New(store *storage.Storage, cacheManager *marshaler.Marshaler) (*DefaultBot, error) {

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

		var auth *proxy.Auth
		if username != "" && password != "" {
			auth = &proxy.Auth{
				User:     username,
				Password: password,
			}
		} else {
			logger.Log.Sugar().Warn("Proxy username or password not set, using anonymous proxy")
		}

		dialer, err := proxy.SOCKS5("tcp", address, auth, proxy.Direct)
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
	// Assign cache manager
	defaultBot.cacheManager = cacheManager

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
	defaultBot.Bot, err = bot.New(config.GetConfig().TelegramBot.Token, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot client: %w", err)
	}

	// Assign browser pool
	defaultBot.browserPool, err = browserpool.NewClient(browserpool.Config{
		Headless:      config.GetConfig().BrowserPool.Headless,
		Proxies:       config.GetConfig().BrowserPool.Proxies,
		PoolSize:      config.GetConfig().BrowserPool.PoolSize,
		TaskQueueSize: config.GetConfig().BrowserPool.TaskQueueSize,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create browser pool: %w", err)
	}

	return defaultBot, nil
}

func (b *DefaultBot) GetMediaSaver(url string) (MediaSaver, error) {
	for saverType, factory := range mediaSaverFactory {
		client, err := factory()
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
	b.Bot.Start(ctx)
}

func getUA() string {
	if config.GetConfig().MediaSaver.UseRandomUA {
		return uarand.GetRandom()
	}

	uas := config.GetConfig().MediaSaver.UserAgents
	if len(uas) == 0 {
		return uarand.GetRandom()
	}

	// Return a random user agent from the list
	randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(uas))))
	if err != nil {
		return uas[0]
	}

	return uas[randomIndex.Int64()]
}

func configMediaSaver(mediaSaver MediaSaver) {
	if mediaSaver == nil {
		return
	}

	mediaSaver.SetUserAgent(getUA())
	mediaSaver.SetQuality(config.GetConfig().MediaSaver.Quality)
	mediaSaver.SetTimeout(time.Duration(config.GetConfig().MediaSaver.Timeout) * time.Second)
}
