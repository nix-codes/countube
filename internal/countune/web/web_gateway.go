package web

import (
	"fmt"
	"image"
	"image/png"
	"net/http"
	"os"
)

type statusError int

func (s statusError) Error() string {
	return fmt.Sprintf("HTTP status %d", int(s))
}

func FetchPic(picNum int) (image.Image, error) {
	url := getCountunePicUrl(picNum)

	return downloadPngImage(url)
}

func getCountunePicUrl(picNum int) string {
	return fmt.Sprintf(CountunePicUrlFmt, picNum)
}

func downloadPngImage(url string) (image.Image, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		os.Exit(1)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to fetch image (no HTTP status):", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code: ", resp.StatusCode)
		os.Exit(1)
	}

	img, err := png.Decode(resp.Body)
	if err != nil {
		fmt.Println("Failed to decode PNG:", err)
		os.Exit(1)
	}

	return img, nil
}
