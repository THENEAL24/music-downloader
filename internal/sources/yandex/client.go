package yandex

import (
	"fmt"
	"io"
	"net/http"
)

func doRequest(client *http.Client, url, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Yandex_Music_Client", "WindowsPhone/3.20")
	if token != "" {
		req.Header.Set("Authorization", "OAuth "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d для %s", resp.StatusCode, url)
	}
	return io.ReadAll(resp.Body)
}

