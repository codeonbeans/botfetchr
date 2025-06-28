package tgbot

import (
	"botvideosaver/config"
	"botvideosaver/internal/client/browserpool"
	"botvideosaver/internal/logger"
	"botvideosaver/internal/utils/common"
	"botvideosaver/internal/utils/download"
	"botvideosaver/internal/utils/ptr"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type MediaResult struct {
	State  string
	Medias []MediaData
}

type MediaData struct {
	Filename string
	Media    io.ReadCloser
}

func (b *DefaultBot) Handler(ctx context.Context, update *models.Update) {
	lines := strings.Split(update.Message.Text, "\n")

	for i, url := range lines {
		go func(url string) {
			isUrl := strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
			if !isUrl {
				return // Skip non-URL lines
			}

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

			// Create update channel for editting message (via MediaResult)
			// Always remember to close the channel whether it's successful or not
			updateMessageChan := make(chan MediaResult)
			defer func() {
				time.Sleep(30 * time.Second) // Give some time for final messages to be processed, for now it is for error handling of handleURL
				close(updateMessageChan)
			}()

			// Create a goroutine to handle updates in channel updateMessageChan
			go func() {
				for result := range updateMessageChan {
					// Complete the status message
					if len(result.Medias) > 0 {
						// Create media group for multiple medias
						var mediaGroup []models.InputMedia

						for _, media := range result.Medias {
							mediaType := DetectFileType(media.Filename)
							mediaType := download.DetectFileType(media.Filename)
							switch mediaType {
							case "video":
								inputVideo := &models.InputMediaVideo{
									Media:           fmt.Sprintf("attach://%s", media.Filename),
									MediaAttachment: media.Media,
									Caption:         fmt.Sprintf("%d. %s\nState: %s", i+1, url, result.State),
								}
								mediaGroup = append(mediaGroup, inputVideo)
							case "photo":
								inputPhoto := &models.InputMediaPhoto{
									Media:           fmt.Sprintf("attach://%s", media.Filename),
									MediaAttachment: media.Media,
									Caption:         fmt.Sprintf("%d. %s\nState: %s", i+1, url, result.State),
								}
								mediaGroup = append(mediaGroup, inputPhoto)
							default:
								logger.Log.Sugar().Errorf("Unsupported media type for file %s", media.Filename)
								continue
							}
						}

						_, err := b.bot.SendMediaGroup(ctx, &bot.SendMediaGroupParams{
							ChatID: update.Message.Chat.ID,
							Media:  mediaGroup,
							ReplyParameters: &models.ReplyParameters{
								MessageID: update.Message.ID,
							},
						})
						if err != nil {
							logger.Log.Sugar().Errorf("Failed to send media group: %v", err)

							state := fmt.Sprintf("❌ failed to send media: %v", err)
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
							// Successfully sent media, now delete the status message
							_, err = b.bot.DeleteMessage(ctx, &bot.DeleteMessageParams{
								ChatID:    update.Message.Chat.ID,
								MessageID: statusMsg.ID,
							})
							if err != nil {
								logger.Log.Sugar().Errorf("Failed to delete status message: %v", err)
							}
						}

						// Close all media streams
						for _, media := range result.Medias {
							media.Media.Close()
						}
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
				updateMessageChan <- MediaResult{
					State: err.Error(),
				}
			}
		}(url)
	}
}

func (b *DefaultBot) handleURL(url string, updateMessageChan chan MediaResult) error {
	attempts := config.GetConfig().VideoSaver.RetryCount

	if err := common.DoWithRetry(common.RetryConfig{
		Attempts: attempts,
		Delay:    2 * time.Second,
	}, func() error {
		var (
			directUrls []string
		)

		saver, err := b.GetVideoSaver(url)
		if err != nil {
			return fmt.Errorf("failed to get media saver: %w", err)
		}

		if err := b.browserPool.UseBrowser(func(ctx context.Context, browser *browserpool.Browser) error {
			updateMessageChan <- MediaResult{
				State: "🔎 getting info...",
			}
			logger.Log.Sugar().Infof("Processing URL: %s", url)

			directUrls, err = saver.GetVideoURLs(ctx, browser.Browser, url)
			if err != nil {
				return fmt.Errorf("failed to get direct video URL: %w", err)
			}

			return nil
		}); err != nil {
			return err
		}

		var totalSize int64
		var medias []MediaData

		for j, directUrl := range directUrls {
			fileSize, _ := download.GetFileSize(directUrl)
			totalSize += fileSize
			sizeStr := getSizeStr(fileSize)

			// Update state to show which URL is being downloaded
			updateMessageChan <- MediaResult{
				State: fmt.Sprintf("⬇️ downloading media %d/%d...%s", j+1, len(directUrls), sizeStr),
			}

			req, err := http.NewRequest("GET", directUrl, nil)
			if err != nil {
				return fmt.Errorf("failed to create request for direct URL %s: %w", directUrl, err)
			}

			// Mimic a real browser
			logger.Log.Sugar().Infof("downloading media from %s with user agent %s", directUrl, saver.GetUA())
			req.Header.Set("User-Agent", saver.GetUA())
			req.Header.Set("Accept", "*/*")
			req.Header.Set("Accept-Language", "en-US,en;q=0.9,ru;q=0.8")
			req.Header.Set("Accept-Encoding", "identity")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("failed to download video from %s: %w", directUrl, err)
			}
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("failed to download video from %s: HTTP %d %s", directUrl, resp.StatusCode, resp.Status)
			}

			medias = append(medias, MediaData{
				Filename: saver.GetFilename(url, directUrl),
				Media:    resp.Body,
			})
		}

		sizeStr := getSizeStr(totalSize)

		// Show final success state with count of medias
		successState := fmt.Sprintf("✅ get media successfully%s", sizeStr)
		if len(medias) > 1 {
			successState = fmt.Sprintf("✅ got %d medias successfully%s", len(medias), sizeStr)
		}

		updateMessageChan <- MediaResult{
			Medias: medias,
			State:  successState,
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to get media info after %d attempts: %w", attempts, err)
	}

	return nil
}

func getSizeStr(fileSize int64) string {
	if fileSize == 0 {
		return ""
	}
	return fmt.Sprintf(" (%s)", download.ByteCountBinary(fileSize))
}
