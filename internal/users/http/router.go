package userhttp

import "net/http"

func UserRouter(h *UserHanler) http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("POST /setIsActive", h.SetIsActive)
	handler.HandleFunc("GET /getReview", h.GetReview)

	return handler
}
