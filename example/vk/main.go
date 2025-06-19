package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Instagram video downloader
func downloadInstagramVideo(url string) error {
	// Extract shortcode from URL
	shortcode, err := extractInstagramShortcode(url)
	if err != nil {
		return fmt.Errorf("failed to extract shortcode: %v", err)
	}

	// Get video info using Instagram's API endpoint
	apiURL := fmt.Sprintf("https://www.instagram.com/p/%s/?__a=1&__d=dis", shortcode)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers to mimic browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch Instagram data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Instagram API returned status: %d", resp.StatusCode)
	}

	// Parse response
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	// Extract video URL from response
	videoURL, err := extractInstagramVideoURL(data)
	if err != nil {
		return fmt.Errorf("failed to extract video URL: %v", err)
	}

	// Download the video
	return downloadFile(videoURL, fmt.Sprintf("instagram_%s.mp4", shortcode))
}

// VK video downloader
func downloadVKVideo(url string) error {
	// Extract video ID from VK URL
	videoID, err := extractVKVideoID(url)
	if err != nil {
		return fmt.Errorf("failed to extract VK video ID: %v", err)
	}

	// Note: VK requires API access token for video downloads
	// This is a simplified example - you'll need to implement OAuth2 flow
	// and get proper access token from VK API

	accessToken := "YOUR_VK_ACCESS_TOKEN" // Replace with actual token
	apiURL := fmt.Sprintf("https://api.vk.com/method/video.get?videos=%s&access_token=%s&v=5.131", videoID, accessToken)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(apiURL)
	if err != nil {
		return fmt.Errorf("failed to fetch VK data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("VK API returned status: %d", resp.StatusCode)
	}

	// Parse VK API response
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("failed to decode VK response: %v", err)
	}

	// Extract video URL from VK response
	videoURL, err := extractVKVideoURL(data)
	if err != nil {
		return fmt.Errorf("failed to extract VK video URL: %v", err)
	}

	// Download the video
	return downloadFile(videoURL, fmt.Sprintf("vk_%s.mp4", strings.ReplaceAll(videoID, "_", "-")))
}

// Helper function to extract Instagram shortcode from URL
func extractInstagramShortcode(url string) (string, error) {
	// Match patterns like: /p/ABC123/, /reel/ABC123/, /tv/ABC123/
	re := regexp.MustCompile(`(?:instagram\.com\/(?:p|reel|tv)\/([A-Za-z0-9_-]+))`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("invalid Instagram URL format")
	}
	return matches[1], nil
}

// Helper function to extract VK video ID from URL
func extractVKVideoID(url string) (string, error) {
	// Match patterns like: video-123456_789012 or video123456_789012
	re := regexp.MustCompile(`video(-?\d+_\d+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("invalid VK video URL format")
	}
	return matches[1], nil
}

// Helper function to extract video URL from Instagram API response
func extractInstagramVideoURL(data map[string]interface{}) (string, error) {
	// Navigate through the nested JSON structure
	items, ok := data["graphql"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected response format")
	}

	shortcodeMedia, ok := items["shortcode_media"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("shortcode_media not found")
	}

	// Check if it's a video
	isVideo, ok := shortcodeMedia["is_video"].(bool)
	if !ok || !isVideo {
		return "", fmt.Errorf("media is not a video")
	}

	// Get video URL
	videoURL, ok := shortcodeMedia["video_url"].(string)
	if !ok {
		return "", fmt.Errorf("video_url not found")
	}

	return videoURL, nil
}

// Helper function to extract video URL from VK API response
func extractVKVideoURL(data map[string]interface{}) (string, error) {
	response, ok := data["response"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid VK response format")
	}

	items, ok := response["items"].([]interface{})
	if !ok || len(items) == 0 {
		return "", fmt.Errorf("no video items found")
	}

	video, ok := items[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid video item format")
	}

	files, ok := video["files"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("video files not found")
	}

	// Try to get the highest quality available
	for _, quality := range []string{"mp4_1080", "mp4_720", "mp4_480", "mp4_360", "mp4_240"} {
		if url, exists := files[quality].(string); exists {
			return url, nil
		}
	}

	return "", fmt.Errorf("no suitable video quality found")
}

// Helper function to download file from URL
func downloadFile(url, filename string) error {
	client := &http.Client{
		Timeout: 5 * time.Minute,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create download request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Create downloads directory if it doesn't exist
	if err := os.MkdirAll("downloads", 0755); err != nil {
		return fmt.Errorf("failed to create downloads directory: %v", err)
	}

	// Create the file
	filePath := filepath.Join("downloads", filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Copy data to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	fmt.Printf("Successfully downloaded: %s\n", filePath)
	return nil
}

func main() {
	// Example usage
	instagramURL := "https://www.instagram.com/p/DLFX_PfRAv_/"
	vkURL := "https://vk.com/video-123456_789012"

	fmt.Println("Downloading Instagram video...")
	if err := downloadInstagramVideo(instagramURL); err != nil {
		fmt.Printf("Instagram download failed: %v\n", err)
	}

	fmt.Println("Downloading VK video...")
	if err := downloadVKVideo(vkURL); err != nil {
		fmt.Printf("VK download failed: %v\n", err)
	}
}

// Additional utility functions for enhanced functionality

// Function to get video metadata
func getVideoMetadata(url string) (map[string]interface{}, error) {
	if strings.Contains(url, "instagram.com") {
		return getInstagramMetadata(url)
	} else if strings.Contains(url, "vk.com") {
		return getVKMetadata(url)
	}
	return nil, fmt.Errorf("unsupported platform")
}

func getInstagramMetadata(url string) (map[string]interface{}, error) {
	shortcode, err := extractInstagramShortcode(url)
	if err != nil {
		return nil, err
	}

	apiURL := fmt.Sprintf("https://www.instagram.com/p/%s/?__a=1&__d=dis", shortcode)

	client := &http.Client{Timeout: 30 * time.Second}
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)

	return data, nil
}

func getVKMetadata(url string) (map[string]interface{}, error) {
	videoID, err := extractVKVideoID(url)
	if err != nil {
		return nil, err
	}

	accessToken := "YOUR_VK_ACCESS_TOKEN"
	apiURL := fmt.Sprintf("https://api.vk.com/method/video.get?videos=%s&access_token=%s&v=5.131", videoID, accessToken)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)

	return data, nil
}
