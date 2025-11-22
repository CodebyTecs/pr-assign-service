package domain

import "time"

type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	ID        string     `db:"pull_request_id"   json:"pull_request_id"`
	Name      string     `db:"pull_request_name" json:"pull_request_name"`
	AuthorID  string     `db:"author_id"         json:"author_id"`
	Status    PRStatus   `db:"status"            json:"status"`
	Reviewers []string   `db:"assigned_reviewers" json:"assigned_reviewers"`
	CreatedAt *time.Time `db:"created_at"        json:"createdAt,omitempty"`
	MergedAt  *time.Time `db:"merged_at"         json:"mergedAt,omitempty"`
}

type PullRequestShort struct {
	ID       string   `json:"pull_request_id"`
	Name     string   `json:"pull_request_name"`
	AuthorID string   `json:"author_id"`
	Status   PRStatus `json:"status"`
}
