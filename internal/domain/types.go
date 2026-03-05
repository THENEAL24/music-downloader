package domain

type TrackQuery struct {
	Raw    string
	Artist string
	Title  string
}

type Track struct {
	Query       TrackQuery
	DisplayName string
	DownloadURL string
	Source      string
}

type DownloadedTrack struct {
	Track    Track
	Data     []byte
	Filename string
}

type FailedTrack struct {
	Query  TrackQuery
	Reason string
}

type DownloadSession struct {
	Downloaded []DownloadedTrack
	Failed     []FailedTrack
}