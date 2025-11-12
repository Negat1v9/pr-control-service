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
	CreatedAt         time.Time         `json:"createdAt" db:"createdAt"`
	MergerAt          *time.Time        `json:"mergedAt,omitempty" db:"mergedAt,omitempty"`
}
