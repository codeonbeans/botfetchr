package instagram

import (
	videosavermdl "botvideosaver/internal/client/videosaver/model"
	"botvideosaver/internal/logger"
	"botvideosaver/internal/utils/common"
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/proto"
)

var shortCodeRegex = regexp.MustCompile(`/(p|tv|reel|reels(?:/videos)?)/([A-Za-z0-9-_]+)`)

type clientImpl struct {
	*videosavermdl.BaseClientImpl
}

func NewClient(browser *rod.Browser) *clientImpl {
	return &clientImpl{
		BaseClientImpl: videosavermdl.NewBaseClient(browser),
	}
}

func (c *clientImpl) GetVideoURL(ctx context.Context, url string) (videoURL string, err error) {
	logger.Log.Sugar().Infof("Opening page %s with user agent %s", url, c.UA)

	fmt.Println("Opening page with user agent:", c.UA, "and timeout:", c.Timeout)

	page, cancel := c.Browser.
		Context(ctx).
		MustPage(url).
		MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
			UserAgent: c.UA,
		}).
		WithCancel()
	defer page.Close()

	go func() {
		time.Sleep(c.Timeout)
		cancel()
	}()

	logger.Log.Sugar().Infof("Setting viewport for page %s", url)
	page.MustSetViewport(devices.Nexus5.Screen.Vertical.Width, devices.Nexus5.Screen.Vertical.Height, 1, true)

	page.MustReload()

	for {
		html := page.MustHTML()

		urls := extractVideoURLs(html)
		if len(urls) > 0 {
			var marshaledURL string
			if c.Quality == "high" {
				marshaledURL = urls[0] // First URL is the highest quality
			} else {
				marshaledURL = urls[len(urls)-1] // Last URL is the lowest quality
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

func (c *clientImpl) GetVideoID(url string) (string, error) {
	// Extract the video ID from the URL
	// This is a placeholder implementation; you may need to adjust it based on your URL structure
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid Instagram URL: %s", url)
	}
	videoID := parts[len(parts)-1]

	return videoID, nil
}

func (c *clientImpl) IsValidURL(url string) bool {
	// Check if the URL is a valid Instagram video URL
	return shortCodeRegex.Match([]byte(url))
}

func extractVideoURLs(text string) []string {
	// Regex pattern to match URLs within video_versions array
	pattern := `(?s)"video_versions":\s*\[(.*?)\]`
	videoVersionsRegex := regexp.MustCompile(pattern)

	// Find the video_versions section
	videoVersionsMatch := videoVersionsRegex.FindStringSubmatch(text)
	if len(videoVersionsMatch) < 2 {
		return nil
	}

	// Extract URLs from the video_versions section
	urlPattern := `"url":\s*"([^"]+)"`
	urlRegex := regexp.MustCompile(urlPattern)

	matches := urlRegex.FindAllStringSubmatch(videoVersionsMatch[1], -1)

	var urls []string
	for _, match := range matches {
		if len(match) > 1 {
			urls = append(urls, match[1])
		}
	}

	return urls
}
