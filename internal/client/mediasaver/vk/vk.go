package vk

import (
	videosavermdl "botvideosaver/internal/client/mediasaver/base"
	"botvideosaver/internal/logger"
	"botvideosaver/internal/utils/common"
	"botvideosaver/internal/utils/download"
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/google/uuid"
)

var shortCodeRegex = regexp.MustCompile(`https:\/\/(m\.)?vkvideo\.ru\/video-(\d+)_(\d+)`)

type clientImpl struct {
	*videosavermdl.BaseClientImpl
}

func NewClient() *clientImpl {
	return &clientImpl{
		BaseClientImpl: videosavermdl.NewBaseClient(),
	}
}

func (c *clientImpl) GetVideoURLs(ctx context.Context, browser *rod.Browser, urlText string) (videoURLs []string, err error) {
	ownerID, videoID, err := getOidAndId(urlText)
	if err != nil {
		return nil, fmt.Errorf("failed to parse VK video URL: %w", err)
	}

	embedUrl := fmt.Sprintf("https://vkvideo.ru/video_ext.php?oid=-%s&id=%s", ownerID, videoID)

	logger.Log.Sugar().Infof("Opening page %s with user agent %s", embedUrl, c.UA)

	c.SetUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")
	page, cancel := browser.
		Context(ctx).
		MustPage(embedUrl).
		MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
			UserAgent: c.UA,
		}).
		WithCancel()
	defer page.Close()

	go func() {
		time.Sleep(c.Timeout)
		cancel()
	}()

	page.MustReload()

	for {
		html := page.MustHTML()
		urls := extractVideoURLs(html)
		if len(urls) > 0 {
			var marshaledURL string
			if c.Quality == "high" {
				marshaledURL = urls[len(urls)-1] // Last URL is the highest quality
			} else {
				marshaledURL = urls[0] // First URL is the lowest quality
			}

			url, err := common.UnmarshalURL(marshaledURL)
			if err != nil {
				return nil, fmt.Errorf("failed to parse video URL: %w", err)
			}

			return []string{url}, nil
		}

		time.Sleep(1 * time.Second)
	}
}

func (c *clientImpl) GetFilename(ogUrl, directUrl string) string {
	var fileID string
	oid, id, err := getOidAndId(ogUrl)
	if err != nil {
		fileID = uuid.NewString()
	} else {
		fileID = fmt.Sprintf("%s_%s", oid, id) // Use owner ID
	}

	filename := fmt.Sprintf("vk_%s_%s%s", c.Quality, fileID, filepath.Ext(download.GetFileName(directUrl)))
	return filename
}

func (c *clientImpl) IsValidURL(url string) bool {
	// Check if the URL matches the VK video format (both desktop and mobile)
	return shortCodeRegex.MatchString(url)
}

func getOidAndId(url string) (ownerID, videoID string, err error) {
	// Get owner ID and video ID from the URL
	// Pattern: https://vk.com/video-123456_7890123 or https://m.vkvideo.ru/video-123456_7890123
	re := regexp.MustCompile(`video-(\d+)_(\d+)`)
	matches := re.FindStringSubmatch(url)

	if len(matches) == 3 {
		ownerID = matches[1]
		videoID = matches[2]
		return ownerID, videoID, nil
	}

	return "", "", fmt.Errorf("invalid VK video URL format")
}

func extractVideoURLs(text string) []string {
	urlRegex := regexp.MustCompile(`"url\d+":"([^"]+)"`)

	matches := urlRegex.FindAllStringSubmatch(text, -1)

	var urls []string
	for _, match := range matches {
		if len(match) > 1 {
			urls = append(urls, match[1])
		}
	}

	return urls
}
