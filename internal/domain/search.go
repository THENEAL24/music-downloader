package domain

import (
	"regexp"
	"strings"
)

func (q TrackQuery) SearchQuery() string {
	return strings.ReplaceAll(strings.ToLower(strings.TrimSpace(q.Raw)), " ", "+")
}

func (q TrackQuery) SearchQueryFirstArtist() string {
	re := regexp.MustCompile(`[,&]`)
	first := strings.TrimSpace(re.Split(q.Artist, -1)[0])
	return strings.ReplaceAll(strings.ToLower(strings.TrimSpace(first+" "+q.Title)), " ", "+")
}

func (q TrackQuery) SearchQueryTitleOnly() string {
	title := q.Title
	if title == "" {
		title = q.Raw
	}
	return strings.ReplaceAll(strings.ToLower(strings.TrimSpace(title)), " ", "+")
}