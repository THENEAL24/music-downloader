package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/THENEAL24/Music-Downloader/config"
	"github.com/THENEAL24/Music-Downloader/internal/api"
	"github.com/THENEAL24/Music-Downloader/internal/downloader"
	"github.com/THENEAL24/Music-Downloader/internal/sources/hitmo"
	"github.com/THENEAL24/Music-Downloader/internal/ytdlp"
	gozip "github.com/THENEAL24/Music-Downloader/internal/zip"
	"github.com/THENEAL24/Music-Downloader/pkg/logger"
)

func main() {
	cfgPath := "config/config.yaml"
	if p := os.Getenv("CONFIG_PATH"); p != "" {
		cfgPath = p
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	log := logger.New(os.Getenv("DEBUG") == "true")

	hitmoParser := hitmo.NewHitmoParser(cfg)

	var ytdlpDl *ytdlp.Downloader
	if cfg.Downloader.UseYtdlpFallback {
		ytdlpDl = ytdlp.NewDownloader(&cfg.Ytdlp)
	}

	svc := downloader.NewService(&cfg.Downloader, hitmoParser, ytdlpDl, log)

	archiver := gozip.NewArchiver(&cfg.Zip)

	mux := http.NewServeMux()
	handler := api.NewHandler(&cfg.API, svc, archiver, log)
	handler.RegisterRoutes(mux)

	chain := api.RequestLogger(log, api.RecoverPanic(log, mux))

	srv := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      chain,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("server started", "addr", cfg.Server.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	<-quit
	log.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("shutdown error", "err", err)
	}

	log.Info("server stopped")
}
