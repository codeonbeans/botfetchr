package videosavermdl

import (
	"botvideosaver/internal/logger"
	"time"

	"github.com/go-rod/rod"
)

const DEFAULT_TIMEOUT = 30 * time.Second
const DEFAULT_USER_AGENT = "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; Trident/4.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; Maxthon 2.0)"

type BaseClientImpl struct {
	UA         string
	Quality    string        // Quality can be "low" or "high"
	RetryCount int           // Number of retries for fetching video
	Timeout    time.Duration // Timeout for each task

	Browser *rod.Browser
}

func NewBaseClient(browser *rod.Browser) *BaseClientImpl {
	return &BaseClientImpl{
		UA:         DEFAULT_USER_AGENT,
		Quality:    "high",          // Default quality
		RetryCount: 3,               // Default retry count
		Timeout:    DEFAULT_TIMEOUT, // Default timeout
		Browser:    browser,
	}
}

func (c *BaseClientImpl) SetUserAgent(ua string) {
	// Set the user agent for the browser
	if ua == "" {
		logger.Log.Sugar().Warn("User agent is empty, using default user agent")
		c.UA = DEFAULT_USER_AGENT
	} else {
		c.UA = ua
	}
}

func (c *BaseClientImpl) SetQuality(quality string) {
	// Set the video quality
	if quality == "low" || quality == "high" {
		c.Quality = quality
	} else {
		logger.Log.Sugar().Warnf("Invalid video quality: %s, defaulting to 'high'", quality)
		c.Quality = "high"
	}
}

func (c *BaseClientImpl) SetRetryCount(count int) {
	// Set the number of retries for failed downloads
	if count < 0 {
		logger.Log.Sugar().Warnf("Invalid retry count: %d, defaulting to 3", count)
		c.RetryCount = 3
	} else {
		c.RetryCount = count
	}
}

func (c *BaseClientImpl) SetTimeout(timeout time.Duration) {
	// Set the timeout for each download
	if timeout <= 0 {
		logger.Log.Sugar().Warnf("Invalid timeout: %s, defaulting to 30 seconds", timeout)
		c.Timeout = DEFAULT_TIMEOUT
	} else {
		c.Timeout = timeout
	}
}

func (c *BaseClientImpl) IsValidURL(url string) bool {
	// This method should be implemented by each specific client
	// to validate the URL format for the respective service
	return false
}
