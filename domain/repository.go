package domain

import (
	"context"
	"time"
)

// UserProfile represents basic user information
type UserProfile struct {
	Username  string
	CreatedAt time.Time
}

// GitHubRepository defines the interface for GitHub data access
type GitHubRepository interface {
	GetUsername(ctx context.Context) (string, error)
	GetUserProfile(ctx context.Context) (*UserProfile, error)
	GetLanguageStats(ctx context.Context, username string) (map[string]int, error)
	GetAllCommits(ctx context.Context, username string) ([]Commit, error)
}

// FileRepository defines the interface for file operations
type FileRepository interface {
	UpdateReadme(content string) error
}
