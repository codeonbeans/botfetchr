package tgbot

import (
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// getFileName retrieves the filename from HTTP response headers
func getFileName(url string) string {
	// Try to get filename from headers first
	if filename := getFileNameFromHeaders(url); filename != "" {
		return filename
	}

	// Fallback to URL parsing if headers don't provide filename
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		lastPart := parts[len(parts)-1]
		// Remove query parameters
		if idx := strings.Index(lastPart, "?"); idx != -1 {
			lastPart = lastPart[:idx]
		}
		if lastPart != "" {
			return lastPart
		}
	}
	return "downloaded_file"
}

// getFileNameFromHeaders makes a HEAD request to get filename from response headers
func getFileNameFromHeaders(url string) string {
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	// Try HEAD request first
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return ""
	}

	// Add common headers that might be needed
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,ru;q=0.8")
	req.Header.Set("Accept-Encoding", "identity")

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	// Check Content-Disposition header first (most reliable)
	if filename := extractFilenameFromContentDisposition(resp.Header.Get("Content-Disposition")); filename != "" {
		return filename
	}

	// Check Content-Type header for file extension
	if filename := extractFilenameFromContentType(resp.Header.Get("Content-Type")); filename != "" {
		return filename
	}

	// If HEAD request didn't work, try GET request (some servers only provide headers on GET)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}

	// Same headers as HEAD request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,ru;q=0.8")
	req.Header.Set("Accept-Encoding", "identity")

	resp, err = client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	// Check Content-Disposition header from GET response
	if filename := extractFilenameFromContentDisposition(resp.Header.Get("Content-Disposition")); filename != "" {
		return filename
	}

	// Check Content-Type header from GET response
	if filename := extractFilenameFromContentType(resp.Header.Get("Content-Type")); filename != "" {
		return filename
	}

	return ""
}

// extractFilenameFromContentDisposition extracts filename from Content-Disposition header
func extractFilenameFromContentDisposition(contentDisposition string) string {
	if contentDisposition == "" {
		return ""
	}

	// Parse the Content-Disposition header
	_, params, err := mime.ParseMediaType(contentDisposition)
	if err != nil {
		return ""
	}

	// Look for filename parameter
	if filename, ok := params["filename"]; ok && filename != "" {
		// Clean the filename
		filename = strings.Trim(filename, `"'`)
		// Remove any path components for security
		return filepath.Base(filename)
	}

	// Look for filename* parameter (RFC 5987)
	if filename, ok := params["filename*"]; ok && filename != "" {
		// Parse RFC 5987 format: filename*=charset''encoded-filename
		if idx := strings.Index(filename, "''"); idx != -1 {
			filename = filename[idx+2:]
			// Clean the filename
			filename = strings.Trim(filename, `"'`)
			// Remove any path components for security
			return filepath.Base(filename)
		}
	}

	return ""
}

// extractFilenameFromContentType extracts a default filename based on Content-Type
func extractFilenameFromContentType(contentType string) string {
	if contentType == "" {
		return ""
	}

	// Parse the Content-Type header
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return ""
	}

	// Map common media types to file extensions
	extensions := map[string]string{
		"video/mp4":                    ".mp4",
		"video/webm":                   ".webm",
		"video/ogg":                    ".ogv",
		"video/avi":                    ".avi",
		"video/quicktime":              ".mov",
		"video/x-msvideo":              ".avi",
		"video/x-ms-wmv":               ".wmv",
		"video/x-flv":                  ".flv",
		"video/3gpp":                   ".3gp",
		"video/3gpp2":                  ".3g2",
		"video/x-matroska":             ".mkv",
		"image/jpeg":                   ".jpg",
		"image/png":                    ".png",
		"image/gif":                    ".gif",
		"image/webp":                   ".webp",
		"image/svg+xml":                ".svg",
		"audio/mpeg":                   ".mp3",
		"audio/ogg":                    ".ogg",
		"audio/wav":                    ".wav",
		"audio/webm":                   ".weba",
		"application/pdf":              ".pdf",
		"application/zip":              ".zip",
		"application/x-rar-compressed": ".rar",
	}

	if ext, ok := extensions[mediaType]; ok {
		return fmt.Sprintf("file%s", ext)
	}

	return ""
}

func DetectFileType(filename string) string {
	// Get file extension and convert to lowercase
	ext := strings.ToLower(filepath.Ext(filename))
	// Remove the dot from extension
	if len(ext) > 0 {
		ext = ext[1:]
	}

	// Common photo extensions
	photoExtensions := map[string]bool{
		"jpg": true, "jpeg": true, "png": true, "gif": true,
		"bmp": true, "tiff": true, "tif": true, "webp": true,
		"svg": true, "ico": true, "raw": true, "cr2": true,
		"nef": true, "arw": true, "dng": true, "orf": true,
		"rw2": true, "pef": true, "srw": true, "heic": true,
		"heif": true,
	}

	// Common video extensions
	videoExtensions := map[string]bool{
		"mp4": true, "avi": true, "mov": true, "wmv": true,
		"flv": true, "webm": true, "mkv": true, "m4v": true,
		"3gp": true, "ogv": true, "mpg": true, "mpeg": true,
		"ts": true, "vob": true, "asf": true, "rm": true,
		"rmvb": true, "f4v": true, "swf": true, "mts": true,
		"m2ts": true,
	}

	if photoExtensions[ext] {
		return "photo"
	} else if videoExtensions[ext] {
		return "video"
	}
	return "unknown"
}

// getFileSize retrieves the file size in bytes from a given URL.
func getFileSize(url string) (int64, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	// Standard approach for other URLs
	size1, _ := getStandardFileSize(client, url)
	// Special handling for OK CDN URLs
	size2, _ := getOKCDNFileSize(client, url)

	size := int64(math.Max(float64(size1), float64(size2)))

	return size, nil
}

func getOKCDNFileSize(client *http.Client, url string) (int64, error) {
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
				return bytes, nil
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
			return bytes, nil
		}
	}

	// Method 3: Read the actual content to determine size
	return readContentSize(resp.Body)
}

func getStandardFileSize(client *http.Client, url string) (int64, error) {
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

	return bytes, nil
}

func readContentSize(body io.ReadCloser) (int64, error) {
	const maxSize = 500 * 1024 * 1024 // 500MB limit
	bytesRead, err := io.Copy(io.Discard, io.LimitReader(body, maxSize))
	if err != nil {
		return 0, fmt.Errorf("failed to read content: %w", err)
	}

	if bytesRead == maxSize {
		return bytesRead, fmt.Errorf("file size exceeds 500MB limit (downloaded: %s)", ByteCountBinary(bytesRead))
	}

	return bytesRead, nil
}

func ByteCountBinary(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
