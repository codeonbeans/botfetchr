package instagram

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/corpix/uarand"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

type Client interface {
	GetVideoURL(url string) (string, error)
	GetVideoID(url string) (string, error)
}

type clientImpl struct {
	m       *sync.Mutex
	browser *rod.Browser
}

func NewClient() Client {
	c := &clientImpl{
		m: &sync.Mutex{},
	}
	c.init()

	return c
}

func (c *clientImpl) init() {
	// Ensure the browser is initialized only once
	c.m.Lock()
	defer c.m.Unlock()

	chrome, found := launcher.LookPath()
	if !found {
		panic("could not find Chrome executable in PATH")
	}

	u := launcher.New().Bin(chrome).Headless(false).MustLaunch()
	c.browser = rod.New().ControlURL(u).MustConnect()
	c.browser.MustPage("about:blank").MustWaitStable()
}

func (c *clientImpl) GetVideoURL(url string) (string, error) {
	goodUA := false

	ua := uarand.GetRandom()

	defer func() {
		if !goodUA {
			fmt.Println("User agent not good, retrying with a new one")
			Append("bad_ua.txt", ua)
		}
	}()

	fmt.Println("Page with user agent:", ua)

	page, cancel := c.browser.
		MustPage(url).
		MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
			// UserAgent: "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36",
			UserAgent: uarand.GetRandom(),
		}).
		WithCancel()
	defer page.Close()

	fmt.Println("Emulating Nexus 5 device")
	page.MustSetViewport(devices.Nexus5.Screen.Vertical.Width, devices.Nexus5.Screen.Vertical.Height, 1, true)

	page.MustReload().MustWaitStable()

	fmt.Println("Waiting for video element")
	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()
	elem, err := page.Element("video")
	if err != nil {
		return "", fmt.Errorf("failed to find video element: %w", err)
	}

	fmt.Println("Getting video source URL")
	src, err := elem.Attribute("src")
	if err != nil {
		return "", fmt.Errorf("failed to get video source attribute: %w", err)
	}

	Append("good_ua.txt", ua)
	goodUA = true

	return *src, nil
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

func Append(path, text string) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	data := fmt.Sprintf("%s\n", text)

	// Write to the file
	if _, err := file.WriteString(data); err != nil {
		log.Fatalf("failed to write to file: %v", err)
	}
}
