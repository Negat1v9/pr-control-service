package userrepository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Negat1v9/pr-review-service/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	userRepo := NewUserRepository()

	t.Run("Create user", func(t *testing.T) {
		user := models.User{
			UserID:   "u1",
			Username: "u1",
			IsActive: true,
			TeamName: "payment",
		}

		mock.ExpectExec(createUserQuery).
			WithArgs(user.UserID, user.Username, user.IsActive, user.TeamName).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = userRepo.CreateUser(context.Background(), sqlxDB, user.TeamName, &user)

		require.NoError(t, err)
	})
}

func TestCreateManyUsers(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	userRepo := NewUserRepository()

	t.Run("Create many users", func(t *testing.T) {
		users := []models.User{
			{UserID: "u1", Username: "U1", IsActive: true},
			{UserID: "u2", Username: "U2", IsActive: true},
			{UserID: "u3", Username: "U3", IsActive: false},
		}

		teamName := "payment"

		var placeholders []string
		for i := range users {
			offset := i * 4
			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d)", offset+1, offset+2, offset+3, offset+4))
		}
		query := fmt.Sprintf(createManyUsersQuery, strings.Join(placeholders, ","))
		mock.ExpectExec(query).
			WithArgs("u1", "U1", true, teamName, "u2", "U2", true, teamName, "u3", "U3", false, teamName).
			WillReturnResult(sqlmock.NewResult(3, 3))

		err = userRepo.CreateManyUsers(context.Background(), sqlxDB, teamName, users)

		require.NoError(t, err)
	})
}

func TestGetUserReviews(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	userRepo := NewUserRepository()

	t.Run("Get user reviews", func(t *testing.T) {
		userID := "u123"
		timeCreatedAt1 := time.Now().Add(-1 * time.Minute)
		timeCreatedAt2 := time.Now().Add(-5 * time.Minute)
		mergedAt2 := time.Now().Add(-3 * time.Minute)

		rows := sqlmock.NewRows([]string{"pull_request_id", "pull_request_name", "author_id", "status", "created_at", "merged_at"}).
			AddRow("pr-1", "feature-1", "u456", "OPEN", timeCreatedAt1, nil).
			AddRow("pr-2", "bugfix-1", "u789", "MERGED", timeCreatedAt2, mergedAt2)
		mock.ExpectQuery(getUserReviewsQuery).WithArgs(userID).WillReturnRows(rows)

		userReviews, err := userRepo.GetUserReviews(context.Background(), sqlxDB, userID)

		require.NoError(t, err)
		require.NotNil(t, userReviews)
		require.Equal(t, userID, userReviews.UserID)
		require.Equal(t, 2, len(userReviews.PullRequests))
	})
}

func TestUpdateUserStatus(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	userRepo := NewUserRepository()

	t.Run("Update user", func(t *testing.T) {
		userID := "u123"
		rows := sqlmock.NewRows([]string{"user_id", "username", "is_active", "team_name"}).
			AddRow(userID, "u1", true, "pay")

		mock.ExpectQuery(updateUserStatusQuery).WithArgs(true, userID).WillReturnRows(rows)

		updatedUser, err := userRepo.UpdateUserStatus(context.Background(), sqlxDB, userID, true)

		require.NoError(t, err)
		require.NotNil(t, updatedUser)
		require.Equal(t, userID, updatedUser.UserID)
		require.Equal(t, true, updatedUser.IsActive)
	})
	t.Run("UpdateNotFound", func(t *testing.T) {
		userID := "nonexistent"

		mock.ExpectQuery(updateUserStatusQuery).WithArgs(true, userID).WillReturnError(sql.ErrNoRows)

		updatedUser, err := userRepo.UpdateUserStatus(context.Background(), sqlxDB, userID, true)

		require.Error(t, err)
		require.NotNil(t, updatedUser)
	})
}
