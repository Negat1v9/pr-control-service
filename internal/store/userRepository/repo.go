package userrepository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/Negat1v9/pr-review-service/internal/models"
)

type userRepository struct {
}

func NewUserRepository() *userRepository {
	return &userRepository{}
}

func (r *userRepository) CreateUser(ctx context.Context, exec sqlx.ExtContext, teamName string, user *models.User) error {
	_, err := exec.ExecContext(ctx, createUserQuery, user.UserID, user.Username, user.IsActive, teamName)
	return err
}

func (r *userRepository) CreateManyUsers(ctx context.Context, exec sqlx.ExtContext, teamName string, users []models.User) error {
	if len(users) == 0 {
		return fmt.Errorf("userRepository.CreateManyUsers: no users")
	}

	var placeholders []string
	var args []any

	for i, user := range users {
		offset := i * 4
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d)", offset+1, offset+2, offset+3, offset+4))
		args = append(args, user.UserID, user.Username, user.IsActive, teamName)
	}

	query := fmt.Sprintf(createManyUsersQuery, strings.Join(placeholders, ","))

	_, err := exec.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetUserReviews(ctx context.Context, exec sqlx.ExtContext, userID string) (*models.UserReviews, error) {

	rows, err := exec.QueryxContext(ctx, getUserReviewsQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userReviews := &models.UserReviews{
		UserID:       userID,
		PullRequests: make([]models.PullRequest, 0),
	}

	for rows.Next() {
		var review models.PullRequest
		if err := rows.StructScan(&review); err != nil {
			return nil, err
		}
		userReviews.PullRequests = append(userReviews.PullRequests, review)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userReviews, nil
}
