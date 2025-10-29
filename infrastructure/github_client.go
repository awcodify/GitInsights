package infrastructure

import (
	"context"
	"fmt"
	"log"
	"sync"

	"GitInsights/domain"

	"github.com/google/go-github/v38/github"
	"golang.org/x/oauth2"
)

// GitHubClient implements domain.GitHubRepository
type GitHubClient struct {
	client *github.Client
}

// NewGitHubClient creates a new GitHub client
func NewGitHubClient(token string) *GitHubClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	return &GitHubClient{
		client: github.NewClient(tc),
	}
}

// GetUsername retrieves the authenticated user's username
func (g *GitHubClient) GetUsername(ctx context.Context) (string, error) {
	user, _, err := g.client.Users.Get(ctx, "")
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	if user.Login == nil {
		return "", fmt.Errorf("user login is nil")
	}

	return *user.Login, nil
}

// GetLanguageStats aggregates language statistics across all repositories
func (g *GitHubClient) GetLanguageStats(ctx context.Context, username string) (map[string]int, error) {
	repos, _, err := g.client.Repositories.List(ctx, username, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories: %w", err)
	}

	languageStats := make(map[string]int)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Fetch languages for each repository concurrently
	for _, repo := range repos {
		wg.Add(1)
		go func(repoName string) {
			defer wg.Done()

			languages, _, err := g.client.Repositories.ListLanguages(ctx, username, repoName)
			if err != nil {
				log.Printf("Error fetching languages for %s: %v\n", repoName, err)
				return
			}

			mu.Lock()
			for lang, bytes := range languages {
				languageStats[lang] += bytes
			}
			mu.Unlock()
		}(*repo.Name)
	}

	wg.Wait()
	return languageStats, nil
}

// GetAllCommits retrieves all commits across all repositories
func (g *GitHubClient) GetAllCommits(ctx context.Context, username string) ([]domain.Commit, error) {
	repos, _, err := g.client.Repositories.List(ctx, username, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories: %w", err)
	}

	var allCommits []domain.Commit
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Fetch commits for each repository concurrently
	for _, repo := range repos {
		wg.Add(1)
		go func(repoName string) {
			defer wg.Done()

			commits, _, err := g.client.Repositories.ListCommits(ctx, username, repoName, nil)
			if err != nil {
				log.Printf("Error fetching commits for %s: %v\n", repoName, err)
				return
			}

			mu.Lock()
			for _, commit := range commits {
				if commit.Commit != nil && commit.Commit.Author != nil && commit.Commit.Author.Date != nil {
					allCommits = append(allCommits, domain.Commit{
						Date: *commit.Commit.Author.Date,
					})
				}
			}
			mu.Unlock()
		}(*repo.Name)
	}

	wg.Wait()
	return allCommits, nil
}
