package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func getFileSizeMB(url string) (float64, error) {
	// Check if URL is expired first
	if isURLExpired(url) {
		return 0, fmt.Errorf("URL has expired")
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			fmt.Printf("Redirected to: %s\n", req.URL)
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	// For Odnoklassniki CDN, we need specific headers
	if strings.Contains(url, "okcdn.ru") {
		fmt.Printf("For Odnoklassniki CDN, we need specific headers\n")
		return getOKCDNFileSize(client, url)
	}

	fmt.Printf("For other URLs, we use standard approach\n")
	// Standard approach for other URLs
	return getStandardFileSize(client, url)
}

func getOKCDNFileSize(client *http.Client, url string) (float64, error) {
	// Method 1: Try HEAD request with proper headers
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create HEAD request: %w", err)
	}

	// Add headers that OK CDN expects
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,ru;q=0.8")
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("HEAD request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		contentLength := resp.Header.Get("Content-Length")
		if contentLength != "" && contentLength != "0" {
			bytes, err := strconv.ParseInt(contentLength, 10, 64)
			if err == nil && bytes > 0 {
				return float64(bytes) / (1024 * 1024), nil
			}
		}
	}

	// Method 2: Try GET request (some CDNs only provide Content-Length on GET)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create GET request: %w", err)
	}

	// Same headers as HEAD request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,ru;q=0.8")
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Referer", "https://ok.ru/")

	resp, err = client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("GET request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// Check Content-Length from GET response
	contentLength := resp.Header.Get("Content-Length")
	if contentLength != "" && contentLength != "0" {
		bytes, err := strconv.ParseInt(contentLength, 10, 64)
		if err == nil && bytes > 0 {
			return float64(bytes) / (1024 * 1024), nil
		}
	}

	// Method 3: Read the actual content to determine size
	return readContentSize(resp.Body)
}

func getStandardFileSize(client *http.Client, url string) (float64, error) {
	// Standard HEAD request
	resp, err := client.Head(url)
	if err != nil {
		return 0, fmt.Errorf("HEAD request failed: %w", err)
	}
	defer resp.Body.Close()

	contentLength := resp.Header.Get("Content-Length")
	if contentLength == "" {
		return 0, fmt.Errorf("Content-Length header not found")
	}

	bytes, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid Content-Length: %w", err)
	}

	return float64(bytes) / (1024 * 1024), nil
}

func readContentSize(body io.ReadCloser) (float64, error) {
	const maxSize = 500 * 1024 * 1024 // 500MB limit
	bytesRead, err := io.Copy(io.Discard, io.LimitReader(body, maxSize))
	if err != nil {
		return 0, fmt.Errorf("failed to read content: %w", err)
	}

	mb := float64(bytesRead) / (1024 * 1024)
	if bytesRead == maxSize {
		return mb, fmt.Errorf("file size exceeds 500MB limit (downloaded: %.2f MB)", mb)
	}

	return mb, nil
}

func isURLExpired(url string) bool {
	if !strings.Contains(url, "expires=") {
		return false
	}

	parts := strings.Split(url, "expires=")
	if len(parts) < 2 {
		return false
	}

	expiresStr := strings.Split(parts[1], "&")[0]
	expires, err := strconv.ParseInt(expiresStr, 10, 64)
	if err != nil {
		return false
	}

	// Convert milliseconds to seconds
	expiryTime := time.Unix(expires/1000, 0)
	return time.Now().After(expiryTime)
}

// Alternative function that tries multiple strategies
func getFileSizeMBRobust(url string) (float64, error) {
	strategies := []string{"HEAD", "GET", "RANGE", "DOWNLOAD"}

	for _, strategy := range strategies {
		size, err := tryStrategy(url, strategy)
		if err == nil && size > 0 {
			return size, nil
		}
		fmt.Printf("Strategy %s failed: %v\n", strategy, err)
	}

	return 0, fmt.Errorf("all strategies failed to determine file size")
}

func tryStrategy(url, strategy string) (float64, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	switch strategy {
	case "HEAD":
		return tryHEAD(client, url)
	case "GET":
		return tryGET(client, url)
	case "RANGE":
		return tryRange(client, url)
	case "DOWNLOAD":
		return tryDownload(client, url)
	default:
		return 0, fmt.Errorf("unknown strategy: %s", strategy)
	}
}

func tryHEAD(client *http.Client, url string) (float64, error) {
	req, _ := http.NewRequest("HEAD", url, nil)
	addCommonHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return parseContentLength(resp.Header.Get("Content-Length"))
}

func tryGET(client *http.Client, url string) (float64, error) {
	req, _ := http.NewRequest("GET", url, nil)
	addCommonHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if size, err := parseContentLength(resp.Header.Get("Content-Length")); err == nil {
		return size, nil
	}

	return readContentSize(resp.Body)
}

func tryRange(client *http.Client, url string) (float64, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Range", "bytes=0-0")
	addCommonHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	contentRange := resp.Header.Get("Content-Range")
	if contentRange == "" {
		return 0, fmt.Errorf("no Content-Range header")
	}

	var totalSize int64
	n, err := fmt.Sscanf(contentRange, "bytes %*d-%*d/%d", &totalSize)
	if err != nil || n != 1 {
		return 0, fmt.Errorf("failed to parse Content-Range: %s", contentRange)
	}

	return float64(totalSize) / (1024 * 1024), nil
}

func tryDownload(client *http.Client, url string) (float64, error) {
	req, _ := http.NewRequest("GET", url, nil)
	addCommonHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return readContentSize(resp.Body)
}

func addCommonHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,ru;q=0.8")
	req.Header.Set("Accept-Encoding", "identity")

	if strings.Contains(req.URL.String(), "okcdn.ru") {
		req.Header.Set("Referer", "https://ok.ru/")
	}
}

func parseContentLength(contentLength string) (float64, error) {
	if contentLength == "" || contentLength == "0" {
		return 0, fmt.Errorf("no valid Content-Length")
	}

	bytes, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid Content-Length: %w", err)
	}

	if bytes <= 0 {
		return 0, fmt.Errorf("Content-Length is zero or negative")
	}

	return float64(bytes) / (1024 * 1024), nil
}

func main() {
	testURL := "https://vkvd196.okcdn.ru/?srcIp=58.187.246.12&pr=40&expires=1751409426919&srcAg=CHROME&fromCache=1&ms=45.136.22.132&type=5&sig=CkQmcYrc-RM&ct=0&urls=185.226.53.133&clientType=13&appId=512000384397&zs=65&id=8718605355762"

	fmt.Printf("Testing URL: %s\n\n", testURL)

	// Check if expired
	if isURLExpired(testURL) {
		fmt.Printf("⚠️  URL has expired!\n")
		return
	}

	// Try the optimized method
	size, err := getFileSizeMB(testURL)
	if err != nil {
		fmt.Printf("Standard method failed: %v\n", err)

		// Try robust method
		fmt.Printf("Trying robust method...\n")
		size, err = getFileSizeMBRobust(testURL)
		if err != nil {
			fmt.Printf("❌ All methods failed: %v\n", err)
			return
		}
	}

	fmt.Printf("✅ File size: %.2f MB\n", size)
}
