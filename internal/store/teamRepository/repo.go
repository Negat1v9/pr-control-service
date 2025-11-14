package teamrepository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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

// return all active users team and empty array if user with userID does nit exists
// if user exists but he is one active user in team return sql.ErrNoRows
func (r *teamRepositiry) GetUsersIDFromUserTeam(ctx context.Context, exec sqlx.ExtContext, userID string, limit int) ([]string, error) {
	// limit+1 select reasoon -> target id selected also
	rows, err := exec.QueryxContext(ctx, getActiveUserIDFromUserTeamQuery, userID, limit+1)
	if err != nil {
		return nil, err
	}

	userTeamMembersID := make([]string, 0, limit)

	defer rows.Close()

	exitstUserID := false
	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		// user with userID exists
		if id == userID {
			exitstUserID = true
			// not add in result
			continue
		}
		userTeamMembersID = append(userTeamMembersID, id)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// user not exists
	if !exitstUserID {
		return nil, sql.ErrNoRows
	}

	return userTeamMembersID, nil
}

func (r *teamRepositiry) GetActiveUsersTeamWithException(ctx context.Context, exec sqlx.ExtContext, userID string, exceptions []string, limit int) ([]string, error) {
	var placeholders []string

	for _, expUserID := range exceptions {
		placeholders = append(placeholders, fmt.Sprintf("'%s'", expUserID))
	}
	query := fmt.Sprintf(getActiveUserFromUserTeamWithException, strings.Join(placeholders, ","))
	// limit select reasoon -> target id selected also
	rows, err := exec.QueryxContext(ctx, query, userID, limit+1)
	if err != nil {
		return nil, err
	}

	// limit select reasoon -> target id selected also
	userTeamMembersID := make([]string, 0, limit)

	defer rows.Close()
	exitstUserID := false

	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		// user with userID exists
		if id == userID {
			exitstUserID = true
			// not add in result
			continue
		}
		userTeamMembersID = append(userTeamMembersID, id)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// user not exists
	if !exitstUserID {
		return nil, sql.ErrNoRows
	}

	if len(userTeamMembersID) == 0 {
		return nil, sql.ErrNoRows
	}

	return userTeamMembersID, nil
}
