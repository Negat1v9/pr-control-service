package prhttp

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Negat1v9/pr-review-service/internal/models"
	prservice "github.com/Negat1v9/pr-review-service/internal/pullRequest/service"
	"github.com/Negat1v9/pr-review-service/pkg/logger"
	"github.com/Negat1v9/pr-review-service/pkg/utils"
)

type PRHanler struct {
	log     *logger.Logger
	service *prservice.PRService
}

func NewPRHanlder(log *logger.Logger, service *prservice.PRService) *PRHanler {
	return &PRHanler{
		log:     log,
		service: service,
	}

}

func (h *PRHanler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()
	var req models.CreatePullRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrResponse(w, utils.NewBadRequestError("invalid request body", nil))
		return
	}

	newPR, err := h.service.CreatePR(ctx, &req)
	if err != nil {
		h.log.Errorf("failed to create pull request: %v", err)
		utils.WriteErrResponse(w, err)
		return
	}

	utils.WriteJsonResponse(w, http.StatusCreated, "pr", newPR)
}
func (h *PRHanler) Merge(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()
	var req models.MergePullRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrResponse(w, utils.NewBadRequestError("invalid request body", nil))
		return
	}

	mergedPR, err := h.service.MergePR(ctx, req.ID)
	if err != nil {
		h.log.Errorf("failed to merge pull request: %v", err)
		utils.WriteErrResponse(w, err)
		return
	}

	utils.WriteJsonResponse(w, http.StatusOK, "pr", mergedPR)
}

func (h *PRHanler) Reassign(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()
	var req models.ReassignPullRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrResponse(w, utils.NewBadRequestError("invalid request body", nil))
		return
	}

	updatedPR, err := h.service.ReassignPR(ctx, req.ID, req.OldReviewerID)
	if err != nil {
		h.log.Errorf("failed to reassign pull request: %v", err)
		utils.WriteErrResponse(w, err)
		return
	}

	utils.WriteJsonResponse(w, 200, "", updatedPR)
}

func (h *PRHanler) Statistics(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	pullRequestsQuantiReviewers, err := h.service.Statistics(ctx)
	if err != nil {
		h.log.Errorf("failed to reassign pull request: %v", err)
		utils.WriteErrResponse(w, err)
	}

	utils.WriteJsonResponse(w, 200, "stat", pullRequestsQuantiReviewers)
}
