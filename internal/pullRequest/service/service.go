package prservice

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Negat1v9/pr-review-service/internal/models"
	"github.com/Negat1v9/pr-review-service/internal/store"
	"github.com/Negat1v9/pr-review-service/pkg/utils"
	"github.com/jmoiron/sqlx"
)

type PRService struct {
	store store.Store
}

func NewPRService(store store.Store) *PRService {
	return &PRService{
		store: store,
	}
}

func (s *PRService) CreatePR(ctx context.Context, pr *models.CreatePullRequest) (*models.PullRequest, error) {
	exitstPR, err := s.store.PRRepo().GetPullRequestByID(ctx, s.store.DB(), pr.ID)
	if err == nil && exitstPR != nil {
		return nil, utils.NewError(409, utils.ErrPrExists, "PR id already exists", nil)
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("CreatePR: unable to get pull request by ID: %v", err)
	}
	newPr := &models.PullRequest{
		ID:        pr.ID,
		Name:      pr.Name,
		AuthorID:  pr.AuthorID,
		CreatedAt: time.Now(),
	}

	err = s.store.DoTx(ctx, func(ctx context.Context, exec sqlx.ExtContext) error {
		// get active team members of PR author to assign as reviewers
		activeAuthorsTeamMembers, err := s.store.TeamRepo().GetUsersIDFromUserTeam(ctx, exec, pr.AuthorID, 2)
		if err != nil {
			if err == sql.ErrNoRows {
				// user with
				return utils.NewNotFoundError("resource not found", nil)
			}
			return fmt.Errorf("CreatePR: unable to get active team members of PR author: %v", err)
		}

		// create PR
		if err := s.store.PRRepo().CreatePullRequest(ctx, exec, newPr); err != nil {
			return fmt.Errorf("CreatePR: unable to create PR: %v", err)
		}

		// assign only if there are active members in author's team
		if len(activeAuthorsTeamMembers) > 0 {
			if err := s.store.PRRepo().AssignManyReviewers(ctx, exec, pr.ID, activeAuthorsTeamMembers); err != nil {
				return fmt.Errorf("CreatePR: unable to assign reviewers to PR: %v", err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// receive created PR
	createdPR, err := s.store.PRRepo().GetPullRequestByID(ctx, s.store.DB(), pr.ID)
	if err != nil {
		return nil, fmt.Errorf("CreatePR: unable to get created PR: %v", err)
	}
	return createdPR, nil
}

func (s *PRService) MergePR(ctx context.Context, prID string) (*models.PullRequest, error) {
	pr, err := s.store.PRRepo().GetPullRequestByID(ctx, s.store.DB(), prID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NewNotFoundError("resource not found", nil)
		}
	}

	// pr alredy merged not merge it twice
	if pr.Status == models.PullRequestStatusMerged {
		return pr, nil
	}

	err = s.store.PRRepo().MergePullRequest(ctx, s.store.DB(), prID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NewNotFoundError("resource not found", nil)
		}
		return nil, fmt.Errorf("MergePR: unable to merge PR: %v", err)
	}

	updatedPR, err := s.store.PRRepo().GetPullRequestByID(ctx, s.store.DB(), prID)
	if err != nil {
		return nil, fmt.Errorf("MergePR: unable to get updated PR: %v", err)
	}

	return updatedPR, nil
}

func (s *PRService) ReassignPR(ctx context.Context, prID string, oldReviewerID string) (*models.ReassignPullRequestResponse, error) {
	pr, err := s.store.PRRepo().GetPullRequestByID(ctx, s.store.DB(), prID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NewNotFoundError("resource not found", nil)
		}
	}
	// pr alredy merged not merge
	if pr.Status == models.PullRequestStatusMerged {
		return nil, utils.NewError(409, utils.ErrPrAlredyMerged, "cannot reassign on merged PR", nil)
	}

	// user not reviewer
	if !isUserIDInReviewers(oldReviewerID, pr.AssignedReviewers) {
		return nil, utils.NewError(409, utils.ErrUserNotReviewer, "reviewer is not assigned to this PR", nil)
	}

	// get active users from author team without oldest
	newActiveUsers, err := s.store.TeamRepo().GetActiveUsersTeamWithException(ctx, s.store.DB(), pr.AuthorID, pr.AssignedReviewers, 1)
	if err != nil {
		fmt.Println("no newActiveUsers")
		if err == sql.ErrNoRows {
			utils.NewError(409, utils.ErrNoCantidate, "no active replacement candidate in team", nil)
		}
		return nil, err
	}

	err = s.store.DoTx(ctx, func(ctx context.Context, exec sqlx.ExtContext) error {
		// delete old
		txErr := s.store.PRRepo().DeleteAssignedByReviewerID(ctx, exec, oldReviewerID)
		if err != nil {
			return txErr
		}
		txErr = s.store.PRRepo().AssignReviewer(ctx, exec, prID, newActiveUsers[0])
		if err != nil {
			return txErr
		}
		return txErr
	})

	if err != nil {
		return nil, err
	}

	pr, err = s.store.PRRepo().GetPullRequestByID(ctx, s.store.DB(), prID)
	if err != nil {
		return nil, err
	}
	return &models.ReassignPullRequestResponse{
		PR:        *pr,
		RepacedBy: newActiveUsers[0],
	}, nil
}

func isUserIDInReviewers(userID string, reviewers []string) bool {
	for _, reviewer := range reviewers {
		if reviewer == userID {
			return true
		}
	}
	return false
}
