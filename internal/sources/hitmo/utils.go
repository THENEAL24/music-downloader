package hitmo

import (
	"regexp"
	"strings"
	"bytes"

	"github.com/THENEAL24/Music-Downloader/internal/domain"
)

func sanitize(name string) string {
	re := regexp.MustCompile(`[<>:"/\\|?*&,]`)
	return strings.TrimSpace(re.ReplaceAllString(name, ""))
}

func score(q domain.TrackQuery, c Candidate) int {
	s := 0
	qt := strings.ToLower(q.Title)
	qa := strings.ToLower(strings.TrimSpace(regexp.MustCompile(`[,&]`).Split(q.Artist, -1)[0]))
	ct := strings.ToLower(c.Title)
	ca := strings.ToLower(c.Artist)

	if qt != "" && strings.Contains(ct, qt) {
		s += 2
	} else if qt != "" {
		for _, w := range strings.Fields(qt) {
			if len(w) > 3 && strings.Contains(ct, w) {
				s++
				break
			}
		}
	}

	if qa != "" && strings.Contains(ca, qa) {
		s++
	}

	return s
}

func best(q domain.TrackQuery, candidates []Candidate) (Candidate, bool) {
	if len(candidates) == 0 {
		return Candidate{}, false
	}
	top := candidates[0]
	topScore := score(q, top)
	for _, c := range candidates[1:] {
		if sc := score(q, c); sc > topScore {
			topScore = sc
			top = c
		}
	}
	return top, true
}

func isMP3(data []byte) bool {
	if len(data) < 3 {
		return false
	}

	if bytes.HasPrefix(data, []byte("ID3")) {
		return true
	}

	if data[0] == 0xFF && (data[1]&0xE0) == 0xE0 {
		return true
	}

	return false
}