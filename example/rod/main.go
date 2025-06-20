package main

import (
	"fmt"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

func main() {
	chrome, found := launcher.LookPath()
	if !found {
		panic("not found")
	}

	u := launcher.New().Bin(chrome).Headless(true).MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()

	page := browser.
		MustPage("https://www.instagram.com/reels/DK_U7Vzh9rJ/").
		MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
			UserAgent: "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36",
		}).
		MustReload().
		MustWaitStable()

	page.Emulate(devices.Nexus5)

	elem, _ := page.Element("video")
	elem.MustWaitVisible()
	src, _ := elem.Attribute("src")
	fmt.Println("Video URL:", *src)

	// body := []byte(page.MustHTML())

	// re := regexp.MustCompile(`<video[^>]+src="([^"]+)"`)
	// match := re.FindSubmatch(body)

	// os.WriteFile("insta.html", body, 0644)
	// if len(match) >= 2 {
	// 	for m := range match {
	// 		fmt.Println("Match:", string(match[m]))
	// 	}
	// 	videoURL := string(match[1])
	// 	fmt.Println("Video URL:", videoURL)
	// } else {
	// 	fmt.Println("Video URL not found.")
	// }
}
