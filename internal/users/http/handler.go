package userhttp

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Negat1v9/pr-review-service/internal/models"
	userservice "github.com/Negat1v9/pr-review-service/internal/users/service"
	"github.com/Negat1v9/pr-review-service/pkg/logger"
	"github.com/Negat1v9/pr-review-service/pkg/utils"
)

type UserHanler struct {
	log     *logger.Logger
	service *userservice.UserService
}

func NewUserHandler(log *logger.Logger, service *userservice.UserService) *UserHanler {
	return &UserHanler{
		log:     log,
		service: service,
	}
}

func (h *UserHanler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()
	var req models.SetUserActiveStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrResponse(w, utils.NewBadRequestError("invalid request body", nil))
		return
	}

	updatedUser, err := h.service.SetUserActiveStatus(ctx, req.UserID, req.IsActive)
	if err != nil {
		h.log.Errorf("failed to set user active status: %v", err)
		utils.WriteErrResponse(w, err)
		return
	}

	utils.WriteJsonResponse(w, http.StatusOK, "user", updatedUser)
}

func (h *UserHanler) GetReview(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		utils.WriteErrResponse(w, utils.NewNotFoundError("resource not found", nil))
		return
	}

	userReviews, err := h.service.GetReview(ctx, userID)
	if err != nil {
		h.log.Errorf("failed to get user reviews: %v", err)
		utils.WriteErrResponse(w, err)
		return
	}

	utils.WriteJsonResponse(w, http.StatusOK, "", userReviews)
}
