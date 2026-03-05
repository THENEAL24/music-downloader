package yandex

type yandexArtist struct {
	Name string `json:"name"`
}

type yandexTrack struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Artists []yandexArtist `json:"artists"`
}

type yandexPlaylistTrack struct {
	Track yandexTrack `json:"track"`
}

type yandexPlaylistResp struct {
	Result struct {
		TrackCount int
		Tracks []yandexPlaylistTrack
	} `json:"result"`
}

type yandexAlbumVolume []yandexTrack

type yandexAlbumResp struct {
	Result struct {
		Volumes []yandexAlbumVolume `json:"volumes"`
	} `json:"result"`
}