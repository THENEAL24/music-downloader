package zip

import (
	"archive/zip"
	"bytes"
	"fmt"
	"strings"

	"github.com/THENEAL24/Music-Downloader/config"
	"github.com/THENEAL24/Music-Downloader/internal/domain"
)

type Archiver struct {
	level int
}

func NewArchiver(cfg *config.ZipConfig) *Archiver {
	return &Archiver{level: cfg.CompressionLevel}
}

func (a *Archiver) Build(results []domain.DownloadResult) ([]byte, error) {
	buf := &bytes.Buffer{}
	w := zip.NewWriter(buf)

	var failedLines []string

	for _, r := range results {
		if !r.Success {
			failedLines = append(failedLines, fmt.Sprintf("%s  ←  %s", r.Query.Raw, r.Error))
			continue
		}

		fh := &zip.FileHeader{
			Name:   r.Filename,
			Method: zip.Deflate,
		}

		fw, err := w.CreateHeader(fh)
		if err != nil {
			return nil, fmt.Errorf("create entry %q: %w", r.Filename, err)
		}

		if _, err := fw.Write(r.Content); err != nil {
			return nil, fmt.Errorf("write entry %q: %w", r.Filename, err)
		}
	}

	if len(failedLines) > 0 {
		fw, err := w.Create("failed_tracks.txt")
		if err != nil {
			return nil, fmt.Errorf("create failed_tracks.txt: %w", err)
		}
		if _, err := fmt.Fprint(fw, strings.Join(failedLines, "\n")); err != nil {
			return nil, fmt.Errorf("write failed_tracks.txt: %w", err)
		}
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("zip close: %w", err)
	}

	return buf.Bytes(), nil
}
