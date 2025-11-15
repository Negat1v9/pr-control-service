package userhttp

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Negat1v9/pr-review-service/internal/models"
	mock_store "github.com/Negat1v9/pr-review-service/internal/store/mock"
	userservice "github.com/Negat1v9/pr-review-service/internal/users/service"
	"github.com/Negat1v9/pr-review-service/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSetIsActive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_store.NewMockUserRepository(ctrl)
	mockStore := mock_store.NewMockStore(ctrl)
	mockStore.EXPECT().UserRepo().Return(mockUserRepo).AnyTimes()

	db := sqlx.DB{}
	mockStore.EXPECT().DB().Return(&db).AnyTimes()

	doReq := func(body any) *httptest.ResponseRecorder {
		service := userservice.NewUserService(mockStore)
		handler := NewUserHandler(logger.NewLogger("local"), service)
		userMux := UserRouter(handler)

		data, err := json.Marshal(body)
		require.NoError(t, err)
		req, err := http.NewRequest("POST", "/setIsActive", bytes.NewBuffer(data))
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		userMux.ServeHTTP(rr, req)
		return rr
	}

	t.Run("Set is active", func(t *testing.T) {
		req := models.SetUserActiveStatusRequest{
			UserID:   "user-1",
			IsActive: false,
		}
		updatedUser := models.User{
			UserID:   "user-1",
			Username: "username-1",
			IsActive: false,
			TeamName: "team-1",
		}

		mockUserRepo.EXPECT().UpdateUserStatus(gomock.Any(), gomock.Any(), req.UserID, false).
			Return(&updatedUser, nil)

		rr := doReq(req)
		r := map[string]any{}
		err := json.Unmarshal(rr.Body.Bytes(), &r)
		require.NoError(t, err)
		require.Equal(t, "user-1", r["user"].(map[string]any)["user_id"])
		require.Equal(t, false, r["user"].(map[string]any)["is_active"])
	})

	t.Run("Not found user", func(t *testing.T) {
		req := models.SetUserActiveStatusRequest{
			UserID:   "user-2",
			IsActive: true,
		}

		mockUserRepo.EXPECT().UpdateUserStatus(gomock.Any(), gomock.Any(), req.UserID, true).Return(nil, sql.ErrNoRows)
		rr := doReq(req)
		require.Equal(t, http.StatusNotFound, rr.Code)
		r := map[string]any{}
		err := json.Unmarshal(rr.Body.Bytes(), &r)
		require.NoError(t, err)
		require.Equal(t, "NOT_FOUND", r["error"].(map[string]any)["code"])
		require.Equal(t, "resource not found", r["error"].(map[string]any)["message"])
	})
}

func TestGetReview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_store.NewMockUserRepository(ctrl)
	mockStore := mock_store.NewMockStore(ctrl)
	mockStore.EXPECT().UserRepo().Return(mockUserRepo).AnyTimes()

	db := sqlx.DB{}
	mockStore.EXPECT().DB().Return(&db).AnyTimes()

	doReq := func(userID string) *httptest.ResponseRecorder {
		service := userservice.NewUserService(mockStore)
		handler := NewUserHandler(logger.NewLogger("local"), service)
		userMux := UserRouter(handler)

		req, err := http.NewRequest("GET", "/getReview?user_id="+userID, nil)

		require.NoError(t, err)
		rr := httptest.NewRecorder()

		userMux.ServeHTTP(rr, req)
		return rr
	}

	t.Run("Get reviews", func(t *testing.T) {
		userID := "user-1"

		userReviews := &models.UserReviews{
			UserID: userID,
			PullRequests: []models.PullRequest{
				{
					ID:                "pr-1",
					Name:              "feature-1",
					AuthorID:          "user-2",
					Status:            models.PullRequestStatusOpen,
					AssignedReviewers: []string{"user-1"},
				},
				{
					ID:                "pr-2",
					Name:              "bugfix-1",
					AuthorID:          "user-3",
					Status:            models.PullRequestStatusMerged,
					AssignedReviewers: []string{"user-1"},
				},
			},
		}

		mockUserRepo.EXPECT().GetUserReviews(gomock.Any(), gomock.Any(), userID).
			Return(userReviews, nil)

		rr := doReq(userID)

		require.Equal(t, http.StatusOK, rr.Code)
		r := map[string]any{}
		err := json.Unmarshal(rr.Body.Bytes(), &r)
		require.NoError(t, err)
		require.Equal(t, userID, r["user_id"])
		require.Equal(t, 2.0, float64(len(r["pull_requests"].([]interface{}))))
	})

	t.Run("Not found user", func(t *testing.T) {
		userID := "user-2"
		mockUserRepo.EXPECT().GetUserReviews(gomock.Any(), gomock.Any(), userID).Return(nil, sql.ErrNoRows)
		rr := doReq(userID)
		require.Equal(t, http.StatusNotFound, rr.Code)
		r := map[string]any{}
		err := json.Unmarshal(rr.Body.Bytes(), &r)
		require.NoError(t, err)
		require.Equal(t, "NOT_FOUND", r["error"].(map[string]any)["code"])
		require.Equal(t, "resource not found", r["error"].(map[string]any)["message"])
	})
}
