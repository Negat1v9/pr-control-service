package prhttp

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Negat1v9/pr-review-service/internal/models"
	prservice "github.com/Negat1v9/pr-review-service/internal/pullRequest/service"
	mock_store "github.com/Negat1v9/pr-review-service/internal/store/mock"
	"github.com/Negat1v9/pr-review-service/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPRRepo := mock_store.NewMockPullRequestRepository(ctrl)
	mockTeamRepo := mock_store.NewMockTeamRepository(ctrl)
	mockStore := mock_store.NewMockStore(ctrl)

	newPR := models.CreatePullRequest{
		ID:       "pr-1",
		Name:     "pr-name",
		AuthorID: "userID",
	}
	prResult := models.PullRequest{
		ID:                "pr-1",
		Name:              "pr-name",
		AuthorID:          "userID",
		Status:            models.PullRequestStatusOpen,
		AssignedReviewers: []string{"u1", "u2"},
	}
	db := sqlx.DB{}

	mockStore.EXPECT().DB().Return(&db).AnyTimes()

	mockStore.EXPECT().PRRepo().Return(mockPRRepo).AnyTimes()

	mockStore.EXPECT().TeamRepo().Return(mockTeamRepo).AnyTimes()

	doReq := func() *httptest.ResponseRecorder {
		service := prservice.NewPRService(mockStore)
		handler := NewPRHanlder(logger.NewLogger("local"), service)
		prMux := PRRouter(handler)

		data, err := json.Marshal(&newPR)
		require.NoError(t, err)
		req, err := http.NewRequest("POST", "/create", bytes.NewBuffer(data))
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		prMux.ServeHTTP(rr, req)
		return rr
	}

	t.Run("Create success", func(t *testing.T) {

		mockPRRepo.EXPECT().GetPullRequestByID(gomock.Any(), gomock.Any(), "pr-1").Return(nil, sql.ErrNoRows).Times(1)

		mockTeamRepo.EXPECT().GetUsersIDFromUserTeam(gomock.Any(), gomock.Any(), "userID", 2).Return([]string{"u1", "u2"}, nil).Times(1)
		mockPRRepo.EXPECT().CreatePullRequest(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
		mockPRRepo.EXPECT().AssignManyReviewers(gomock.Any(), gomock.Any(), "pr-1", []string{"u1", "u2"}).Return(nil).Times(1)

		mockStore.EXPECT().DoTx(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(context.Context, sqlx.ExtContext) error) error {
				return fn(ctx, &db)
			},
		).Times(1)

		mockPRRepo.EXPECT().GetPullRequestByID(gomock.Any(), gomock.Any(), "pr-1").Return(&prResult, nil).Times(1)
		rr := doReq()
		require.Equal(t, 201, rr.Code)
	})

	t.Run("Alredy created", func(t *testing.T) {
		mockPRRepo.EXPECT().GetPullRequestByID(gomock.Any(), gomock.Any(), "pr-1").Return(&prResult, nil).Times(1)
		rr := doReq()
		require.Equal(t, 409, rr.Code)
	})

	t.Run("Author not found", func(t *testing.T) {
		mockPRRepo.EXPECT().GetPullRequestByID(gomock.Any(), gomock.Any(), "pr-1").Return(nil, sql.ErrNoRows).Times(1)
		mockTeamRepo.EXPECT().GetUsersIDFromUserTeam(gomock.Any(), gomock.Any(), "userID", 2).Return([]string{}, sql.ErrNoRows).Times(1)
		mockStore.EXPECT().DoTx(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(context.Context, sqlx.ExtContext) error) error {
				return fn(ctx, &db)
			},
		).Times(1)
		rr := doReq()
		require.Equal(t, 404, rr.Code)
	})
}

