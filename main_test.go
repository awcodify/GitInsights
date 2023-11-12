package main

import (
	"context"
	"testing"

	"github.com/google/go-github/v38/github"
)

// MockGitHubAPI is a mock implementation of the GitHubClient interface
type MockGitHubAPI struct {
	UserResponse         *github.User
	ListRepositoriesResp []*github.Repository
	ListLanguagesResp    map[string]int
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
