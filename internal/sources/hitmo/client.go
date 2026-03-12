package hitmo

import (
	"net/http"

	"github.com/THENEAL24/Music-Downloader/config"
)

func newHTTPClientMP3(cfg *config.ClientConfig) (searchClient *http.Client, mp3Client *http.Client) {
	searchClient = &http.Client{Timeout: cfg.ClientTimeout}
	mp3Client = &http.Client{Timeout: cfg.MP3Timeout}
	return
}
