package userservice

import (
	"context"
	"database/sql"

	"github.com/Negat1v9/pr-review-service/internal/models"
	"github.com/Negat1v9/pr-review-service/internal/store"
	"github.com/Negat1v9/pr-review-service/pkg/utils"
)

type UserService struct {
	store store.Store
}

func NewUserService(store store.Store) *UserService {
	return &UserService{
		store: store,
	}
}

func (s *UserService) SetUserActiveStatus(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	updatedUser, err := s.store.UserRepo().UpdateUserStatus(ctx, s.store.DB(), userID, isActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NewNotFoundError("resource not found", nil)
		}
		return nil, err
	}
	return updatedUser, nil
}

func (s *UserService) GetReview(ctx context.Context, userID string) (*models.UserReviews, error) {
	userReviews, err := s.store.UserRepo().GetUserReviews(ctx, s.store.DB(), userID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NewNotFoundError("resource not found", nil)
		}
		return nil, err
	}
	return userReviews, nil
}
