package tgbot

import (
	"botmediasaver/config"
	"botmediasaver/internal/client/browserpool"
	"botmediasaver/internal/logger"
	"botmediasaver/internal/utils/common"
	"botmediasaver/internal/utils/download"
	"botmediasaver/internal/utils/ptr"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// MediaProcessor handles the processing of individual URLs
type MediaProcessor struct {
	bot        *DefaultBot
	processCtx *ProcessingContext
	updateChan chan MediaResult
}

func (mp *MediaProcessor) handleStatusUpdates() {
	for result := range mp.updateChan {
		mp.updateStatus(result)
	}
}

func (mp *MediaProcessor) processURL() error {
	attempts := config.GetConfig().MediaSaver.RetryCount

	return common.DoWithRetry(common.RetryConfig{
		Attempts: attempts,
		Delay:    2 * time.Second,
	}, mp.attemptDownload)
}

func (mp *MediaProcessor) attemptDownload() error {
	// Get video saver
	saver, err := mp.bot.GetMediaSaver(mp.processCtx.url)
	if err != nil {
		return fmt.Errorf("failed to get media saver: %w", err)
	}

	// Get direct URLs
	directUrls, err := mp.getDirectURLs(saver)
	if err != nil {
		return err
	}

	// Download media files
	medias, err := mp.downloadMedias(saver, directUrls)
	if err != nil {
		return err
	}

	// Send success result
	mp.sendSuccessResult(medias)
	return nil
}

func (mp *MediaProcessor) getDirectURLs(saver MediaSaver) ([]string, error) {
	var directUrls []string

	err := mp.bot.browserPool.UseBrowser(func(ctx context.Context, browser *browserpool.Browser) error {
		mp.updateChan <- MediaResult{State: "🔎 getting info..."}
		logger.Log.Sugar().Infof("Processing URL: %s", mp.processCtx.url)

		var err error
		directUrls, err = saver.GetVideoURLs(ctx, browser.Browser, mp.processCtx.url)

		if err != nil {
			return fmt.Errorf("failed to get direct video URL: %w", err)
		}

		return nil
	})

	return directUrls, err
}

func (mp *MediaProcessor) downloadMedias(saver MediaSaver, directUrls []string) ([]MediaData, error) {
	var medias []MediaData

	for j, directUrl := range directUrls {
		media, err := mp.downloadSingleMedia(saver, directUrl, j, len(directUrls))
		if err != nil {
			return nil, err
		}
		medias = append(medias, media)
	}

	return medias, nil
}

func (mp *MediaProcessor) downloadSingleMedia(saver MediaSaver, directUrl string, index, total int) (MediaData, error) {
	fileSize, _ := download.GetFileSize(directUrl)
	sizeStr := getSizeStr(fileSize)

	// Update download progress
	mp.updateChan <- MediaResult{
		State: fmt.Sprintf("⬇️ downloading media %d/%d...%s", index+1, total, sizeStr),
	}

	// Create and configure request
	req, err := http.NewRequest("GET", directUrl, nil)
	if err != nil {
		return MediaData{}, fmt.Errorf("failed to create request for direct URL %s: %w", directUrl, err)
	}

	mp.configureRequest(req, saver, directUrl)

	// Execute request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return MediaData{}, fmt.Errorf("failed to download video from %s: %w", directUrl, err)
	}

	if resp.StatusCode != http.StatusOK {
		return MediaData{}, fmt.Errorf("failed to download video from %s: HTTP %d %s", directUrl, resp.StatusCode, resp.Status)
	}

	return MediaData{
		Filename:  saver.GetFilename(mp.processCtx.url, directUrl),
		Size:      fileSize,
		Media:     resp.Body,
		DirectURL: directUrl,
	}, nil
}

func (mp *MediaProcessor) configureRequest(req *http.Request, saver MediaSaver, directUrl string) {
	userAgent := saver.GetUA()
	logger.Log.Sugar().Infof("downloading media from %s with user agent %s", directUrl, userAgent)

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,ru;q=0.8")
	req.Header.Set("Accept-Encoding", "identity")
}

func (mp *MediaProcessor) sendSuccessResult(medias []MediaData) {
	totalSize := mp.calculateTotalSize(medias)
	sizeStr := getSizeStr(totalSize)

	mp.updateChan <- MediaResult{
		State: fmt.Sprintf("📲 sending media%s", sizeStr),
	}

	successState := mp.getSuccessMessage(medias, sizeStr)
	mp.updateChan <- MediaResult{
		Medias: medias,
		State:  successState,
	}
}

func (mp *MediaProcessor) calculateTotalSize(medias []MediaData) int64 {
	var totalSize int64
	for _, media := range medias {
		totalSize += media.Size
	}
	return totalSize
}

func (mp *MediaProcessor) getSuccessMessage(medias []MediaData, sizeStr string) string {
	if len(medias) > 1 {
		return fmt.Sprintf("✅ got %d medias successfully%s", len(medias), sizeStr)
	}
	return fmt.Sprintf("✅ get media successfully%s", sizeStr)
}

