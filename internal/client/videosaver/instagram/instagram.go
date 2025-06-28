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

func (c *clientImpl) GetVideoURLs(ctx context.Context, url string) (videoURLs []string, err error) {
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
	page.MustSetViewport(100, 100, 1, true)

	page.MustReload()

	for {
		html := page.MustHTML()

		var urls []string

		videoUrls := extractVideoURLs(html)
		if len(videoUrls) > 0 {
			var marshaledURL string
			if c.Quality == "high" {
				marshaledURL = videoUrls[0] // First URL is the highest quality
			} else {
				marshaledURL = videoUrls[len(videoUrls)-1] // Last URL is the lowest quality
			}

			url, err := common.UnmarshalURL(marshaledURL)
			if err != nil {
				return nil, fmt.Errorf("failed to parse video URL: %w", err)
			}

			urls = append(urls, url)
		}

		imageUrls := extractImageURLs(html)
		if len(imageUrls) > 0 {
			for _, marshaledURL := range imageUrls {
				url, err := common.UnmarshalURL(marshaledURL)
				if err != nil {
					return nil, fmt.Errorf("failed to parse image URL: %w", err)
				}
				urls = append(urls, url)
			}
		}

		if len(urls) > 0 {
			return urls, nil
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

func extractImageURLs(text string) []string {
	// Regex pattern to match image_versions2 section that comes after original_width
	pattern := `(?s)"original_width":[^,}]*,.*?"image_versions2":\{"candidates":\[\{[^}]*"url":"([^"]+)"`
	urlRegex := regexp.MustCompile(pattern)

	matches := urlRegex.FindAllStringSubmatch(text, -1)

	var urls []string
	for _, match := range matches {
		if len(match) > 1 {
			urls = append(urls, match[1])
		}
	}

	return urls
}
