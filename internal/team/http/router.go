package teamhttp

import "net/http"

func TeamRouter(h *TeamHanler) http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("POST /add", h.Add)
	handler.HandleFunc("GET /get", h.Get)

	return handler
}
