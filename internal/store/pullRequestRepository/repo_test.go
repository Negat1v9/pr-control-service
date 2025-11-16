package pullrequestrepository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Negat1v9/pr-review-service/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestCreatePullRequest(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	prRepo := NewPullRequestRepository()

	t.Run("Create", func(t *testing.T) {
		pr := models.PullRequest{
			ID:       "pr-111",
			Name:     "Author",
			AuthorID: "user-123",
		}

		mock.ExpectExec(createPullRequestQuery).WithArgs(&pr.ID, &pr.Name, &pr.AuthorID).WillReturnResult(sqlmock.NewResult(1, 1))

		err = prRepo.CreatePullRequest(context.Background(), sqlxDB, &pr)

		require.NoError(t, err)
	})
}

func TestGetPullRequestByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	prRepo := NewPullRequestRepository()

	t.Run("Get", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"pull_request_id", "pull_request_name", "author_id", "status"}).
			AddRow("pr-1000", "payment", "u1", "OPEN")

		rowsReviewers := sqlmock.NewRows([]string{"reviewer_user_id"}).AddRow("u2").AddRow("u3")
		pr := models.PullRequest{
			ID:       "pr-1000",
			Name:     "payment",
			AuthorID: "u1",
			Status:   models.PullRequestStatusOpen,
		}

		mock.ExpectQuery(getPullRequestByIDQuery).WithArgs(&pr.ID).WillReturnRows(rows)
		mock.ExpectQuery(getPullRequestReviewersQuery).WithArgs(&pr.ID).WillReturnRows(rowsReviewers)

		pullRequest, err := prRepo.GetPullRequestByID(context.Background(), sqlxDB, "pr-1000")

		require.NoError(t, err)
		require.NotNil(t, pullRequest)
		require.Equal(t, 2, len(pullRequest.AssignedReviewers))
	})

	t.Run("GetNotFound", func(t *testing.T) {
		mock.ExpectQuery(getPullRequestByIDQuery).WithArgs("nonexistent").WillReturnError(sql.ErrNoRows)

		pullRequest, err := prRepo.GetPullRequestByID(context.Background(), sqlxDB, "nonexistent")

		require.Error(t, err)
		require.Nil(t, pullRequest)
	})
}

func TestMergePullRequest(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	prRepo := NewPullRequestRepository()

	t.Run("Merge", func(t *testing.T) {
		mock.ExpectExec(mergePullRequestQuery).WithArgs("pr-123").WillReturnResult(sqlmock.NewResult(1, 1))

		err = prRepo.MergePullRequest(context.Background(), sqlxDB, "pr-123")

		require.NoError(t, err)
	})

	t.Run("MergeNotFound", func(t *testing.T) {
		mock.ExpectExec(mergePullRequestQuery).WithArgs("nonexistent").WillReturnResult(sqlmock.NewResult(0, 0))

		err = prRepo.MergePullRequest(context.Background(), sqlxDB, "nonexistent")

		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestGetQuantityPRReviewers(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	prRepo := NewPullRequestRepository()
	t.Run("GetQuantityReviewers", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"pull_request_id", "quantity_reviewers"}).
			AddRow("pr-1", 4).
			AddRow("pr-2", 1).
			AddRow("pr-3", 0)

		mock.ExpectQuery(getPullRequestsQuantityAssignedReviewers).WillReturnRows(rows)

		prInfo, err := prRepo.GetQuantityPRReviewers(context.Background(), sqlxDB)
		require.NoError(t, err)
		require.NotNil(t, prInfo)
	})
	t.Run("GetQuantityReviewers no one", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"pull_request_id", "quantity_reviewers"})

		mock.ExpectQuery(getPullRequestsQuantityAssignedReviewers).WillReturnRows(rows)

		prInfo, err := prRepo.GetQuantityPRReviewers(context.Background(), sqlxDB)
		require.NoError(t, err)
		require.Equal(t, 0, len(prInfo))
	})

}

func TestAssignReviewer(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	prRepo := NewPullRequestRepository()

	t.Run("Assign", func(t *testing.T) {
		mock.ExpectExec(createAssignedQuery).WithArgs("user-1", "pr-123").WillReturnResult(sqlmock.NewResult(1, 1))

		err = prRepo.AssignReviewer(context.Background(), sqlxDB, "pr-123", "user-1")

		require.NoError(t, err)
	})
}

func TestAssignManyReviewers(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	prRepo := NewPullRequestRepository()

	t.Run("AssignMany", func(t *testing.T) {
		reviewerIDs := []string{"user-1", "user-2", "user-3"}
		var placeholders []string
		for i := range reviewerIDs {
			offset := i * 2
			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", offset+1, offset+2))
		}
		query := fmt.Sprintf(createManyAssignedQuery, strings.Join(placeholders, ","))
		mock.ExpectExec(query).
			WithArgs("user-1", "pr-123", "user-2", "pr-123", "user-3", "pr-123").
			WillReturnResult(sqlmock.NewResult(3, 3))

		err = prRepo.AssignManyReviewers(context.Background(), sqlxDB, "pr-123", reviewerIDs)

		require.NoError(t, err)
	})
}

func TestDeleteAssignedByReviewerID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	prRepo := NewPullRequestRepository()

	t.Run("Delete", func(t *testing.T) {
		mock.ExpectExec(deleteAssignedByReviewerIDQuery).WithArgs("user-1").WillReturnResult(sqlmock.NewResult(1, 1))

		err = prRepo.DeleteAssignedByReviewerID(context.Background(), sqlxDB, "user-1")

		require.NoError(t, err)
	})

	t.Run("DeleteNotFound", func(t *testing.T) {
		mock.ExpectExec(deleteAssignedByReviewerIDQuery).WithArgs("nonexistent").WillReturnResult(sqlmock.NewResult(0, 0))

		err = prRepo.DeleteAssignedByReviewerID(context.Background(), sqlxDB, "nonexistent")

		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}
