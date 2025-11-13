package teamhttp

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Negat1v9/pr-review-service/internal/models"
	teamservice "github.com/Negat1v9/pr-review-service/internal/team/service"
	"github.com/Negat1v9/pr-review-service/pkg/logger"
	"github.com/Negat1v9/pr-review-service/pkg/utils"
)

type TeamHanler struct {
	log     *logger.Logger
	service *teamservice.TeamService
}

func NewTeamHanlder(log *logger.Logger, service *teamservice.TeamService) *TeamHanler {
	return &TeamHanler{
		log:     log,
		service: service,
	}

}

func (h *TeamHanler) Add(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	var newTeam models.Team
	if err := json.NewDecoder(r.Body).Decode(&newTeam); err != nil {
		utils.WriteErrResponse(w, utils.NewBadRequestError("invalid request body", nil))
		return
	}

	createdTeam, err := h.service.AddTeam(ctx, &newTeam)
	if err != nil {
		h.log.Errorf("failed to create team: %v", err)
		utils.WriteErrResponse(w, err)
		return
	}

	utils.WriteJsonResponse(w, http.StatusCreated, "team", createdTeam)
}
func (h *TeamHanler) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		utils.WriteErrResponse(w, utils.NewNotFoundError("resource not found", nil))
		return
	}

	team, err := h.service.GetTeam(ctx, teamName)
	if err != nil {
		h.log.Errorf("failed to get team: %v", err)
		utils.WriteErrResponse(w, err)
		return
	}

	utils.WriteJsonResponse(w, http.StatusOK, "", team)
}
