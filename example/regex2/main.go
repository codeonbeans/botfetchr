package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Print(isReel("https://www.instagram.com/reels/DLAi8xlySmh/"))

}

func isReel(url string) bool {
	return strings.Contains(url, "/reel/") || strings.Contains(url, "/reels/")
}
