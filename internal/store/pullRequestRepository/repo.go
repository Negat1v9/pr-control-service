package pullrequestrepository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Negat1v9/pr-review-service/internal/models"
	"github.com/jmoiron/sqlx"
)

type pullRequestRepository struct{}

func NewPullRequestRepository() *pullRequestRepository {
	return &pullRequestRepository{}
}

func (r *pullRequestRepository) CreatePullRequest(ctx context.Context, exec sqlx.ExtContext, pr *models.PullRequest) error {
	err := exec.QueryRowxContext(ctx, createPullRequestQuery, pr.ID, pr.Name, pr.AuthorID).
		StructScan(pr)
	return err
}

func (r *pullRequestRepository) GetPullRequestByID(ctx context.Context, exec sqlx.ExtContext, prID string) (*models.PullRequest, error) {
	var pr models.PullRequest
	// recieve PR data
	if err := exec.QueryRowxContext(ctx, getPullRequestByIDQuery, prID).
		StructScan(&pr); err != nil {
		return nil, err
	}

	// then recieve reviewers
	rows, err := exec.QueryxContext(ctx, getPullRequestReviewersQuery, prID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviewers := make([]string, 0)
	for rows.Next() {
		var reviewerID string
		if err := rows.Scan(&reviewerID); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, reviewerID)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	pr.AssignedReviewers = reviewers
	return &pr, nil
}

func (r *pullRequestRepository) MergePullRequest(ctx context.Context, exec sqlx.ExtContext, prID string) error {
	res, err := exec.ExecContext(ctx, mergePullRequestQuery, prID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *pullRequestRepository) AssignReviewer(ctx context.Context, exec sqlx.ExtContext, prID, reviewerID string) error {
	_, err := exec.ExecContext(ctx, createAssignedQuery, reviewerID, prID)
	return err
}

func (r *pullRequestRepository) AssignManyReviewers(ctx context.Context, exec sqlx.ExtContext, prID string, reviewerIDs []string) error {
	if len(reviewerIDs) == 0 {
		return nil
	}

	var placeholders []string
	var args []any

	for i, reviewerID := range reviewerIDs {
		offset := i * 2
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", offset+1, offset+2))
		args = append(args, reviewerID, prID)
	}

	query := fmt.Sprintf(createManyAssignedQuery, strings.Join(placeholders, ","))

	_, err := exec.ExecContext(ctx, query, args...)
	return err
}

func (r *pullRequestRepository) DeleteAssignedByReviewerID(ctx context.Context, exec sqlx.ExtContext, reviewerID string) error {
	res, err := exec.ExecContext(ctx, deleteAssignedByReviewerIDQuery, reviewerID)

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
