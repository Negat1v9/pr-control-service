package teamservice

import (
	"context"
	"database/sql"

	"github.com/Negat1v9/pr-review-service/internal/models"
	"github.com/Negat1v9/pr-review-service/internal/store"
	"github.com/Negat1v9/pr-review-service/pkg/utils"
	"github.com/jmoiron/sqlx"
)

type TeamService struct {
	store store.Store
}

func NewTeamService(store store.Store) *TeamService {
	return &TeamService{
		store: store,
	}
}

func (s *TeamService) AddTeam(ctx context.Context, newTeam *models.Team) (*models.Team, error) {

	_, err := s.store.TeamRepo().GetTeamWithMembers(ctx, s.store.DB(), newTeam.TeamName)
	if err == nil {
		return nil, utils.NewError(400, utils.ErrTeamExists, "team_name already exists", nil)
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	err = s.store.DoTx(ctx, func(ctx context.Context, exec sqlx.ExtContext) error {

		if err := s.store.TeamRepo().CreateTeam(ctx, exec, newTeam.TeamName); err != nil {
			return err
		}

		if err := s.store.UserRepo().CreateManyUsers(ctx, exec, newTeam.TeamName, newTeam.Members); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	createdTeam, err := s.store.TeamRepo().GetTeamWithMembers(ctx, s.store.DB(), newTeam.TeamName)
	if err != nil {
		return nil, err
	}

	return createdTeam, nil
}

func (s *TeamService) GetTeam(ctx context.Context, teamName string) (*models.Team, error) {
	team, err := s.store.TeamRepo().GetTeamWithMembers(ctx, s.store.DB(), teamName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NewNotFoundError("resource not found", nil)
		}
		return nil, err
	}

	return team, nil
}
