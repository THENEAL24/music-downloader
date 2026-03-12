package api

import (
	"log/slog"
	"net/http"

	"github.com/THENEAL24/Music-Downloader/config"
	"github.com/THENEAL24/Music-Downloader/internal/downloader"
	gozip "github.com/THENEAL24/Music-Downloader/internal/zip"
)

type Handler struct {
	cfg      *config.APIConfig
	svc      *downloader.Service
	archiver *gozip.Archiver
	store    *jobStore
	logger   *slog.Logger
}

func NewHandler(
	cfg *config.APIConfig,
	svc *downloader.Service,
	archiver *gozip.Archiver,
	logger *slog.Logger,
) *Handler {
	store := newJobStore(cfg.ResultTTL, cfg.EvictInterval)

	return &Handler{
		cfg:      cfg,
		svc:      svc,
		archiver: archiver,
		store:    store,
		logger:   logger,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /", newStaticHandler(h.cfg.StaticDir))

	mux.HandleFunc("POST /api/stream/txt",    h.streamTxt)
	mux.HandleFunc("POST /api/stream/yandex", h.streamYandex)
	mux.HandleFunc("GET /api/result/{id}",    h.getResult)
}
