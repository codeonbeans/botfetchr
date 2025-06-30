package tgbot

import (
	"botmediasaver/generated/sqlc"
	"botmediasaver/internal/logger"
	"botmediasaver/internal/model"
	"botmediasaver/internal/utils/ptr"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/jackc/pgx/v5/pgtype"
)

type MediaResult struct {
	State  string
	Medias []MediaData
}

type MediaData struct {
	Filename  string
	Size      int64
	Media     io.ReadCloser
	DirectURL string
}

// ProcessingContext encapsulates all the context needed for processing a URL
type ProcessingContext struct {
	ctx           context.Context
	chatID        int64
	originalMsgID int
	urlIndex      int
	url           string
	statusMsg     *models.Message
}

func (b *DefaultBot) Handler(ctx context.Context, update *models.Update) {
	account, err := b.storage.GetAccountTelegram(ctx, sqlc.GetAccountTelegramParams{
		TelegramID: pgtype.Int8{Int64: int64(update.Message.From.ID), Valid: true},
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			account, err = b.storage.CreateAccountTelegram(ctx, sqlc.CreateAccountTelegramParams{
				TelegramID:   update.Message.From.ID,
				IsBot:        update.Message.From.IsBot,
				FirstName:    update.Message.From.FirstName,
				LastName:     update.Message.From.LastName,
				Username:     pgtype.Text{String: update.Message.From.Username, Valid: true},
				LanguageCode: update.Message.From.LanguageCode,
				IsPremium:    update.Message.From.IsPremium,
			})
			if err != nil {
				logger.Log.Sugar().Errorf("Failed to create account: %v", err)
				return
			}
		} else {
			logger.Log.Sugar().Errorf("Failed to get account: %v", err)
			return
		}
	}

	urls := extractURLs(update.Message.Text)
	for i, url := range urls {
		go b.processURLAsync(ctx, account, update, url, i)
	}
}

func (b *DefaultBot) processURLAsync(ctx context.Context, account sqlc.AccountTelegram, update *models.Update, url string, index int) {
	// Check subscription
	if err := b.IsAccountAllow(ctx, account.ID, model.FeatureGetMedia, 1); err != nil {
		logger.Log.Sugar().Errorf("Account %d is not allowed to download: %v", update.Message.From.ID, err)
		// Send error message to user
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("You are not allowed to download media: %v", err),
			ReplyParameters: &models.ReplyParameters{
				MessageID: update.Message.ID,
			},
		}); err != nil {
			logger.Log.Sugar().Errorf("Failed to send error message: %v", err)
			return
		}

		return
	}

	processCtx := &ProcessingContext{
		ctx:           ctx,
		chatID:        update.Message.Chat.ID,
		originalMsgID: update.Message.ID,
		urlIndex:      index,
		url:           url,
	}

	// Send initial status message
	statusMsg, err := b.sendInitialStatus(ctx, processCtx)
	if err != nil {
		logger.Log.Sugar().Errorf("Failed to send initial status: %v", err)
		return
	}
	processCtx.statusMsg = statusMsg

	processor := &MediaProcessor{
		bot:        b,
		processCtx: processCtx,
		updateChan: make(chan MediaResult, 10),
	}

	// Start status updater goroutine
	go processor.handleStatusUpdates()

	// Process the URL
	if err := processor.processURL(); err != nil {
		processor.updateChan <- MediaResult{State: err.Error()}
	}

	// Clean up
	time.Sleep(30 * time.Second)
	close(processor.updateChan)
}

func (b *DefaultBot) sendInitialStatus(ctx context.Context, processCtx *ProcessingContext) (*models.Message, error) {
	return b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: processCtx.chatID,
		Text:   fmt.Sprintf("%d. %s\nState: ⌛ queued...", processCtx.urlIndex+1, processCtx.url),
		ReplyParameters: &models.ReplyParameters{
			MessageID: processCtx.originalMsgID,
		},
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: ptr.ToPtr(true),
		},
	})
}

func extractURLs(text string) []string {
	lines := strings.Split(text, "\n")
	var urls []string

	for _, line := range lines {
		if strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") {
			urls = append(urls, line)
		}
	}

	return urls
}
