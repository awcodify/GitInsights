package main

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-github/v38/github"
)

// MockGitHubAPI is a mock implementation of the GitHubClient interface
type MockGitHubAPI struct {
	UserResponse         *github.User
	ListRepositoriesResp []*github.Repository
	ListLanguagesResp    map[string]int
	ListAllCommitsResp   []*github.RepositoryCommit
}

func (m *MockGitHubAPI) GetUser(ctx context.Context, user string) (*github.User, *github.Response, error) {
	return m.UserResponse, nil, nil
}

func (m *MockGitHubAPI) ListRepositories(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error) {
	return m.ListRepositoriesResp, nil, nil
}

func (m *MockGitHubAPI) ListLanguages(ctx context.Context, user, repo string) (map[string]int, *github.Response, error) {
	return m.ListLanguagesResp, nil, nil
}

func (m *MockGitHubAPI) ListAllCommits(ctx context.Context, user string, opts *github.CommitsListOptions) ([]*github.RepositoryCommit, *github.Response, error) {
	return m.ListAllCommitsResp, nil, nil
}

func TestSummarizeGitHubProfile(t *testing.T) {
	// Mock data for testing
	mockUser := &github.User{
		Login: github.String("testuser"),
		Name:  github.String("Test User"),
	}

	mockRepos := []*github.Repository{
		{ID: github.Int64(1), Name: github.String("repo1")},
	}

	mockLanguages := map[string]int{
		"Go":    1000,
		"Java":  500,
		"Other": 200,
	}

	// Create an instance of the MockGitHubAPI
	mockGitHubAPI := &MockGitHubAPI{
		UserResponse:         mockUser,
		ListRepositoriesResp: mockRepos,
		ListLanguagesResp:    mockLanguages,
	}

	// Call the function to be tested
	_, totalSize, err := summarizeGitHubProfile(context.Background(), mockGitHubAPI)
	if err != nil {
		t.Errorf("Error in summarizeGitHubProfile: %v", err)
	}

	// Add assertions to verify the expected output
	if totalSize != 1700 {
		t.Errorf("Total size is incorrect, got: %d, want: %d", totalSize, 1700)
	}
}

func TestSortLanguagesByPercentage(t *testing.T) {
	// Mock data for testing
	mockStats := map[string]int{
		"Go":    1000,
		"Java":  500,
		"Other": 200,
	}

	mockTotalSize := 1700

	// Call the function to be tested
	sortedLanguages := sortLanguagesByPercentage(mockStats, mockTotalSize)

	// Add assertions to verify the expected output
	expectedOrder := []string{"Go", "Java", "Other"}
	for i, lang := range sortedLanguages {
		if lang != expectedOrder[i] {
			t.Errorf("Language order is incorrect, got: %s, want: %s", lang, expectedOrder[i])
		}
	}
}

func TestCombineLanguages(t *testing.T) {
	// Mock data for testing
	mockSortedLanguages := []string{"Go", "Java", "Other"}
	mockStats := map[string]int{
		"Go":         1000,
		"Java":       500,
		"Javascript": 10,
		"Rust":       10,
	}

	mockTotalSize := 1520

	// Call the function to be tested
	combinedLanguages := combineLanguages(mockSortedLanguages, mockStats, mockTotalSize)

	// Add assertions to verify the expected output
	expectedCombined := []string{"Go", "Java", "Other"}
	for i, lang := range combinedLanguages {
		if lang != expectedCombined[i] {
			t.Errorf("Combined languages are incorrect, got: %s, want: %s", lang, expectedCombined[i])
		}
	}
}

func TestCalculateMostProductiveDay(t *testing.T) {
	// Mock data for testing
	mockUser := &github.User{
		Login: github.String("testuser"),
		Name:  github.String("Test User"),
	}

	// Create time objects with the desired dates
	sunday := time.Date(2023, time.November, 12, 0, 0, 0, 0, time.UTC) // 12 Nov 2023 => Sunday
	monday := time.Date(2023, time.November, 13, 0, 0, 0, 0, time.UTC) // 13 Nov 2023 => Monday

	// Create mock commits with the specified dates
	// 4 m
	mockCommits := []*github.RepositoryCommit{
		{Commit: &github.Commit{Author: &github.CommitAuthor{Date: &sunday}}},
		{Commit: &github.Commit{Author: &github.CommitAuthor{Date: &sunday}}},
		{Commit: &github.Commit{Author: &github.CommitAuthor{Date: &sunday}}},
		{Commit: &github.Commit{Author: &github.CommitAuthor{Date: &sunday}}},
		{Commit: &github.Commit{Author: &github.CommitAuthor{Date: &monday}}},
	}

	// Create an instance of the MockGitHubAPI
	mockGitHubAPI := &MockGitHubAPI{
		UserResponse:         mockUser,
		ListRepositoriesResp: nil,         // Repositories list not needed for this test
		ListLanguagesResp:    nil,         // Languages list not needed for this test
		ListAllCommitsResp:   mockCommits, // Set the mock commits for ListAllCommits
	}

	// Call the function to be tested
	mostProductiveDay, err := calculateMostProductiveDay(context.Background(), mockGitHubAPI)
	if err != nil {
		t.Errorf("Error in calculateMostProductiveDay: %v", err)
	}

	// Add assertions to verify the expected output
	expectedDay := "Sunday" // Assuming the first commit is on a Monday
	if mostProductiveDay != expectedDay {
		t.Errorf("Most productive day is incorrect, got: %s, want: %s", mostProductiveDay, expectedDay)
	}
}

func TestCalculateMostProductiveTime(t *testing.T) {
	// Mock data for testing
	mockUser := &github.User{
		Login: github.String("testuser"),
		Name:  github.String("Test User"),
	}

	// Create time objects with the desired dates
	time1 := time.Date(2023, time.November, 12, 10, 0, 0, 0, time.UTC)
	time2 := time.Date(2023, time.November, 13, 0, 0, 0, 0, time.UTC)
	time3 := time.Date(2023, time.November, 14, 10, 0, 0, 0, time.UTC)
	time4 := time.Date(2023, time.November, 15, 10, 0, 0, 0, time.UTC)
	time5 := time.Date(2023, time.November, 16, 10, 0, 0, 0, time.UTC)

	// Create mock commits with the specified dates
	// 4 m
	mockCommits := []*github.RepositoryCommit{
		{Commit: &github.Commit{Author: &github.CommitAuthor{Date: &time1}}},
		{Commit: &github.Commit{Author: &github.CommitAuthor{Date: &time2}}},
		{Commit: &github.Commit{Author: &github.CommitAuthor{Date: &time3}}},
		{Commit: &github.Commit{Author: &github.CommitAuthor{Date: &time4}}},
		{Commit: &github.Commit{Author: &github.CommitAuthor{Date: &time5}}},
	}

	// Create an instance of the MockGitHubAPI
	mockGitHubAPI := &MockGitHubAPI{
		UserResponse:         mockUser,
		ListRepositoriesResp: nil, // Repositories list not needed for this test
		ListLanguagesResp:    nil, // Languages list not needed for this test
		ListAllCommitsResp:   mockCommits,
	}

	// Call the function to be tested
	mostProductiveTime, err := calculateMostProductiveTime(context.Background(), mockGitHubAPI)
	if err != nil {
		t.Errorf("Error in calculateMostProductiveTime: %v", err)
	}

	// Add assertions to verify the expected output
	expectedTime := "10:00 - 11:00" // Assuming the first commit is at 10:00
	if mostProductiveTime != expectedTime {
		t.Errorf("Most productive time is incorrect, got: %s, want: %s", mostProductiveTime, expectedTime)
	}
}
