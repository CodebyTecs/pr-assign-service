package domain

import "errors"

var (
	ErrTeamExists  = errors.New("team already exists")
	ErrPRExists    = errors.New("pull request already exists")
	ErrPRMerged    = errors.New("pull request already merged")
	ErrNotAssigned = errors.New("reviewer not assigned to pull request")
	ErrNoCandidate = errors.New("no active candidate available for review")
	ErrNotFound    = errors.New("resource not found")
)
