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

	u := launcher.New().Bin(chrome).Headless(false).MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()

	page := browser.
		MustPage("https://www.instagram.com/stories/msgr.ig/").
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
}