func TestMerge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPRRepo := mock_store.NewMockPullRequestRepository(ctrl)
	mockTeamRepo := mock_store.NewMockTeamRepository(ctrl)
	mockStore := mock_store.NewMockStore(ctrl)

	newPR := models.CreatePullRequest{
		ID:       "pr-1",
		Name:     "pr-name",
		AuthorID: "userID",
	}
	prResult := models.PullRequest{
		ID:                "pr-1",
		Name:              "pr-name",
		AuthorID:          "userID",
		Status:            models.PullRequestStatusOpen,
		AssignedReviewers: []string{"u1", "u2"},
	}
	db := sqlx.DB{}

	mockStore.EXPECT().DB().Return(&db).AnyTimes()

	mockStore.EXPECT().PRRepo().Return(mockPRRepo).AnyTimes()

	mockStore.EXPECT().TeamRepo().Return(mockTeamRepo).AnyTimes()

	doReq := func() *httptest.ResponseRecorder {
		service := prservice.NewPRService(mockStore)
		handler := NewPRHanlder(logger.NewLogger("local"), service)
		prMux := PRRouter(handler)

		data, err := json.Marshal(&newPR)
		require.NoError(t, err)
		req, err := http.NewRequest("POST", "/merge", bytes.NewBuffer(data))
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		prMux.ServeHTTP(rr, req)
		return rr
	}

	t.Run("Merge success", func(t *testing.T) {
		mergedPR := models.PullRequest{
			ID:       "pr-1",
			Name:     "pr-name",
			AuthorID: "userID",
			Status:   models.PullRequestStatusMerged,
		}
		mockPRRepo.EXPECT().GetPullRequestByID(gomock.Any(), gomock.Any(), "pr-1").Return(&prResult, nil).Times(1)

		mockPRRepo.EXPECT().MergePullRequest(gomock.Any(), gomock.Any(), "pr-1").Return(nil).Times(1)

		mockPRRepo.EXPECT().GetPullRequestByID(gomock.Any(), gomock.Any(), "pr-1").Return(&mergedPR, nil).Times(1)

		rr := doReq()
		r := map[string]any{}
		err := json.Unmarshal(rr.Body.Bytes(), &r)
		require.NoError(t, err)
		require.Equal(t, 200, rr.Code)
		require.Equal(t, "pr-1", r["pr"].(map[string]any)["pull_request_id"])
		require.Equal(t, "MERGED", r["pr"].(map[string]any)["status"])
	})

	t.Run("PR not found", func(t *testing.T) {
		mockPRRepo.EXPECT().GetPullRequestByID(gomock.Any(), gomock.Any(), "pr-1").Return(nil, sql.ErrNoRows).Times(1)

		rr := doReq()
		require.Equal(t, 404, rr.Code)
	})

	t.Run("PR already merged", func(t *testing.T) {
		mergedPR := models.PullRequest{
			ID:       "pr-1",
			Name:     "pr-name",
			AuthorID: "userID",
			Status:   models.PullRequestStatusMerged,
		}
		mockPRRepo.EXPECT().GetPullRequestByID(gomock.Any(), gomock.Any(), "pr-1").Return(&mergedPR, nil).Times(1)

		rr := doReq()
		require.Equal(t, 200, rr.Code)
	})
}

