package store

import (
	"context"

	"github.com/Negat1v9/pr-review-service/internal/models"
	pullrequestrepository "github.com/Negat1v9/pr-review-service/internal/store/pullRequestRepository"
	teamrepository "github.com/Negat1v9/pr-review-service/internal/store/teamRepository"
	userrepository "github.com/Negat1v9/pr-review-service/internal/store/userRepository"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(ctx context.Context, exec sqlx.ExtContext, teamName string, user *models.User) error
	CreateManyUsers(ctx context.Context, exec sqlx.ExtContext, teamName string, users []models.User) error
	GetUserReviews(ctx context.Context, exec sqlx.ExtContext, userID string) (*models.UserReviews, error)
	UpdateUserStatus(ctx context.Context, exec sqlx.ExtContext, userID string, isActive bool) (*models.User, error)
}

type TeamRepository interface {
	CreateTeam(ctx context.Context, exec sqlx.ExtContext, teamName string) error
	GetTeamWithMembers(ctx context.Context, exec sqlx.ExtContext, teamName string) (*models.Team, error)
	// returns all userIDs from user team
	GetUsersIDFromUserTeam(ctx context.Context, exec sqlx.ExtContext, userID string, limit int) ([]string, error)
	// return active usersID from userID team without exceptions users
	// if there are no active users or the user himself, it returns sql.ErrNoRows error
	GetActiveUsersTeamWithException(ctx context.Context, exec sqlx.ExtContext, userID string, exceptions []string, limit int) ([]string, error)
}

type PullRequestRepository interface {
	CreatePullRequest(ctx context.Context, exec sqlx.ExtContext, pr *models.PullRequest) error
	GetPullRequestByID(ctx context.Context, exec sqlx.ExtContext, prID string) (*models.PullRequest, error)
	MergePullRequest(ctx context.Context, exec sqlx.ExtContext, prID string) error
	AssignReviewer(ctx context.Context, exec sqlx.ExtContext, prID, reviewerID string) error
	AssignManyReviewers(ctx context.Context, exec sqlx.ExtContext, prID string, reviewerIDs []string) error
	DeleteAssignedByReviewerID(ctx context.Context, exec sqlx.ExtContext, reviewerID string) error
}

type Store struct {
	Db       *sqlx.DB
	TeamRepo TeamRepository
	UserRepo UserRepository
	PRRepo   PullRequestRepository
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		Db:       db,
		TeamRepo: teamrepository.NewTeamRepositiry(),
		UserRepo: userrepository.NewUserRepository(),
		PRRepo:   pullrequestrepository.NewPullRequestRepository(),
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
