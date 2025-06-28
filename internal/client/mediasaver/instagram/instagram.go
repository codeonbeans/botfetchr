package instagram

import (
	videosavermdl "botvideosaver/internal/client/mediasaver/base"
	"botvideosaver/internal/logger"
	"botvideosaver/internal/utils/common"
	"botvideosaver/internal/utils/download"
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/google/uuid"
)

var shortCodeRegex = regexp.MustCompile(`/(p|tv|reel|reels(?:/videos)?)/([A-Za-z0-9-_]+)`)

type clientImpl struct {
	*videosavermdl.BaseClientImpl
}

func NewClient() *clientImpl {
	return &clientImpl{
		BaseClientImpl: videosavermdl.NewBaseClient(),
	}
}

func (c *clientImpl) GetVideoURLs(ctx context.Context, browser *rod.Browser, url string) (videoURLs []string, err error) {
	logger.Log.Sugar().Infof("Opening page %s with user agent %s", url, c.UA)
	page, cancel := browser.
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
	page.MustSetViewport(1000, 1000, 1, true)

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

		if isPost(url) {
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
		}

		if len(urls) > 0 {
			return urls, nil
		}

		time.Sleep(1 * time.Second)
	}
}

func (c *clientImpl) GetFilename(ogUrl, directUrl string) string {
	var fileID string
	matches := shortCodeRegex.FindStringSubmatch(ogUrl)
	if len(matches) < 3 || isPost(ogUrl) {
		fileID = uuid.NewString()
	} else {
		fileID = matches[2]
	}

	filename := fmt.Sprintf("instagram_%s_%s%s", c.Quality, fileID, filepath.Ext(download.GetFileName(directUrl)))
	return filename
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

func isPost(url string) bool {
	return strings.Contains(url, "/p/") || strings.Contains(url, "/tv/")
}

func isReel(url string) bool {
	return strings.Contains(url, "/reel/") || strings.Contains(url, "/reels/videos/")
}
