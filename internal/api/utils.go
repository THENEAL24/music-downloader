package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

func splitLines(s string) []string {
	lines := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	result := make([]string, 0, len(lines))
	for _, l := range lines {
		result = append(result, strings.TrimRight(l, "\r"))
	}
	return result
}

func filterEmpty(lines []string) []string {
	out := lines[:0]
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			out = append(out, l)
		}
	}
	return out
}

func roundPct(f float64) float64 {
	return float64(int(f*10+0.5)) / 10
}

func decodeJSON(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}