func TestReassign(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPRRepo := mock_store.NewMockPullRequestRepository(ctrl)
	mockTeamRepo := mock_store.NewMockTeamRepository(ctrl)
	mockStore := mock_store.NewMockStore(ctrl)

	newReviewers := []string{"u3"}
	pr := models.PullRequest{
		ID:                "pr-1",
		Name:              "pr-name",
		AuthorID:          "userID",
		Status:            models.PullRequestStatusOpen,
		AssignedReviewers: []string{"u1", "u2"},
	}

	reasignReq := models.ReassignPullRequest{
		ID:            "pr-1",
		OldReviewerID: "u1",
	}

	db := sqlx.DB{}

	mockStore.EXPECT().DB().Return(&db).AnyTimes()

	mockStore.EXPECT().PRRepo().Return(mockPRRepo).AnyTimes()

	mockStore.EXPECT().TeamRepo().Return(mockTeamRepo).AnyTimes()

	doReq := func() *httptest.ResponseRecorder {
		service := prservice.NewPRService(mockStore)
		handler := NewPRHanlder(logger.NewLogger("local"), service)
		prMux := PRRouter(handler)

		data, err := json.Marshal(&reasignReq)
		require.NoError(t, err)
		req, err := http.NewRequest("POST", "/reassign", bytes.NewBuffer(data))
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		prMux.ServeHTTP(rr, req)
		return rr
	}

	t.Run("Reassign success", func(t *testing.T) {
		updPR := models.PullRequest{
			ID:                "pr-1",
			Name:              "pr-name",
			AuthorID:          "userID",
			Status:            models.PullRequestStatusOpen,
			AssignedReviewers: []string{"u3", "u2"},
		}
		mockPRRepo.EXPECT().GetPullRequestByID(gomock.Any(), gomock.Any(), "pr-1").Return(&pr, nil).Times(1)
		mockTeamRepo.EXPECT().GetActiveUsersTeamWithException(gomock.Any(), gomock.Any(), "userID", []string{"u1", "u2"}, 1).Return(newReviewers, nil)
		mockStore.EXPECT().DoTx(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(context.Context, sqlx.ExtContext) error) error {
				return fn(ctx, &db)
			},
		).Times(1)
		mockPRRepo.EXPECT().DeleteAssignedByReviewerID(gomock.Any(), gomock.Any(), "u1").Return(nil)
		mockPRRepo.EXPECT().AssignReviewer(gomock.Any(), gomock.Any(), "pr-1", "u3").Return(nil)
		mockPRRepo.EXPECT().GetPullRequestByID(gomock.Any(), gomock.Any(), "pr-1").Return(&updPR, nil).Times(1)
		rr := doReq()
		require.Equal(t, 200, rr.Code)
		r := make(map[string]any, 0)
		err := json.Unmarshal(rr.Body.Bytes(), &r)
		require.NoError(t, err)
		require.Equal(t, "u3", r["replaced_by"])
		require.Equal(t, []any{"u3", "u2"}, r["pr"].(map[string]any)["assigned_reviewers"])
	})

	t.Run("Reassign alredy merged", func(t *testing.T) {
		merged := models.PullRequest{
			ID:     "pr-1",
			Status: models.PullRequestStatusMerged,
		}
		mockPRRepo.EXPECT().GetPullRequestByID(gomock.Any(), gomock.Any(), "pr-1").Return(&merged, nil).Times(1)
		rr := doReq()
		require.Equal(t, 409, rr.Code)
		r := make(map[string]any, 0)
		err := json.Unmarshal(rr.Body.Bytes(), &r)
		require.NoError(t, err)
		require.Equal(t, "PR_MERGED", r["error"].(map[string]any)["code"])
		require.Equal(t, "cannot reassign on merged PR", r["error"].(map[string]any)["message"])
	})
	t.Run("Reassign not found", func(t *testing.T) {
		mockPRRepo.EXPECT().GetPullRequestByID(gomock.Any(), gomock.Any(), "pr-1").Return(nil, sql.ErrNoRows).Times(1)
		rr := doReq()
		require.Equal(t, 404, rr.Code)
		r := make(map[string]any, 0)
		err := json.Unmarshal(rr.Body.Bytes(), &r)
		require.NoError(t, err)
		require.Equal(t, "NOT_FOUND", r["error"].(map[string]any)["code"])
		require.Equal(t, "resource not found", r["error"].(map[string]any)["message"])
	})

	t.Run("Reassign old user not reviewer", func(t *testing.T) {
		notReviewwerPR := models.PullRequest{
			ID:                "pr-1",
			AssignedReviewers: []string{"u10", "u100"},
		}
		mockPRRepo.EXPECT().GetPullRequestByID(gomock.Any(), gomock.Any(), "pr-1").Return(&notReviewwerPR, nil).Times(1)
		rr := doReq()
		require.Equal(t, 409, rr.Code)
		r := make(map[string]any, 0)
		err := json.Unmarshal(rr.Body.Bytes(), &r)
		require.NoError(t, err)
		require.Equal(t, "NOT_ASSIGNED", r["error"].(map[string]any)["code"])
		require.Equal(t, "reviewer is not assigned to this PR", r["error"].(map[string]any)["message"])
	})

	t.Run("Reassign no candidates", func(t *testing.T) {
		noCandidate := models.PullRequest{
			ID:                "pr-1",
			AuthorID:          "userID",
			Status:            models.PullRequestStatusOpen,
			AssignedReviewers: []string{"u1", "u2"},
		}
		mockPRRepo.EXPECT().GetPullRequestByID(gomock.Any(), gomock.Any(), "pr-1").Return(&noCandidate, nil).Times(1)
		mockTeamRepo.EXPECT().GetActiveUsersTeamWithException(gomock.Any(), gomock.Any(), "userID", []string{"u1", "u2"}, 1).Return(nil, sql.ErrNoRows)
		rr := doReq()
		require.Equal(t, 409, rr.Code)
		r := make(map[string]any, 0)
		err := json.Unmarshal(rr.Body.Bytes(), &r)
		require.NoError(t, err)
		require.Equal(t, "NO_CANDIDATE", r["error"].(map[string]any)["code"])
		require.Equal(t, "no active replacement candidate in team", r["error"].(map[string]any)["message"])

	})
}
