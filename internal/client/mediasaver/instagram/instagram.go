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

func (c *clientImpl) GetVideoURLs(ctx context.Context, browser *rod.Browser, ogUrl string) (videoURLs []string, err error) {
	logger.Log.Sugar().Infof("Opening page %s with user agent %s", ogUrl, c.UA)
	page, cancel := browser.
		Context(ctx).
		MustPage(ogUrl).
		MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
			UserAgent: c.UA,
		}).
		WithCancel()
	defer page.Close()

	go func() {
		time.Sleep(c.Timeout)
		cancel()
	}()

	logger.Log.Sugar().Infof("Setting viewport for page %s", ogUrl)
	page.MustSetViewport(1000, 1000, 1, true)

	page.MustReload()

	for {
		html := page.MustHTML()

		var urls []string

		videoUrls := c.extractVideoURLs(html)
		for _, videoUrl := range videoUrls {
			url, err := common.UnmarshalURL(videoUrl)
			if err != nil {
				continue
			}
			urls = append(urls, url)

			// Only take one video URL if it's a reel
			if isReel(ogUrl) {
				break
			}
		}

		if isPost(ogUrl) {
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

func (c *clientImpl) extractVideoURLs(text string) []string {
	// Regex pattern to match URLs within video_versions array
	pattern := `"video_versions"\s*:\s*\[(.*?)\]`
	videoVersionsRegex := regexp.MustCompile(pattern)

	// Find all video_versions sections
	videoVersionsMatches := videoVersionsRegex.FindAllStringSubmatch(text, -1)
	if len(videoVersionsMatches) == 0 {
		return nil
	}

	var urls []string

	// Extract URLs from each video_versions section
	urlPattern := `"url"\s*:\s*"([^"]+)"`
	urlRegex := regexp.MustCompile(urlPattern)

	for _, videoVersionsMatch := range videoVersionsMatches {
		if len(videoVersionsMatch) < 2 {
			continue
		}

		matches := urlRegex.FindAllStringSubmatch(videoVersionsMatch[1], -1)
		// for _, match := range matches {
		// 	if len(match) > 1 {
		// 		urls = append(urls, UnmarshalURL(match[1]))
		// 	}
		// }
		if len(matches) > 1 {
			//First url is highest quality, last is lowest quality
			// urls = append(urls, UnmarshalURL(matches[0][1]))
			if c.Quality == "high" {
				urls = append(urls, matches[0][1]) // First URL is the highest quality
			} else {
				urls = append(urls, matches[len(matches)-1][1]) // Last URL is the lowest quality
			}
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
	// Match common Instagram post patterns
	return strings.Contains(url, "/p/") || strings.Contains(url, "/tv/") || strings.Contains(url, "/post/")
}

func isReel(url string) bool {
	return strings.Contains(url, "/reel/") || strings.Contains(url, "/reels/")
}
