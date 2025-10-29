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
	client       *github.Client
	includeForks bool
}

// NewGitHubClient creates a new GitHub client
func NewGitHubClient(token string, includeForks bool) *GitHubClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	return &GitHubClient{
		client:       github.NewClient(tc),
		includeForks: includeForks,
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

// filterRepositories filters out forks if includeForks is false
func (g *GitHubClient) filterRepositories(repos []*github.Repository) []*github.Repository {
	if g.includeForks {
		return repos
	}

	var filtered []*github.Repository
	for _, repo := range repos {
		if repo.Fork == nil || !*repo.Fork {
			filtered = append(filtered, repo)
		}
	}
	return filtered
}

// GetLanguageStats aggregates language statistics across all repositories
func (g *GitHubClient) GetLanguageStats(ctx context.Context, username string) (map[string]int, error) {
	// Get all repositories with pagination
	var allRepos []*github.Repository
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for {
		repos, resp, err := g.client.Repositories.List(ctx, username, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list repositories: %w", err)
		}
		allRepos = append(allRepos, repos...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	// Filter out forks if needed
	allRepos = g.filterRepositories(allRepos)

	log.Printf("Analyzing languages across %d repositories%s...\n",
		len(allRepos),
		func() string {
			if !g.includeForks {
				return " (excluding forks)"
			}
			return ""
		}())

	languageStats := make(map[string]int)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Fetch languages for each repository concurrently
	for _, repo := range allRepos {
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
	// Get all repositories with pagination
	var allRepos []*github.Repository
	repoOpts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for {
		repos, resp, err := g.client.Repositories.List(ctx, username, repoOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to list repositories: %w", err)
		}
		allRepos = append(allRepos, repos...)

		if resp.NextPage == 0 {
			break
		}
		repoOpts.Page = resp.NextPage
	}

	// Filter out forks if needed
	allRepos = g.filterRepositories(allRepos)

	var allCommits []domain.Commit
	var mu sync.Mutex
	var wg sync.WaitGroup

	log.Printf("Fetching commits from %d repositories%s...\n",
		len(allRepos),
		func() string {
			if !g.includeForks {
				return " (excluding forks)"
			}
			return ""
		}())

	// Fetch commits for each repository concurrently
	for _, repo := range allRepos {
		wg.Add(1)
		go func(repoName string) {
			defer wg.Done()

			// Fetch all commits with pagination
			opts := &github.CommitsListOptions{
				ListOptions: github.ListOptions{
					PerPage: 100, // Maximum allowed by GitHub API
				},
			}

			repoCommitCount := 0
			for {
				commits, resp, err := g.client.Repositories.ListCommits(ctx, username, repoName, opts)
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
						repoCommitCount++
					}
				}
				mu.Unlock()

				// Check if there are more pages
				if resp.NextPage == 0 {
					break
				}
				opts.Page = resp.NextPage
			}

			if repoCommitCount > 0 {
				log.Printf("  ✓ %s: %d commits\n", repoName, repoCommitCount)
			}
		}(*repo.Name)
	}

	wg.Wait()
	log.Printf("Total commits analyzed: %d\n", len(allCommits))
	return allCommits, nil
}
