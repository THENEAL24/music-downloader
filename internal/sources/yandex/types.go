package yandex

import "encoding/json"

type yandexArtist struct {
	Name string `json:"name"`
}

// yandexTrack — id может прийти как number или string, используем json.Number.
type yandexTrack struct {
	ID      json.Number    `json:"id"`
	Title   string         `json:"title"`
	Artists []yandexArtist `json:"artists"`
}

type yandexPlaylistTrack struct {
	Track yandexTrack `json:"track"`
}

type yandexPlaylistResp struct {
	Result struct {
		TrackCount int                   `json:"trackCount"`
		Tracks     []yandexPlaylistTrack `json:"tracks"`
	} `json:"result"`
}

type yandexAlbumVolume []yandexTrack

type yandexAlbumResp struct {
	Result struct {
		Volumes []yandexAlbumVolume `json:"volumes"`
	} `json:"result"`
}
