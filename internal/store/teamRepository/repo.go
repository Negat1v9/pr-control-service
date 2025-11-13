package teamrepository

import (
	"context"
	"database/sql"

	"github.com/Negat1v9/pr-review-service/internal/models"
	"github.com/jmoiron/sqlx"
)

type teamRepositiry struct{}

func NewTeamRepositiry() *teamRepositiry {
	return &teamRepositiry{}
}

func (r *teamRepositiry) CreateTeam(ctx context.Context, exec sqlx.ExtContext, teamName string) error {
	if _, err := exec.ExecContext(ctx, createTeamQuery, teamName); err != nil {
		return err
	}
	return nil
}

func (r *teamRepositiry) GetTeamWithMembers(ctx context.Context, exec sqlx.ExtContext, teamName string) (*models.Team, error) {

	rows, err := exec.QueryxContext(ctx, getTeamWithUsersByNameQuery, teamName)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var team = models.Team{}

	for rows.Next() {
		var member models.User
		var scannedTeamName string

		if err = rows.Scan(&scannedTeamName, &member.UserID, &member.Username, &member.IsActive); err != nil {
			return nil, err
		}

		if team.TeamName == "" {
			team.TeamName = scannedTeamName
		}

		team.Members = append(team.Members, member)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if team.TeamName == "" {
		return nil, sql.ErrNoRows
	}

	return &team, nil
}

func (r *teamRepositiry) GetUsersIDFromUserTeam(ctx context.Context, exec sqlx.ExtContext, userID string) ([]string, error) {
	rows, err := exec.QueryxContext(ctx, getActiveUserIDFromUserTeamQuery, userID)
	if err != nil {
		return nil, err
	}

	userTeamMembersID := make([]string, 0)

	defer rows.Close()

	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		userTeamMembersID = append(userTeamMembersID, id)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userTeamMembersID, nil
}
