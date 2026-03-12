package domain

import (
	"net/http"
	
	"github.com/THENEAL24/Music-Downloader/config"
)

func NewHTTPClient(cfg *config.ClientConfig) *http.Client {
	return &http.Client{Timeout: cfg.ClientTimeout}
}