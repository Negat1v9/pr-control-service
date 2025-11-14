package models

import "time"

type PullRequestStatus string

const (
	PullRequestStatusOpen   PullRequestStatus = "OPEN"
	PullRequestStatusMerged PullRequestStatus = "MERGED"
)

type PullRequest struct {
	ID                string            `json:"pull_request_id" db:"pull_request_id"`
	Name              string            `json:"pull_request_name" db:"pull_request_name"`
	AuthorID          string            `json:"author_id" db:"author_id"`
	Status            PullRequestStatus `json:"status" db:"status"`
	AssignedReviewers []string          `json:"assigned_reviewers" db:"assigned_reviewers"`
	CreatedAt         time.Time         `json:"-" db:"created_at"`
	MergerAt          *time.Time        `json:"mergedAt,omitempty" db:"merged_at,omitempty"`
}

// CreatePullRequest represents the data needed to create a new pull request.
type CreatePullRequest struct {
	ID       string `json:"pull_request_id" db:"pull_request_id"`
	Name     string `json:"pull_request_name" db:"pull_request_name"`
	AuthorID string `json:"author_id" db:"author_id"`
}

type MergePullRequest struct {
	ID string `json:"pull_request_id"`
}

type ReassignPullRequest struct {
	ID            string `json:"pull_request_id"`
	OldReviewerID string `json:"old_reviewer_id"`
}

type ReassignPullRequestResponse struct {
	PR        PullRequest `json:"pr"`
	RepacedBy string      `json:"replaced_by"`
}
type UserReviews struct {
	UserID       string        `json:"user_id" db:"user_id"`
	PullRequests []PullRequest `json:"pull_requests" db:"pull_requests"`
}
