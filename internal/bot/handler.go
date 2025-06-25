package tgbot

import (
	"botvideosaver/config"
	"botvideosaver/internal/client/browserpool"
	"botvideosaver/internal/logger"
	"botvideosaver/internal/utils/common"
	"botvideosaver/internal/utils/ptr"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type VideoResult struct {
	ID    string
	State string
	Media io.ReadCloser
}

func (b *DefaultBot) Handler(ctx context.Context, update *models.Update) {
	lines := strings.Split(update.Message.Text, "\n")

	for i, url := range lines {
		go func(url string) {
			// Send initial status message for this URL
			statusMsg, _ := b.bot.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("%d. %s\nState: ⌛ queued...", i+1, url),
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
				LinkPreviewOptions: &models.LinkPreviewOptions{
					IsDisabled: ptr.ToPtr(true),
				},
			})

			// Create update channel for editting message (via VideoResult)
			// Always remember to close the channel whether it's successful or not
			updateMessageChan := make(chan VideoResult)
			defer func() {
				time.Sleep(30 * time.Second) // Give some time for final messages to be processed, for now it is for error handling of handleURL
				close(updateMessageChan)
			}()

			// Create a goroutine to handle updates in channel updateMessageChan
			go func() {
				for result := range updateMessageChan {
					// Complete the status message
					if result.Media != nil {
						inputVideo := &models.InputMediaVideo{
							Media:           fmt.Sprintf("attach://%s.mp4", result.ID),
							MediaAttachment: result.Media,
							Caption:         fmt.Sprintf("%d. %s\nState: %s", i+1, url, result.State),
						}

						_, err := b.bot.SendMediaGroup(ctx, &bot.SendMediaGroupParams{
							ChatID: update.Message.Chat.ID,
							Media:  []models.InputMedia{inputVideo},
							ReplyParameters: &models.ReplyParameters{
								MessageID: update.Message.ID,
							},
						})
						if err != nil {
							logger.Log.Sugar().Errorf("Failed to send media group: %v", err)

							state := fmt.Sprintf("❌ failed to send video: %v", err)
							if _, err = b.bot.EditMessageText(ctx, &bot.EditMessageTextParams{
								ChatID:    update.Message.Chat.ID,
								MessageID: statusMsg.ID,
								Text:      fmt.Sprintf("%d. %s\nState: %s", i+1, url, state),
								LinkPreviewOptions: &models.LinkPreviewOptions{
									IsDisabled: ptr.ToPtr(true),
								},
							}); err != nil {
								logger.Log.Sugar().Errorf("Failed to edit message text: %v", err)
							}
						} else {
							// Successfully sent video, now delete the status message
							_, err = b.bot.DeleteMessage(ctx, &bot.DeleteMessageParams{
								ChatID:    update.Message.Chat.ID,
								MessageID: statusMsg.ID,
							})
							if err != nil {
								logger.Log.Sugar().Errorf("Failed to delete status message: %v", err)
							}
						}

						result.Media.Close() // Close the media stream (request body)
					} else {
						// Update status message with the current state
						_, err := b.bot.EditMessageText(ctx, &bot.EditMessageTextParams{
							ChatID:    update.Message.Chat.ID,
							MessageID: statusMsg.ID,
							Text:      fmt.Sprintf("%d. %s\nState: %s", i+1, url, result.State),
							LinkPreviewOptions: &models.LinkPreviewOptions{
								IsDisabled: ptr.ToPtr(true),
							},
						})
						if err != nil {
							logger.Log.Sugar().Errorf("Failed to edit message text: %v", err)
						}
					}
				}
			}()

			if err := b.handleURL(url, updateMessageChan); err != nil {
				updateMessageChan <- VideoResult{
					State: err.Error(),
				}
			}
		}(url)
	}
}

func (b *DefaultBot) handleURL(url string, updateMessageChan chan VideoResult) error {
	attempts := config.GetConfig().VideoSaver.RetryCount

	if err := common.DoWithRetry(common.RetryConfig{
		Attempts: attempts,
		Delay:    2 * time.Second,
	}, func() error {
		var (
			directUrl string
			videoID   string
		)

		if err := b.browserPool.UseBrowser(func(ctx context.Context, browser *browserpool.Browser) error {
			updateMessageChan <- VideoResult{
				State: "🔎 getting video info...",
			}
			logger.Log.Sugar().Infof("Processing URL: %s", url)

			saver, err := b.GetVideoSaver(url, browser.Browser)
			if err != nil {
				return fmt.Errorf("failed to get video saver: %w", err)
			}

			videoID, err = saver.GetVideoID(url)
			if err != nil {
				return fmt.Errorf("failed to get video ID: %w", err)
			}

			directUrl, err = saver.GetVideoURL(ctx, url)
			if err != nil {
				return fmt.Errorf("failed to get direct video URL: %w", err)
			}

			return nil
		}); err != nil {
			return err
		}

		sizeMB, err := getFileSizeMB(directUrl)
		sizeMBStr := fmt.Sprintf("%.2f MB", sizeMB)
		if err != nil {
			sizeMBStr = "unknown size"
		}

		updateMessageChan <- VideoResult{
			// State: "⬇️ downloading video...",
			State: fmt.Sprintf("⬇️ downloading video... (%s)\nOr you can download it manually: %s", sizeMBStr, directUrl),
		}

		req, err := http.NewRequest("GET", directUrl, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		// Mimic a real browser
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("failed to download video (%s): %w", sizeMBStr, err)
		}

		updateMessageChan <- VideoResult{
			State: fmt.Sprintf("➤ video downloaded (%s), sending video...\nOr you can download it manually: %s", sizeMBStr, directUrl),
		}

		updateMessageChan <- VideoResult{
			ID:    videoID,
			Media: resp.Body,
			State: fmt.Sprintf("✅ get video successfully (%s)", sizeMBStr),
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to get video info after %d attempts: %w", attempts, err)
	}

	return nil
}

func getFileSizeMB(url string) (float64, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	contentLength := resp.Header.Get("Content-Length")
	if contentLength == "" {
		return 0, fmt.Errorf("Content-Length header not found")
	}

	bytes, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		return 0, err
	}

	// Convert bytes to MB
	mb := float64(bytes) / (1024 * 1024)
	return mb, nil
}