func (mp *MediaProcessor) updateStatus(result MediaResult) {
	defer mp.closeMediaStreams(result.Medias)

	if len(result.Medias) > 0 {
		mp.handleMediaSending(result)
	} else {
		mp.updateStatusMessage(result.State)
	}
}

func (mp *MediaProcessor) closeMediaStreams(medias []MediaData) {
	for _, media := range medias {
		if media.Media != nil {
			media.Media.Close()
		}
	}
}

func (mp *MediaProcessor) handleMediaSending(result MediaResult) {
	groups := mp.createMediaGroups(result.Medias)

	if err := mp.sendMediaGroups(groups); err != nil {
		mp.updateStatusMessage(fmt.Sprintf("❌ failed to send media: %v", err))
	} else {
		mp.deleteStatusMessage()
	}
}

func (mp *MediaProcessor) createMediaGroups(medias []MediaData) [][]models.InputMedia {
	maxGroupSize := int64(config.GetConfig().MediaSaver.MaxGroupMediaSize * 1024 * 1024)

	var groups [][]models.InputMedia
	var currentGroup []models.InputMedia
	var currentGroupSize int64

	for mediaIdx, media := range medias {
		if media.Size >= maxGroupSize {
			mp.sendOversizedMediaURL(media, mediaIdx)
			continue
		}

		inputMedia := mp.createInputMedia(media)
		if inputMedia == nil {
			continue
		}

		if currentGroupSize+media.Size > maxGroupSize && len(currentGroup) > 0 {
			groups = append(groups, currentGroup)
			currentGroup = nil
			currentGroupSize = 0
		}

		currentGroup = append(currentGroup, inputMedia)
		currentGroupSize += media.Size
	}

	if len(currentGroup) > 0 {
		groups = append(groups, currentGroup)
	}

	return groups
}

func (mp *MediaProcessor) sendOversizedMediaURL(media MediaData, index int) {
	text := fmt.Sprintf("%d. %s\nFile (%d) too large to send directly (%.2f MB). Direct URL: %s",
		mp.processCtx.urlIndex+1, mp.processCtx.url, index+1,
		float64(media.Size)/1024.0/1024.0, media.DirectURL)

	_, err := mp.bot.SendMessage(mp.processCtx.ctx, &bot.SendMessageParams{
		ChatID: mp.processCtx.chatID,
		Text:   text,
		ReplyParameters: &models.ReplyParameters{
			MessageID: mp.processCtx.originalMsgID,
		},
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: ptr.ToPtr(true),
		},
	})

	if err != nil {
		logger.Log.Sugar().Errorf("Failed to send DirectURL message: %v", err)
	}
}

func (mp *MediaProcessor) createInputMedia(media MediaData) models.InputMedia {
	mediaType := download.DetectFileType(media.Filename)

	switch mediaType {
	case "video":
		return &models.InputMediaVideo{
			Media:           fmt.Sprintf("attach://%s", media.Filename),
			MediaAttachment: media.Media,
		}
	case "photo":
		return &models.InputMediaPhoto{
			Media:           fmt.Sprintf("attach://%s", media.Filename),
			MediaAttachment: media.Media,
		}
	default:
		logger.Log.Sugar().Errorf("Unsupported media type for file %s", media.Filename)
		return nil
	}
}

func (mp *MediaProcessor) sendMediaGroups(groups [][]models.InputMedia) error {
	for _, group := range groups {
		if len(group) == 0 {
			continue
		}

		_, err := mp.bot.SendMediaGroup(mp.processCtx.ctx, &bot.SendMediaGroupParams{
			ChatID: mp.processCtx.chatID,
			Media:  group,
			ReplyParameters: &models.ReplyParameters{
				MessageID: mp.processCtx.originalMsgID,
			},
		})

		if err != nil {
			logger.Log.Sugar().Errorf("Failed to send media group: %v", err)
			return err
		}
	}

	return nil
}

func (mp *MediaProcessor) updateStatusMessage(state string) {
	text := fmt.Sprintf("%d. %s\nState: %s", mp.processCtx.urlIndex+1, mp.processCtx.url, state)

	_, err := mp.bot.EditMessageText(mp.processCtx.ctx, &bot.EditMessageTextParams{
		ChatID:    mp.processCtx.chatID,
		MessageID: mp.processCtx.statusMsg.ID,
		Text:      text,
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: ptr.ToPtr(true),
		},
	})

	if err != nil {
		logger.Log.Sugar().Errorf("Failed to edit message text: \"%s\" %v", text, err)
	}
}

func (mp *MediaProcessor) deleteStatusMessage() {
	_, err := mp.bot.DeleteMessage(mp.processCtx.ctx, &bot.DeleteMessageParams{
		ChatID:    mp.processCtx.chatID,
		MessageID: mp.processCtx.statusMsg.ID,
	})

	if err != nil {
		logger.Log.Sugar().Errorf("Failed to delete status message: %v", err)
	}
}

func getSizeStr(fileSize int64) string {
	if fileSize == 0 {
		return ""
	}
	return fmt.Sprintf(" (%s)", download.ByteCountBinary(fileSize))
}
