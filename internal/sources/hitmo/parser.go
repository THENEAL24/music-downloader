package hitmo

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/THENEAL24/Music-Downloader/config"
)

func NewHitmoParser(cfg *config.Config) *HitmoParser {
	search, mp3 := newHTTPClientMP3(&cfg.Client)
	return &HitmoParser{
		Cfg:       &cfg.App,
		Track:     &cfg.Track,
		Client:    search,
		MP3Client: mp3,
	}
}

func (p *HitmoParser) parseHTML(body []byte) ([]Candidate, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	var out []Candidate
	doc.Find(p.Track.ItemsCss).Each(func(_ int, s *goquery.Selection) {
		artist := strings.TrimSpace(s.Find(p.Track.ArtistCss).Text())
		title := strings.TrimSpace(s.Find(p.Track.TitleCss).Text())
		url, exists := s.Find(p.Track.DlBtnCss).Attr("href")
		if artist != "" && title != "" && exists && url != "" {
			out = append(out, Candidate{Artist: artist, Title: title, URL: url})
		}
	})
	return out, nil
}
