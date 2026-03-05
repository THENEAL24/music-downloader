package yandex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"github.com/THENEAL24/Music-Downloader/internal/domain"
)

func FetchYandexPlaylist(rawURL, token string) ([]domain.TrackQuery, error) {
	client := newHTTPClient(cfg.Client)

	rePlaylist := regexp.MustCompile(`users/([^/]+)/playlists/(\d+)`)
	if m := rePlaylist.FindStringSubmatch(rawURL); m != nil {
		return fetchUserPlaylist(client, m[1], m[2], token)
	}

	reAlbum := regexp.MustCompile(`album/(\d+)`)
	if m := reAlbum.FindStringSubmatch(rawURL); m != nil {
		return fetchAlbum(client, m[1], token)
	}

	return nil, fmt.Errorf(
		"неподдерживаемый URL\nОжидается:\n"+
			"  https://music.yandex.ru/users/{login}/playlists/{id}\n"+
			"  https://music.yandex.ru/album/{id}",
	)
}

func fetchUserPlaylist(client *http.Client, login, kind, token string) ([]domain.TrackQuery, error) {
	url := fmt.Sprintf(
		"https://api.music.yandex.net/users/%s/playlists/%s?rich-tracks=true",
		login, kind,
	)

	body, err := doRequest(client, url, token)
	if err != nil {
		return nil, fmt.Errorf("playlist request: %w", err)
	}

	var resp yandexPlaylistResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("playlist parse: %w", err)
	}

	tracks := make([]domain.TrackQuery, 0, len(resp.Result.Tracks))
	for _, pt := range resp.Result.Tracks {
		t := pt.Track
		tracks = append(tracks, trackToQuery(t.Artists, t.Title))
	}
	return tracks, nil
}

func fetchAlbum(client *http.Client, albumID, token string) ([]domain.TrackQuery, error) {
	url := fmt.Sprintf("https://api.music.yandex.net/albums/%s/with-tracks", albumID)
	body, err := doRequest(client, url, token)
	if err != nil {
		return nil, fmt.Errorf("album request: %w", err)
	}

	var resp yandexAlbumResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("album parse: %w", err)
	}

	totalTracks := 0
	for _, vol := range resp.Result.Volumes {
		totalTracks += len(vol)
	}

	tracks := make([]domain.TrackQuery, 0, totalTracks)
	for _, vol := range resp.Result.Volumes {
		for _, t := range vol {
			tracks = append(tracks, trackToQuery(t.Artists, t.Title))
		}
	}
	return tracks, nil
}
