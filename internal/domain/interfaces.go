package domain

import "time"

type TrackDownloader interface {
	SearchAndDownload(q TrackQuery, delay time.Duration) (DownloadedTrack, error)
}

type Archiver interface {
	Build(results []DownloadResult) ([]byte, error)
}
