package vk

import (
	"botvideosaver/internal/logger"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/go-rod/rod"
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

func (c *clientImpl) getOidAndId(url string) (ownerID, videoID string, err error) {
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

func (c *clientImpl) GetVideoURL(urlText string) (videoURL string, err error) {
	ownerID, videoID, err := c.getOidAndId(urlText)
	if err != nil {
		return "", fmt.Errorf("failed to parse VK video URL: %w", err)
	}

	embedUrl := fmt.Sprintf("https://vkvideo.ru/video_ext.php?oid=-%s&id=%s", ownerID, videoID)

	logger.Log.Sugar().Infof("Opening page %s with user agent %s", embedUrl, c.ua)

	page, cancel := c.browser.
		MustPage(embedUrl).
		MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36",
		}).
		WithCancel()
	defer page.Close()

	// logger.Log.Sugar().Infof("Setting viewport for page %s", url)
	// page.MustSetViewport(devices.Nexus5.Screen.Vertical.Width, devices.Nexus5.Screen.Vertical.Height, 1, true)

	go func() {
		time.Sleep(TIMEOUT)
		cancel()
	}()

	page.MustReload()

	for {
		html := page.MustHTML()
		re := regexp.MustCompile(`"url\d+":"([^"]+)"`)
		matches := re.FindAllStringSubmatch(html, -1)

		logger.Log.Sugar().Infof("Found %d video URL matches in the page", len(matches))
		if len(matches) > 0 {
			// Take the best quality video (usually 1080)
			lastMatch := matches[len(matches)-1][1] // group 1 is the URL
			lastMatch = matches[0][1]

			fmt.Println("NIGGA", fmt.Sprintf(`"%s"`, lastMatch))

			url, err := UnmarshalURL(fmt.Sprintf(`"%s"`, lastMatch))
			if err != nil {
				return "", fmt.Errorf("failed to parse video URL: %w", err)
			}

			url, err = FixEscapedURL(url)
			if err != nil {
				return "", fmt.Errorf("failed to parse video URL: %w", err)
			}

			fmt.Println("final url", url)

			return url, err
		}

		time.Sleep(1 * time.Second)
	}
}

func (c *clientImpl) GetVideoID(url string) (videoID string, err error) {
	ownerID, videoID, err := c.getOidAndId(url)
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

func UnmarshalURL(marshalledURL string) (string, error) {
	var result string
	err := json.Unmarshal([]byte(marshalledURL), &result)
	return result, err
}

func FixEscapedURL(escapedURL string) (string, error) {
	// Replace escaped forward slashes
	fixedURL := strings.ReplaceAll(escapedURL, `\/`, `/`)

	// Validate the URL
	_, err := url.Parse(fixedURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL after fixing: %v", err)
	}

	return fixedURL, nil
}
