package teamhttp

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Negat1v9/pr-review-service/internal/models"
	mock_store "github.com/Negat1v9/pr-review-service/internal/store/mock"
	teamservice "github.com/Negat1v9/pr-review-service/internal/team/service"
	"github.com/Negat1v9/pr-review-service/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAdd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTeamRepo := mock_store.NewMockTeamRepository(ctrl)
	mockUserRepo := mock_store.NewMockUserRepository(ctrl)
	mockStore := mock_store.NewMockStore(ctrl)

	mockStore.EXPECT().TeamRepo().Return(mockTeamRepo).AnyTimes()
	mockStore.EXPECT().UserRepo().Return(mockUserRepo).AnyTimes()

	db := sqlx.DB{}
	mockStore.EXPECT().DB().Return(&db).AnyTimes()

	doReq := func(body any) *httptest.ResponseRecorder {
		service := teamservice.NewTeamService(mockStore)
		handler := NewTeamHanlder(logger.NewLogger("local"), service)
		teamMux := TeamRouter(handler)

		data, err := json.Marshal(body)
		require.NoError(t, err)
		req, err := http.NewRequest("POST", "/add", bytes.NewBuffer(data))
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		teamMux.ServeHTTP(rr, req)
		return rr
	}

	t.Run("Add team success", func(t *testing.T) {
		newTeam := models.Team{
			TeamName: "bb",
			Members: []models.User{
				{UserID: "u1", Username: "u1", IsActive: true},
				{UserID: "u2", Username: "u2", IsActive: false},
			},
		}

		mockTeamRepo.EXPECT().GetTeamWithMembers(gomock.Any(), gomock.Any(), "bb").
			Return(nil, sql.ErrNoRows)

		mockStore.EXPECT().DoTx(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(context.Context, sqlx.ExtContext) error) error {
				return fn(ctx, &db)
			},
		)

		mockTeamRepo.EXPECT().CreateTeam(gomock.Any(), gomock.Any(), "bb").Return(nil)
		mockUserRepo.EXPECT().CreateManyUsers(gomock.Any(), gomock.Any(), "bb", newTeam.Members).Return(nil)
		mockTeamRepo.EXPECT().GetTeamWithMembers(gomock.Any(), gomock.Any(), "bb").Return(&newTeam, nil)

		rr := doReq(newTeam)

		require.Equal(t, http.StatusCreated, rr.Code)
		r := map[string]any{}
		err := json.Unmarshal(rr.Body.Bytes(), &r)
		require.NoError(t, err)
		require.Equal(t, "bb", r["team"].(map[string]any)["team_name"])
		require.Equal(t, 2, len(r["team"].(map[string]any)["members"].([]any)))

	})

	t.Run("Add team already exists", func(t *testing.T) {
		newTeam := models.Team{
			TeamName: "bb",
		}

		existingTeam := &models.Team{
			TeamName: "bb",
		}

		mockTeamRepo.EXPECT().GetTeamWithMembers(gomock.Any(), gomock.Any(), "bb").
			Return(existingTeam, nil)

		rr := doReq(newTeam)

		require.Equal(t, 400, rr.Code)
		r := map[string]any{}
		err := json.Unmarshal(rr.Body.Bytes(), &r)
		require.NoError(t, err)
		require.Equal(t, "TEAM_EXISTS", r["error"].(map[string]any)["code"])
		require.Equal(t, "team_name already exists", r["error"].(map[string]any)["message"])
	})
}

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTeamRepo := mock_store.NewMockTeamRepository(ctrl)
	mockStore := mock_store.NewMockStore(ctrl)
	mockStore.EXPECT().TeamRepo().Return(mockTeamRepo).AnyTimes()

	db := sqlx.DB{}
	mockStore.EXPECT().DB().Return(&db).AnyTimes()

	doReq := func(teamName string) *httptest.ResponseRecorder {
		service := teamservice.NewTeamService(mockStore)
		handler := NewTeamHanlder(logger.NewLogger("local"), service)
		teamMux := TeamRouter(handler)

		req, err := http.NewRequest("GET", "/get?team_name="+teamName, nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		teamMux.ServeHTTP(rr, req)
		return rr
	}

	t.Run("Get team success", func(t *testing.T) {
		teamName := "bb"
		team := &models.Team{
			TeamName: teamName,
			Members: []models.User{
				{UserID: "u1", Username: "u1", IsActive: true, TeamName: teamName},
				{UserID: "u2", Username: "u2", IsActive: true, TeamName: teamName},
				{UserID: "u3", Username: "u3", IsActive: false, TeamName: teamName},
			},
		}

		mockTeamRepo.EXPECT().GetTeamWithMembers(gomock.Any(), gomock.Any(), teamName).
			Return(team, nil)

		rr := doReq(teamName)

		require.Equal(t, 200, rr.Code)
		r := map[string]any{}
		err := json.Unmarshal(rr.Body.Bytes(), &r)
		require.NoError(t, err)
		require.Equal(t, teamName, r["team_name"])
		require.Equal(t, 3, len(r["members"].([]any)))
	})

	t.Run("Get team not found", func(t *testing.T) {
		teamName := "bb"

		mockTeamRepo.EXPECT().GetTeamWithMembers(gomock.Any(), gomock.Any(), teamName).
			Return(nil, sql.ErrNoRows)

		rr := doReq(teamName)

		require.Equal(t, 404, rr.Code)
		r := map[string]any{}
		err := json.Unmarshal(rr.Body.Bytes(), &r)
		require.NoError(t, err)
		require.Equal(t, "NOT_FOUND", r["error"].(map[string]any)["code"])
		require.Equal(t, "resource not found", r["error"].(map[string]any)["message"])
	})
}
