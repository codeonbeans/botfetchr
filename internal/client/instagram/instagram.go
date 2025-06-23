package instagram

import (
	"botvideosaver/internal/logger"
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/proto"
)

const (
	TIMEOUT = 10 * time.Second
)

type clientImpl struct {
	ua      string
	browser *rod.Browser
}

func NewClient(ua string, browser *rod.Browser) (*clientImpl, error) {
	return &clientImpl{
		ua:      ua,
		browser: browser,
	}, nil
}

func (c *clientImpl) getVideoUrlFallback(url string) (videoURL string, err error) {
	logger.Log.Sugar().Infof("Opening page %s with user agent %s", url, c.ua)

	page, cancel := c.browser.
		MustPage(url).
		MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
			UserAgent: c.ua,
		}).
		WithCancel()
	defer page.Close()

	logger.Log.Sugar().Infof("Setting viewport for page %s", url)
	page.MustSetViewport(devices.Nexus5.Screen.Vertical.Width, devices.Nexus5.Screen.Vertical.Height, 1, true)

	go func() {
		time.Sleep(TIMEOUT)
		cancel()
	}()

	page.MustReload()

	logger.Log.Sugar().Infof("Waiting for video element on page %s", url)
	elem := page.MustElement("video")

	logger.Log.Sugar().Infof("Found video element on page %s", url)
	src := elem.MustAttribute("src")

	if src == nil {
		return "", fmt.Errorf("video element has no src attribute")
	}

	return *src, nil
}

func (c *clientImpl) GetVideoURL(url string) (videoURL string, err error) {
	code, err := ExtractShortCodeFromLink(url)
	if err != nil {
		return "", fmt.Errorf("failed to extract shortcode from link: %w", err)
	}

	extractor := NewExtractor()
	extract, err := extractor.GetPostWithCode(code)
	if err != nil {
		// If the extractor fails, try another method
		return c.getVideoUrlFallback(url)
	}

	return extract.URL, err
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
