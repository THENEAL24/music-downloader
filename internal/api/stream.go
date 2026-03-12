package api

import (
	"io"
	"net/http"
	"sync"

	"github.com/THENEAL24/Music-Downloader/internal/domain"
	"github.com/THENEAL24/Music-Downloader/internal/downloader"
	"github.com/THENEAL24/Music-Downloader/internal/sources/yandex"
)

func (h *Handler) streamTxt(w http.ResponseWriter, r *http.Request) {
	maxBytes := h.cfg.MaxUploadMB << 20
	if err := r.ParseMultipartForm(maxBytes); err != nil {
		http.Error(w, "bad form: "+err.Error(), http.StatusBadRequest)
		return
	}

	f, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file required", http.StatusBadRequest)
		return
	}
	defer f.Close()

	raw, err := io.ReadAll(io.LimitReader(f, maxBytes))
	if err != nil {
		http.Error(w, "read error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tracks := filterEmpty(splitLines(string(raw)))
	if len(tracks) == 0 {
		http.Error(w, "file is empty", http.StatusBadRequest)
		return
	}

	h.streamDownload(w, r, tracks, true)
}

func (h *Handler) streamYandex(w http.ResponseWriter, r *http.Request) {
	var body struct {
		URL   string `json:"url"`
		Token string `json:"token"`
	}
	if err := decodeJSON(r, &body); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	if body.URL == "" {
		http.Error(w, "url required", http.StatusBadRequest)
		return
	}

	sseHeaders(w)

	queries, err := yandex.FetchYandexPlaylist(body.URL, body.Token)
	if err != nil {
		sseWrite(w, msgError{Type: "error", Message: "❌ playlist: " + err.Error()})
		return
	}

	raw := make([]string, 0, len(queries))
	for _, q := range queries {
		raw = append(raw, q.Raw)
	}

	h.streamDownload(w, r, raw, false)
}

func (h *Handler) streamDownload(w http.ResponseWriter, r *http.Request, rawTracks []string, writeHeaders bool) {
	if writeHeaders {
		sseHeaders(w)
	}

	queries := make([]domain.TrackQuery, 0, len(rawTracks))
	for _, t := range rawTracks {
		queries = append(queries, domain.ParseTrackQuery(t))
	}

	total := len(queries)
	sseWrite(w, msgStart{Type: "start", Total: total, TrackList: rawTracks})

	var (
		mu      sync.Mutex
		counter int
	)

	ctx := r.Context()
	results := h.svc.Run(queries, func(res domain.DownloadResult) {
		select {
		case <-ctx.Done():
			return
		default:
		}

		mu.Lock()
		counter++
		done := counter
		mu.Unlock()

		sseWrite(w, msgTrack{
			Type:   "track",
			Done:   done,
			Total:  total,
			Query:  res.Query.Raw,
			OK:     res.Success,
			Source: res.Source,
			Error:  res.Error,
		})
	})

	if ctx.Err() != nil {
		return
	}

	data, err := h.archiver.Build(results)
	if err != nil {
		sseWrite(w, msgError{Type: "error", Message: "archive: " + err.Error()})
		return
	}

	jobID := h.store.save(data)
	stats := downloader.CalcStats(results, 0)

	var pct float64
	if total > 0 {
		pct = roundPct(float64(stats.OK) / float64(total) * 100)
	}

	sseWrite(w, msgDone{
		Type:  "done",
		OK:    stats.OK,
		Fail:  stats.Fail,
		Pct:   pct,
		JobID: jobID,
	})

	h.logger.Info("job finished", "job_id", jobID, "ok", stats.OK, "fail", stats.Fail)
}
