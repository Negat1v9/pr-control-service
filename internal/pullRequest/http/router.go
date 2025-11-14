package prhttp

import "net/http"

func PRRouter(h *PRHanler) http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("POST /create", h.Create)
	handler.HandleFunc("POST /merge", h.Merge)
	handler.HandleFunc("POST /reassign", h.Reassign)

	return handler
}
