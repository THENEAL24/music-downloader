package hitmo

import (
	"fmt"
	"io"
	"net/http"
)

func (p *HitmoParser) fetchMP3(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Referer", p.Cfg.HitmoBaseUrl+"/")

	resp, err := p.MP3Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if !isMP3(data) {
		return nil, fmt.Errorf("downloaded file is not mp3")
	}

	return data, nil
}
