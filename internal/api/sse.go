package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type msgStart struct {
	Type      string   `json:"type"`
	Total     int      `json:"total"`
	TrackList []string `json:"track_list"`
}

type msgTrack struct {
	Type   string `json:"type"`
	Done   int    `json:"done"`
	Total  int    `json:"total"`
	Query  string `json:"query"`
	OK     bool   `json:"ok"`
	Source string `json:"source,omitempty"`
	Error  string `json:"error,omitempty"`
}

type msgDone struct {
	Type  string  `json:"type"`
	OK    int     `json:"ok"`
	Fail  int     `json:"fail"`
	Pct   float64 `json:"pct"`
	JobID string  `json:"job_id"`
}

type msgError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}


func sseWrite(w http.ResponseWriter, v any) {
	data, _ := json.Marshal(v)
	fmt.Fprintf(w, "data: %s\n\n", data)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

func sseHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
