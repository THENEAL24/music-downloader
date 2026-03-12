package downloader

import (
	"log/slog"
	"sync"
	"time"

	"github.com/THENEAL24/Music-Downloader/config"
	"github.com/THENEAL24/Music-Downloader/internal/domain"
	"github.com/THENEAL24/Music-Downloader/internal/sources/hitmo"
	"github.com/THENEAL24/Music-Downloader/internal/ytdlp"
)

type Service struct {
	cfg    *config.DownloaderConfig
	hitmo  *hitmo.HitmoParser
	ytdlp  *ytdlp.Downloader
	logger *slog.Logger
}

func NewService(
	cfg *config.DownloaderConfig,
	hitmoParser *hitmo.HitmoParser,
	ytdlpDl *ytdlp.Downloader,
	logger *slog.Logger,
) *Service {
	return &Service{
		cfg:    cfg,
		hitmo:  hitmoParser,
		ytdlp:  ytdlpDl,
		logger: logger,
	}
}

type ProgressFunc func(result domain.DownloadResult)

func (s *Service) Run(tracks []domain.TrackQuery, onProgress ProgressFunc) []domain.DownloadResult {
	results := make([]domain.DownloadResult, len(tracks))
	jobs := make(chan indexedQuery, len(tracks))
	var wg sync.WaitGroup

	for i := 0; i < s.cfg.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ij := range jobs {
				r := s.downloadOne(ij.query)
				results[ij.index] = r
				if onProgress != nil {
					onProgress(r)
				}
			}
		}()
	}

	for i, q := range tracks {
		jobs <- indexedQuery{index: i, query: q}
	}
	close(jobs)
	wg.Wait()

	return results
}

type indexedQuery struct {
	index int
	query domain.TrackQuery
}

func (s *Service) downloadOne(q domain.TrackQuery) domain.DownloadResult {
	dt, err := s.hitmo.SearchAndDownload(q, s.cfg.HitmoDelay)
	if err == nil {
		return hitmoTrackToResult(dt)
	}

	s.logger.Debug("hitmo miss", "track", q.Raw, "err", err)

	if s.cfg.UseYtdlpFallback && s.ytdlp != nil {
		return s.ytdlp.Download(q)
	}

	return domain.DownloadResult{
		Query:   q,
		Success: false,
		Error:   err.Error(),
	}
}

func hitmoTrackToResult(dt domain.DownloadedTrack) domain.DownloadResult {
	return domain.DownloadResult{
		Query:    dt.Track.Query,
		Success:  true,
		Filename: dt.Filename,
		Content:  dt.Data,
		Source:   dt.Track.Source,
	}
}

type Stats struct {
	OK       int
	Fail     int
	BySource map[string]int
	Duration time.Duration
}

func CalcStats(results []domain.DownloadResult, dur time.Duration) Stats {
	s := Stats{BySource: make(map[string]int)}
	s.Duration = dur
	for _, r := range results {
		if r.Success {
			s.OK++
			s.BySource[r.Source]++
		} else {
			s.Fail++
		}
	}
	return s
}
