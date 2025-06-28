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
	"path/filepath"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
)

type VideoResult struct {
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
					if len(result.Medias) > 0 {
						// Create media group for multiple videos
						var mediaGroup []models.InputMedia

						for _, media := range result.Medias {
							mediaType := DetectFileType(media.Filename)
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
			directUrls []string
			ua         string
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

			directUrls, err = saver.GetVideoURLs(ctx, url)
			if err != nil {
				return fmt.Errorf("failed to get direct video URL: %w", err)
			}

			// VK client used fixed user agent, so we need to get ua here
			ua = saver.GetUA()

			return nil
		}); err != nil {
			return err
		}

		var totalSize int64
		var medias []MediaData

		for j, directUrl := range directUrls {
			fileSize, _ := getFileSize(directUrl)
			totalSize += fileSize
			sizeStr := getSizeStr(fileSize)

			// Update state to show which URL is being downloaded
			updateMessageChan <- VideoResult{
				State: fmt.Sprintf("⬇️ downloading media %d/%d...%s", j+1, len(directUrls), sizeStr),
			}

			req, err := http.NewRequest("GET", directUrl, nil)
			if err != nil {
				return fmt.Errorf("failed to create request for direct URL %s: %w", directUrl, err)
			}

			// Mimic a real browser
			logger.Log.Sugar().Infof("downloading media from %s with user agent %s", directUrl, ua)
			req.Header.Set("User-Agent", ua)
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

			filename := fmt.Sprintf("%s%s", uuid.NewString(), filepath.Ext(getFileName(directUrl)))

			medias = append(medias, MediaData{
				Filename: filename,
				Media:    resp.Body,
			})
		}

		sizeStr := getSizeStr(totalSize)

		// Show final success state with count of videos
		successState := fmt.Sprintf("✅ get video successfully%s", sizeStr)
		if len(medias) > 1 {
			successState = fmt.Sprintf("✅ got %d videos successfully%s", len(medias), sizeStr)
		}

		updateMessageChan <- VideoResult{
			Medias: medias,
			State:  successState,
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to get video info after %d attempts: %w", attempts, err)
	}

	return nil
}

func getSizeStr(fileSize int64) string {
	if fileSize == 0 {
		return " (unknown size)"
	}
	return fmt.Sprintf(" (%s)", ByteCountBinary(fileSize))
}
