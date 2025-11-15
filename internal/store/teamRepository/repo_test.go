package teamrepository

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestCreateTeam(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	teamRepo := NewTeamRepositiry()

	t.Run("Create team", func(t *testing.T) {

		mock.ExpectExec(createTeamQuery).WithArgs("team-1").WillReturnResult(sqlmock.NewResult(1, 1))

		err := teamRepo.CreateTeam(context.Background(), sqlxDB, "team-1")
		require.NoError(t, err)
	})

}

func TestGetTeamWithMembers(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	teamRepo := NewTeamRepositiry()

	t.Run("Get team with memmbers", func(t *testing.T) {

		rowsTeam := sqlmock.NewRows([]string{"team_name", "user_id", "username", "is_active"}).
			AddRow("team-1", "userID", "username", true).
			AddRow("team-1", "userID-1", "username2", false)

		mock.ExpectQuery(getTeamWithUsersByNameQuery).
			WithArgs("team-1").WillReturnRows(rowsTeam)

		team, err := teamRepo.GetTeamWithMembers(context.Background(), sqlxDB, "team-1")
		require.NoError(t, err)
		require.Equal(t, 2, len(team.Members))
	})
}

func TestGetUsersIDFromUserTeam(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	teamRepo := NewTeamRepositiry()

	t.Run("Get User IDs from user team", func(t *testing.T) {

		rowsTeamUsers := sqlmock.NewRows([]string{"user_id"}).
			AddRow("userID").
			AddRow("userID-1").
			AddRow("userID-2")

		mock.ExpectQuery(getActiveUserIDFromUserTeamQuery).
			WithArgs("userID", 3).WillReturnRows(rowsTeamUsers)

		userIDs, err := teamRepo.GetUsersIDFromUserTeam(context.Background(), sqlxDB, "userID", 2)
		require.NoError(t, err)
		require.Equal(t, 2, len(userIDs))
	})
}

func TestGetActiveUsersTeamWithException(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	teamRepo := NewTeamRepositiry()

	t.Run("Get Active Users Team With Exception", func(t *testing.T) {
		exceptionUsers := []string{"userID-1", "userID-2"}
		rowsTeamUsers := sqlmock.NewRows([]string{"user_id"}).
			AddRow("userID").
			AddRow("userID-4").
			AddRow("userID-3")

		placeholders := make([]string, 0, len(exceptionUsers))
		for _, expUserID := range exceptionUsers {
			placeholders = append(placeholders, fmt.Sprintf("'%s'", expUserID))
		}
		query := fmt.Sprintf(getActiveUserFromUserTeamWithException, strings.Join(placeholders, ","))

		mock.ExpectQuery(query).
			WithArgs("userID", 3).WillReturnRows(rowsTeamUsers)

		userIDs, err := teamRepo.GetActiveUsersTeamWithException(context.Background(), sqlxDB, "userID", exceptionUsers, 2)
		require.NoError(t, err)
		require.Equal(t, 2, len(userIDs))
	})
}
