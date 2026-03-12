package hitmo

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/THENEAL24/Music-Downloader/internal/domain"
)

func (p *HitmoParser) searchURL(q string) string {
	return p.Cfg.HitmoBaseUrl + "/search/start/0?q=" + q
}

func (p *HitmoParser) search(q string) ([]Candidate, error) {
	req, err := http.NewRequest(http.MethodGet, p.searchURL(q), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9")
	req.Header.Set("Referer", p.Cfg.HitmoBaseUrl+"/")

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return p.parseHTML(body)
}

func (p *HitmoParser) SearchAndDownload(q domain.TrackQuery, delay time.Duration) (domain.DownloadedTrack, error) {
	attempts := uniqueStrings([]string{
		q.SearchQuery(),
		q.SearchQueryFirstArtist(),
		q.SearchQueryTitleOnly(),
	})

	lastErr := fmt.Errorf("не найдено на Hitmo")

	for i, query := range attempts {
		p.Lock.Lock()
		if i == 0 {
			time.Sleep(delay)
		} else {
			time.Sleep(time.Duration(float64(delay) * 0.4))
		}
		p.Lock.Unlock()

		candidates, err := p.search(query)
		if err != nil {
			lastErr = fmt.Errorf("поиск: %w", err)
			continue
		}

		match, ok := best(q, candidates)
		if !ok {
			lastErr = fmt.Errorf("нет результатов (попытка %d)", i+1)
			continue
		}

		data, err := p.fetchMP3(match.URL)
		if err != nil {
			lastErr = fmt.Errorf("загрузка файла: %w", err)
			continue
		}

		filename := sanitize(match.Artist + " - " + match.Title) + ".mp3"
		return domain.DownloadedTrack{
			Track: domain.Track{
				Query:       q,
				DisplayName: match.Artist + " - " + match.Title,
				DownloadURL: match.URL,
				Source:      "hitmo",
			},
			Data:     data,
			Filename: filename,
		}, nil
	}

	return domain.DownloadedTrack{}, lastErr
}

func uniqueStrings(ss []string) []string {
	seen := make(map[string]struct{}, len(ss))
	out := ss[:0:0]
	for _, s := range ss {
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
