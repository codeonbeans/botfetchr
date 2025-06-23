package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	// text := `t":0.000000,"url144":"https:\/\/vk6-7.vkuser.net\/?srcIp=58.187.246.12&pr=40&expires=1751116459527&srcAg=CHROME&fromCache=1&ms=95.142.206.166&type=4&subId=8151226255876&sig=UODdnK5TZkU&ct=0&urls=45.136.22.165%3B185.226.52.209&clientType=13&appId=512000384397&zs=43&id=8150980758020"`

	// re := regexp.MustCompile(`"url\d+":"([^"]+)"`)
	// match := re.FindStringSubmatch(text)
	// if len(match) > 1 {
	// 	fmt.Println("Found URL:", match[1])
	// } else {
	// 	fmt.Println("No match found")
	// }

	text := `"https:\\/\\/vkvd296.okcdn.ru\\/?srcIp=58.187.246.12&pr=40&expires=1751129922986&srcAg=CHROME&fromCache=1&ms=185.226.53.133&type=5&subId=7268260383481&sig=RSu0I7SJ5BE&ct=0&urls=45.136.21.154&clientType=13&appId=512000384397&zs=65&id=7266165197561"`
	var idk string
	json.Unmarshal([]byte(text), &idk)
	idk, _ = FixEscapedURL(idk)
	fmt.Println("NIGGA KYS", idk)

	resp, err := http.Get(idk)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

}

func UnmarshalURL(marshalledURL string) (string, error) {
	var result string
	err := json.Unmarshal([]byte(marshalledURL), &result)
	return result, err
}

func FixEscapedURL(escapedURL string) (string, error) {
	// Replace escaped forward slashes
	fixedURL := strings.ReplaceAll(escapedURL, `\/`, `/`)

	// Validate the URL
	_, err := url.Parse(fixedURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL after fixing: %v", err)
	}

	return fixedURL, nil
}
