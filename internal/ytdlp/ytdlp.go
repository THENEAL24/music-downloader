package ytdlp

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"time"

	"github.com/THENEAL24/Music-Downloader/config"
	"github.com/THENEAL24/Music-Downloader/internal/domain"
)

func NewDownloader(cfg *config.YtdlpConfig) *Downloader {
	d := &Downloader{
		cfg:     cfg,
		jobs:    make(chan job, cfg.Workers*2),
		limiter: time.NewTicker(time.Second / time.Duration(cfg.RPS)),
	}

	for i := 0; i < cfg.Workers; i++ {
		go d.worker()
	}

	return d
}

func (d *Downloader) Download(query domain.TrackQuery) domain.DownloadResult {
	resCh := make(chan domain.DownloadResult, 1)

	d.jobs <- job{
		query: query,
		resCh: resCh,
	}

	return <-resCh
}

func (d *Downloader) worker() {
	for j := range d.jobs {
		<-d.limiter.C
		res := d.downloadWithRetry(j.query)
		j.resCh <- res
	}
}

func (d *Downloader) downloadWithRetry(query domain.TrackQuery) domain.DownloadResult {
	attempts := uniqueAttempts(query)

	var lastErr error

	for _, q := range attempts {
		res, err := d.downloadOnce(q)
		if err == nil {
			res.Query = query
			return res
		}
		lastErr = err
	}

	return domain.DownloadResult{
		Query:   query,
		Success: false,
		Error:   lastErr.Error(),
	}
}

func (d *Downloader) downloadOnce(search string) (domain.DownloadResult, error) {
	maxSize := int64(d.cfg.MaxAudioSizeMB) << 20

	ctx, cancel := context.WithTimeout(context.Background(), d.cfg.Timeout)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		"yt-dlp",
		"-f", "bestaudio/best",
		"-x",
		"--audio-format", d.cfg.AudioFormat,
		"--audio-quality", d.cfg.AudioQuality,
		"-o", "-",
		"--no-playlist",
		"ytsearch1:"+search,
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return domain.DownloadResult{}, err
	}

	if err := cmd.Start(); err != nil {
		return domain.DownloadResult{}, err
	}

	data, err := io.ReadAll(io.LimitReader(stdout, maxSize))
	if err != nil {
		return domain.DownloadResult{}, err
	}

	if err := cmd.Wait(); err != nil {
		return domain.DownloadResult{}, err
	}

	if len(data) == 0 {
		return domain.DownloadResult{}, fmt.Errorf("empty output")
	}

	return domain.DownloadResult{
		Success:  true,
		Filename: sanitize(search) + "." + d.cfg.AudioFormat,
		Content:  data,
		Source:   "ytdlp",
	}, nil
}

func uniqueAttempts(q domain.TrackQuery) []string {
	candidates := []string{
		q.Raw,
		q.Artist + " " + q.Title,
		q.Title,
	}
	seen := make(map[string]struct{}, len(candidates))
	out := make([]string, 0, len(candidates))
	for _, s := range candidates {
		if s == "" {
			continue
		}
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, s)
		}
	}
	return out
}

func sanitize(name string) string {
	re := regexp.MustCompile(`[<>:"/\\|?*&,]`)
	return re.ReplaceAllString(name, "")
}
