package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/google/go-github/v38/github"
	"golang.org/x/oauth2"
)

// GitHubClient is an interface for GitHub-related functionality
type GitHubClient interface {
	GetUser(ctx context.Context, user string) (*github.User, *github.Response, error)
	ListRepositories(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error)
	ListLanguages(ctx context.Context, user, repo string) (map[string]int, *github.Response, error)
}

// GitHubAPI implements GitHubClient interface using go-github library
type GitHubAPI struct {
	Client *github.Client
}

// GetUser retrieves user details from GitHub
func (g *GitHubAPI) GetUser(ctx context.Context, user string) (*github.User, *github.Response, error) {
	return g.Client.Users.Get(ctx, user)
}

// ListRepositories retrieves a list of repositories for a user from GitHub
func (g *GitHubAPI) ListRepositories(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error) {
	return g.Client.Repositories.List(ctx, user, opts)
}

// ListLanguages retrieves the languages used in a repository from GitHub
func (g *GitHubAPI) ListLanguages(ctx context.Context, user, repo string) (map[string]int, *github.Response, error) {
	return g.Client.Repositories.ListLanguages(ctx, user, repo)
}

func main() {
	// Step 1: Get GitHub token
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN is not set")
	}

	// Step 2: Initialize GitHub client
	ctx := context.Background()
	client := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)))

	// Step 3: Create an instance of GitHubAPI
	gitHubAPI := &GitHubAPI{Client: client}

	// Step 4: Summarize GitHub profile
	stats, totalSize, err := summarizeGitHubProfile(ctx, gitHubAPI)
	if err != nil {
		log.Fatal("Error summarizing GitHub profile:", err)
	}

	// Step 5: Generate markdown content
	markdownContent := generateMarkdown(stats, totalSize)

	// Step 6: Update README.md
	err = updateReadme(markdownContent)
	if err != nil {
		log.Fatal("Error updating README.md:", err)
	}
}

func summarizeGitHubProfile(ctx context.Context, client GitHubClient) (map[string]int, int, error) {
	// Get the authenticated user
	user, _, err := client.GetUser(ctx, "")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user details: %w", err)
	}

	fmt.Printf("GitHub Profile Summary for %s (%s)\n", *user.Login, *user.Name)

	// Get the list of repositories for the authenticated user
	repos, _, err := client.ListRepositories(ctx, *user.Login, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user repositories: %w", err)
	}

	// Initialize a map to store language statistics
	languageStats := make(map[string]int)

	// Iterate through each repository and count the languages
	for _, repo := range repos {
		languages, _, err := client.ListLanguages(ctx, *user.Login, *repo.Name)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get languages for repository %s: %w", *repo.Name, err)
		}

		// Increment language count in the map
		for lang, bytes := range languages {
			languageStats[lang] += bytes
		}
	}

	// Calculate the total size of all repositories
	totalSize := 0
	for _, size := range languageStats {
		totalSize += size
	}

	return languageStats, totalSize, nil
}

func sortLanguagesByPercentage(languageStats map[string]int, totalSize int) []string {
	var sortedLanguages []string
	for lang := range languageStats {
		sortedLanguages = append(sortedLanguages, lang)
	}

	// Sort languages by percentage in descending order
	sort.Slice(sortedLanguages, func(i, j int) bool {
		return float64(languageStats[sortedLanguages[i]])/float64(totalSize) > float64(languageStats[sortedLanguages[j]])/float64(totalSize)
	})

	return sortedLanguages
}

func combineLanguages(sortedLanguages []string, languageStats map[string]int, totalSize int) []string {
	var combinedLanguages []string
	threshold := 5.0

	// Iterate through sortedLanguages and combine those < 5% into "Other"
	for _, lang := range sortedLanguages {
		percentage := float64(languageStats[lang]) / float64(totalSize) * 100
		if percentage >= threshold {
			combinedLanguages = append(combinedLanguages, lang)
		} else {
			languageStats["Other"] += languageStats[lang]
		}
	}

	return append(combinedLanguages, "Other")
}

func generateProgressBar(percentage float64) string {
	const barWidth = 40
	numFilled := int(percentage / 100 * barWidth)
	return fmt.Sprintf("%s%s", strings.Repeat("█", numFilled), strings.Repeat("░", barWidth-numFilled))
}

func generateMarkdown(stats map[string]int, totalSize int) string {
	var lines []string
	lines = append(lines, "<!--START_SECTION:GitInsights-->")
	lines = append(lines, "### Git Insight")
	lines = append(lines, "\nLanguage Statistics:")
	lines = append(lines, "```")

	// Sort languages by percentage
	sortedLanguages := sortLanguagesByPercentage(stats, totalSize)

	// Combine languages that make up less than 5% into "Other"
	combinedLanguages := combineLanguages(sortedLanguages, stats, totalSize)

	maxLangLength := maxLanguageLength(combinedLanguages)

	for _, lang := range combinedLanguages {
		percentage := float64(stats[lang]) / float64(totalSize) * 100
		lines = append(lines, fmt.Sprintf("%-*s [%-30s] %5.2f%%", maxLangLength, lang, generateProgressBar(percentage), percentage))
	}

	lines = append(lines, "```")
	lines = append(lines, "<!--END_SECTION:GitInsights-->")

	return strings.Join(lines, "\n")
}

func maxLanguageLength(languages []string) int {
	maxLength := 0
	for _, lang := range languages {
		if len(lang) > maxLength {
			maxLength = len(lang)
		}
	}
	return maxLength
}

func updateReadme(content string) error {
	filePath := "README.md"
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open README.md: %w", err)
	}
	defer file.Close()

	// Read the existing content
	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read README.md: %w", err)
	}

	// Find the GitInsights section in the README.md content
	startIdx := bytes.Index(data, []byte("<!--START_SECTION:GitInsights-->"))
	endIdx := bytes.Index(data, []byte("<!--END_SECTION:GitInsights-->"))
	if startIdx == -1 || endIdx == -1 {
		return fmt.Errorf("GitInsights section not found in README.md")
	}

	// Replace the GitInsights section with the updated content
	updatedContent := append(data[:startIdx], []byte(content)...)
	updatedContent = append(updatedContent, data[endIdx+len("<!--END_SECTION:GitInsights-->"):]...)

	// Truncate the file and write the updated content
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate README.md: %w", err)
	}
	if _, err := file.WriteAt(updatedContent, 0); err != nil {
		return fmt.Errorf("failed to write to README.md: %w", err)
	}

	return nil
}
