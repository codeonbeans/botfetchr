package tgbot

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func getFileSizeMB(url string) (float64, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	// For Odnoklassniki CDN, we need specific headers
	if strings.Contains(url, "okcdn.ru") {
		return getOKCDNFileSize(client, url)
	}

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
