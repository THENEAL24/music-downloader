package api

import "net/http"

func newStaticHandler(dir string) http.Handler {
	return http.FileServer(http.Dir(dir))
}
