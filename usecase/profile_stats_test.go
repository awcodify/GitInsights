package usecase_test

import (
	"context"
	"testing"
	"time"

	"GitInsights/domain"
	"GitInsights/usecase"
)

// MockGitHubRepository for testing
type MockGitHubRepository struct {
	Username      string
	LanguageStats map[string]int
	Commits       []domain.Commit
	Err           error
}

func (m *MockGitHubRepository) GetUsername(ctx context.Context) (string, error) {
	return m.Username, m.Err
}

func (m *MockGitHubRepository) GetLanguageStats(ctx context.Context, username string) (map[string]int, error) {
	return m.LanguageStats, m.Err
}

func (m *MockGitHubRepository) GetAllCommits(ctx context.Context, username string) ([]domain.Commit, error) {
	return m.Commits, m.Err
}

func TestGetProfileStats(t *testing.T) {
	mockRepo := &MockGitHubRepository{
		Username: "testuser",
		LanguageStats: map[string]int{
			"Go":   1000,
			"Java": 500,
		},
		Commits: []domain.Commit{
			{Date: time.Date(2023, 11, 12, 10, 0, 0, 0, time.UTC)},
			{Date: time.Date(2023, 11, 13, 14, 0, 0, 0, time.UTC)},
		},
	}

	uc := usecase.NewProfileStatsUseCase(mockRepo, 10)
	stats, err := uc.GetProfileStats(context.Background())

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if stats.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got: %s", stats.Username)
	}

	if stats.TotalBytes != 1500 {
		t.Errorf("Expected total bytes 1500, got: %d", stats.TotalBytes)
	}
}
