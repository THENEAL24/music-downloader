package hitmo

import (
	"net/http"
	"sync"

	"github.com/THENEAL24/Music-Downloader/config"
)

type HitmoParser struct {
	Cfg       *config.AppConfig
	Track     *config.TrackConfig
	Client    *http.Client
	MP3Client *http.Client
	Lock      sync.Mutex
}

type Candidate struct {
	Artist string
	Title  string
	URL    string
}
