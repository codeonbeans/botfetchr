package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	link := `https://vkvd296.okcdn.ru/?expires=1750928104071&srcIp=58.187.246.12&pr=40&srcAg=CHROME&ms=185.226.53.133&type=4&subId=7268260383481&sig=Ro6ulqoxn5c&ct=0&urls=45.136.21.154&clientType=13&appId=512000384397&zs=65&id=7266165197561`

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		panic(err)
	}

	// Mimic a real browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	// Read the entire response into memory
	videoData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	os.WriteFile("test.mp4", videoData, 0644)
}
