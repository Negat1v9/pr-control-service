package userservice

import (
	"context"

	"github.com/Negat1v9/pr-review-service/internal/models"
	"github.com/Negat1v9/pr-review-service/internal/store"
)

type UserService struct {
	store *store.Store
}

func NewUserService(store *store.Store) *UserService {
	return &UserService{
		store: store,
	}
}

func (s *UserService) SetUserActiveStatus(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	updatedUser, err := s.store.UserRepo.UpdateUserStatus(ctx, s.store.Db, userID, isActive)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (s *UserService) GetReview(ctx context.Context, userID string) (*models.UserReviews, error) {

	return s.store.UserRepo.GetUserReviews(ctx, s.store.Db, userID)
}
