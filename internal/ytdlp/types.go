package ytdlp

import (
	"time"

	"github.com/THENEAL24/Music-Downloader/config"
	"github.com/THENEAL24/Music-Downloader/internal/domain"
)

type Downloader struct {
	cfg     *config.YtdlpConfig
	workers int
	jobs    chan job
	limiter *time.Ticker
}

type job struct {
	query domain.TrackQuery
	resCh chan domain.DownloadResult
}
