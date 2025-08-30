package common

import (
	// standard
	"fmt"
	"io"
	"net/http"
	"os"
)

func UrlExists(url string) bool {
	r, e := http.Head(url)
	return e == nil && r.StatusCode == 200
}

func DownloadFile(URL, fileName string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Received non 200 response code: %d", resp.StatusCode)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
