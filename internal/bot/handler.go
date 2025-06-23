package tgbot

import (
	"botvideosaver/internal/client/browserpool"
	"botvideosaver/internal/logger"
	"botvideosaver/internal/utils/ptr"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/corpix/uarand"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *DefaultBot) Handler(ctx context.Context, update *models.Update) error {
	lines := strings.Split(update.Message.Text, "\n")

	for i, url := range lines {
		go func(url string) {
			// Send initial status message for this URL
			statusMsg, _ := b.bot.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("%d. [%s]\nState: ⌛ queued...", i+1, url),
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
				LinkPreviewOptions: &models.LinkPreviewOptions{
					IsDisabled: ptr.ToPtr(true),
				},
			})

			// Create update channel for editting message (via VideoResult)
			// Remember to send Result with final media to completely close the channel
			updateMessageChan := make(chan VideoResult)
			go func() {
				for result := range updateMessageChan {
					if result.Media != nil {
						inputVideo := &models.InputMediaVideo{
							Media:           fmt.Sprintf("attach://%s.mp4", result.ID),
							MediaAttachment: result.Media,
							Caption:         url,
						}

						_, err := b.bot.SendMediaGroup(ctx, &bot.SendMediaGroupParams{
							ChatID: update.Message.Chat.ID,
							Media:  []models.InputMedia{inputVideo},
							ReplyParameters: &models.ReplyParameters{
								MessageID: update.Message.ID,
							},
						})
						if err != nil {
							fmt.Println(err)
						}

						_, err = b.bot.DeleteMessage(ctx, &bot.DeleteMessageParams{
							ChatID:    update.Message.Chat.ID,
							MessageID: statusMsg.ID,
						})
						if err != nil {
							fmt.Println(err)
						}

						result.Media.Close()
						close(updateMessageChan)
					} else {
						b.bot.EditMessageText(ctx, &bot.EditMessageTextParams{
							ChatID:    update.Message.Chat.ID,
							MessageID: statusMsg.ID,
							Text:      fmt.Sprintf("%d. [%s]\nState: %s", i+1, url, result.State),
						})
					}
				}
			}()

			if err := b.handleURL(url, updateMessageChan); err != nil {
				updateMessageChan <- VideoResult{
					State: fmt.Sprintf("❌ failed: %s", err.Error()),
				}
			}
		}(url)
	}

	return nil
}

func (b *DefaultBot) handleURL(url string, updateMessageChan chan VideoResult) error {
	var (
		directUrl string
		videoID   string
	)

	updateMessageChan <- VideoResult{
		State: "🔎 getting video info...",
	}

	if err := b.browserPool.UseBrowser(func(browser *browserpool.Browser) error {
		logger.Log.Sugar().Infof("Processing URL: %s", url)

		ua := uarand.GetRandom()

		saver, err := b.GetVideoSaver(url, ua, browser.Browser)
		if err != nil {
			return fmt.Errorf("failed to get video saver: %w", err)
		}

		videoID, err = saver.GetVideoID(url)
		if err != nil {
			return fmt.Errorf("failed to get video ID for %s: %w", url, err)
		}

		directUrl, err = saver.GetVideoURL(url)
		if err != nil {
			return fmt.Errorf("failed to get video URL for %s: %w", url, err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to process URL %s: %w", url, err)
	}

	updateMessageChan <- VideoResult{
		State: "⬇️ downloading video...",
	}

	req, err := http.NewRequest("GET", directUrl, nil)
	if err != nil {
		return fmt.Errorf("failed to download video from %s: %w", url, err)
	}

	// Mimic a real browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download video from %s: %w", url, err)
	}

	updateMessageChan <- VideoResult{
		ID:    videoID,
		Media: resp.Body,
		State: "➤ sending video...",
	}

	return nil
}

type VideoResult struct {
	ID    string
	State string
	Media io.ReadCloser
}
