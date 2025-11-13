package store

import (
	"context"

	"github.com/Negat1v9/pr-review-service/internal/models"
	teamrepository "github.com/Negat1v9/pr-review-service/internal/store/teamRepository"
	userrepository "github.com/Negat1v9/pr-review-service/internal/store/userRepository"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(ctx context.Context, exec sqlx.ExtContext, user *models.User) error
	CreateManyUsers(ctx context.Context, exec sqlx.ExtContext, users []models.User) error
	GetUserReviews(ctx context.Context, exec sqlx.ExtContext, userID string) (*models.UserReviews, error)
}

type TeamRepository interface {
	CreateTeam(ctx context.Context, exec sqlx.ExtContext, teamName string) error
	GetTeamWithMembers(ctx context.Context, exec sqlx.ExtContext, teamName string) (*models.Team, error)
	GetUsersIDFromUserTeam(ctx context.Context, exec sqlx.ExtContext, userID string) ([]string, error)
	CreateTeamMember(ctx context.Context, exec sqlx.ExtContext, userID, teamName string) error
	CreateManyTeamMembers(ctx context.Context, exec sqlx.ExtContext, teamMembers *models.Team) error
}

type PullRequestRepository interface {
}

type Store struct {
	Db       *sqlx.DB
	TeamRepo TeamRepository
	UserRepo UserRepository
	// pullRequestRepo PullRequestRepository
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		Db:       db,
		TeamRepo: teamrepository.NewTeamRepositiry(),
		UserRepo: userrepository.NewUserRepository(),
		// pullRequestRepo: pullrequestrepository,
	}
}

func (s *Store) DoTx(ctx context.Context, fn func(ctx context.Context, exec sqlx.ExtContext) error) error {
	tx, err := s.Db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(ctx, tx); err != nil {
		if rberr := tx.Rollback(); rberr != nil {
			// FIXME: logging
		}
		return err
	}

	return tx.Commit()
}
