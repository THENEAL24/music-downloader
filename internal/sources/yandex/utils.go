package yandex

import (
	"fmt"
	"strings"

	"github.com/THENEAL24/Music-Downloader/internal/domain"
)

func artistNames(artists []yandexArtist) string {
	names := make([]string, 0, len(artists))
	for _, a := range artists {
		names = append(names, a.Name)
	}
	return strings.Join(names, ", ")
}

func trackToQuery(artists []yandexArtist, title string) domain.TrackQuery {
	raw := fmt.Sprintf("%s - %s", artistNames(artists), title)
	return domain.TrackQuery{
		Raw:    raw,
		Artist: artistNames(artists),
		Title:  title,
	}
}
