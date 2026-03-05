package domain

import (
	"strings"
)

func ParseTrackQuery(raw string) TrackQuery {
	q := TrackQuery{Raw: strings.TrimSpace(raw)}
	parts := strings.SplitN(q.Raw, " - ", 2)
	if len(parts) == 2 {
		q.Artist = strings.TrimSpace(parts[0])
		q.Title = strings.TrimSpace(parts[1])
	} else {
		q.Artist = strings.TrimSpace(q.Raw)
	}
	return q
}
