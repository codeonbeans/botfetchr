package vk

import (
	videosavermdl "botvideosaver/internal/client/videosaver/model"
	"botvideosaver/internal/logger"
	"botvideosaver/internal/utils/common"
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type clientImpl struct {
	*videosavermdl.BaseClientImpl
}

func NewClient(browser *rod.Browser) *clientImpl {
	return &clientImpl{
		BaseClientImpl: videosavermdl.NewBaseClient(browser),
	}
}

func (c *clientImpl) GetVideoURL(ctx context.Context, urlText string) (videoURL string, err error) {
	ownerID, videoID, err := getOidAndId(urlText)
	if err != nil {
		return "", fmt.Errorf("failed to parse VK video URL: %w", err)
	}

	embedUrl := fmt.Sprintf("https://vkvideo.ru/video_ext.php?oid=-%s&id=%s", ownerID, videoID)

	logger.Log.Sugar().Infof("Opening page %s with user agent %s", embedUrl, c.UA)

	page := c.Browser.
		MustPage(embedUrl).
		MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36",
		}).
		Timeout(c.Timeout).
		Context(ctx)
	defer page.Close()

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
				return "", fmt.Errorf("failed to parse video URL: %w", err)
			}

			return url, nil
		}

		time.Sleep(1 * time.Second)
	}
}

func (c *clientImpl) GetVideoID(url string) (videoID string, err error) {
	ownerID, videoID, err := getOidAndId(url)
	if err != nil {
		return "", fmt.Errorf("failed to parse VK video URL: %w", err)
	}

	return fmt.Sprintf("%s_%s", ownerID, videoID), nil
}

func (c *clientImpl) IsValidURL(url string) bool {
	// Check if the URL matches the VK video format
	re := regexp.MustCompile(`https:\/\/vkvideo\.ru\/video-(\d+)_(\d+)`)
	return re.MatchString(url)
}

func getOidAndId(url string) (ownerID, videoID string, err error) {
	// Get owner ID and video ID from the URL
	// Pattern: https://vk.com/video-123456_7890123
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
