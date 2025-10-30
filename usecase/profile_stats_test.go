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
	UserProfile   *domain.UserProfile
	LanguageStats map[string]int
	Commits       []domain.Commit
	Err           error
}

func (m *MockGitHubRepository) GetUsername(ctx context.Context) (string, error) {
	return m.Username, m.Err
}

func (m *MockGitHubRepository) GetUserProfile(ctx context.Context) (*domain.UserProfile, error) {
	return m.UserProfile, m.Err
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
		UserProfile: &domain.UserProfile{
			Username:  "testuser",
			CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		LanguageStats: map[string]int{
			"Go":   1000,
			"Java": 500,
		},
		Commits: []domain.Commit{
			{Date: time.Date(2023, 11, 12, 10, 0, 0, 0, time.UTC)},
			{Date: time.Date(2023, 11, 13, 14, 0, 0, 0, time.UTC)},
		},
	}

	uc := usecase.NewProfileStatsUseCase(mockRepo, 10, "")
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

func TestExcludeLanguages(t *testing.T) {
	mockRepo := &MockGitHubRepository{
		Username: "testuser",
		UserProfile: &domain.UserProfile{
			Username:  "testuser",
			CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		LanguageStats: map[string]int{
			"Go":   1000,
			"Java": 500,
			"SCSS": 300,
			"HTML": 200,
			"CSS":  100,
		},
		Commits: []domain.Commit{
			{Date: time.Date(2023, 11, 12, 10, 0, 0, 0, time.UTC)},
		},
	}

	// Test excluding SCSS and HTML (case-insensitive)
	uc := usecase.NewProfileStatsUseCase(mockRepo, 10, "scss,html")
	stats, err := uc.GetProfileStats(context.Background())

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check that SCSS and HTML are not in the languages list
	for _, lang := range stats.Languages {
		if lang.Language == "SCSS" || lang.Language == "HTML" {
			t.Errorf("Expected %s to be excluded, but it's in the results", lang.Language)
		}
	}

	// Check that other languages are present
	expectedLanguages := map[string]bool{"Go": false, "Java": false, "CSS": false}
	for _, lang := range stats.Languages {
		if _, ok := expectedLanguages[lang.Language]; ok {
			expectedLanguages[lang.Language] = true
		}
	}

	for lang, found := range expectedLanguages {
		if !found {
			t.Errorf("Expected %s to be in the results, but it's missing", lang)
		}
	}

	// Total bytes should be 1600 (1000 + 500 + 100), excluding SCSS (300) and HTML (200)
	if stats.TotalBytes != 2100 {
		t.Errorf("Expected total bytes 2100 (original total), got: %d", stats.TotalBytes)
	}

	// Check that percentages add up to ~100%
	totalPercentage := 0.0
	for _, lang := range stats.Languages {
		totalPercentage += lang.Percentage
	}
	if totalPercentage < 99.9 || totalPercentage > 100.1 {
		t.Errorf("Expected total percentage to be ~100%%, got: %.2f%%", totalPercentage)
	}
}

func TestExcludeLanguagesCaseInsensitive(t *testing.T) {
	mockRepo := &MockGitHubRepository{
		Username: "testuser",
		UserProfile: &domain.UserProfile{
			Username:  "testuser",
			CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		LanguageStats: map[string]int{
			"Go":   1000,
			"SCSS": 500,
		},
		Commits: []domain.Commit{
			{Date: time.Date(2023, 11, 12, 10, 0, 0, 0, time.UTC)},
		},
	}

	// Test case-insensitive matching (lowercase input, uppercase language)
	uc := usecase.NewProfileStatsUseCase(mockRepo, 10, "scss")
	stats, err := uc.GetProfileStats(context.Background())

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// SCSS should be excluded
	for _, lang := range stats.Languages {
		if lang.Language == "SCSS" {
			t.Errorf("Expected SCSS to be excluded (case-insensitive), but it's in the results")
		}
	}

	// Only Go should remain
	if len(stats.Languages) != 1 {
		t.Errorf("Expected 1 language, got: %d", len(stats.Languages))
	}

	if stats.Languages[0].Language != "Go" {
		t.Errorf("Expected Go to be the only language, got: %s", stats.Languages[0].Language)
	}
}

func TestExcludeLanguagesWithSpaces(t *testing.T) {
	mockRepo := &MockGitHubRepository{
		Username: "testuser",
		UserProfile: &domain.UserProfile{
			Username:  "testuser",
			CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		LanguageStats: map[string]int{
			"Go":   1000,
			"Java": 500,
			"HTML": 200,
		},
		Commits: []domain.Commit{
			{Date: time.Date(2023, 11, 12, 10, 0, 0, 0, time.UTC)},
		},
	}

	// Test with spaces around commas
	uc := usecase.NewProfileStatsUseCase(mockRepo, 10, "html , java ")
	stats, err := uc.GetProfileStats(context.Background())

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Only Go should remain
	if len(stats.Languages) != 1 {
		t.Errorf("Expected 1 language, got: %d", len(stats.Languages))
	}

	if stats.Languages[0].Language != "Go" {
		t.Errorf("Expected Go to be the only language, got: %s", stats.Languages[0].Language)
	}
}

func TestNoExcludeLanguages(t *testing.T) {
	mockRepo := &MockGitHubRepository{
		Username: "testuser",
		UserProfile: &domain.UserProfile{
			Username:  "testuser",
			CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		LanguageStats: map[string]int{
			"Go":   1000,
			"Java": 500,
		},
		Commits: []domain.Commit{
			{Date: time.Date(2023, 11, 12, 10, 0, 0, 0, time.UTC)},
		},
	}

	// Test with empty exclude string
	uc := usecase.NewProfileStatsUseCase(mockRepo, 10, "")
	stats, err := uc.GetProfileStats(context.Background())

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// All languages should be present
	if len(stats.Languages) != 2 {
		t.Errorf("Expected 2 languages, got: %d", len(stats.Languages))
	}
}
