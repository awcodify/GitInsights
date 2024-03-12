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
	"sync"
	"time"

	"github.com/google/go-github/v38/github"
	"golang.org/x/oauth2"
)

// GitHubClient is an interface for GitHub-related functionality
type GitHubClient interface {
	GetUser(ctx context.Context, user string) (*github.User, *github.Response, error)
	ListRepositories(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error)
	ListLanguages(ctx context.Context, user, repo string) (map[string]int, *github.Response, error)
	ListAllCommits(ctx context.Context, user string, opts *github.CommitsListOptions) ([]*github.RepositoryCommit, *github.Response, error)
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

// Implement the new method ListAllCommits in the GitHubAPI struct
func (g *GitHubAPI) ListAllCommits(ctx context.Context, user string, opts *github.CommitsListOptions) ([]*github.RepositoryCommit, *github.Response, error) {
	var allCommits []*github.RepositoryCommit

	repos, _, err := g.Client.Repositories.List(ctx, user, nil)
	if err != nil {
		return nil, nil, err
	}

	for _, repo := range repos {
		commits, _, err := g.Client.Repositories.ListCommits(ctx, user, *repo.Name, opts)
		if err != nil {
			return nil, nil, err
		}

		allCommits = append(allCommits, commits...)
	}

	return allCommits, nil, nil
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

	// Calculate most productive day
	mostProductiveDay, err := calculateMostProductiveDay(ctx, gitHubAPI)
	if err != nil {
		log.Printf("Error calculating most productive day: %v\n", err)
	}

	// Calculate most productive time
	mostProductiveHour, err := calculateMostProductiveTime(ctx, gitHubAPI)
	if err != nil {
		log.Printf("Error calculating most productive time: %v\n", err)
	}

	// Step 5: Generate markdown content
	markdownContent := generateMarkdown(stats, totalSize, mostProductiveDay, mostProductiveHour)

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

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Create a channel to receive language statistics from goroutines
	statsCh := make(chan map[string]int)

	// Iterate through each repository and fetch languages concurrently
	for _, repo := range repos {
		wg.Add(1)
		go func(repo *github.Repository) {
			defer wg.Done()
			languages, _, err := client.ListLanguages(ctx, *user.Login, *repo.Name)
			if err != nil {
				log.Printf("Error fetching languages for repository %s: %v\n", *repo.Name, err)
				return
			}
			statsCh <- languages
		}(repo)
	}

	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(statsCh)
	}()

	// Collect language statistics from the channel
	for stats := range statsCh {
		// Increment language count in the map
		for lang, bytes := range stats {
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

func generateMarkdown(stats map[string]int, totalSize int, mostProductiveDay string, mostProductiveHour string) string {
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
	lines = append(lines, "\n📅 Most Productive Day: "+mostProductiveDay)
	lines = append(lines, "\n⌚️ Most Productive Hour: "+mostProductiveHour)
	lines = append(lines, "\n _Last update: "+fmt.Sprint(time.Now().Format("2006-01-02 15:04:05")+"_"))
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

func calculateMostProductiveDay(ctx context.Context, client GitHubClient) (string, error) {
	user, _, err := client.GetUser(ctx, "")
	if err != nil {
		return "", fmt.Errorf("failed to get user details: %w", err)
	}

	commits, _, err := client.ListAllCommits(ctx, *user.Name, nil)
	if err != nil {
		return "", err
	}

	commitsPerDay := make(map[string]int)
	var wg sync.WaitGroup
	ch := make(chan string)

	// Process commits concurrently
	for _, commit := range commits {
		wg.Add(1)
		go func(c *github.RepositoryCommit) {
			defer wg.Done()
			day := c.Commit.Author.Date.Weekday().String()
			ch <- day
		}(commit)
	}

	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Collect results from goroutines
	for day := range ch {
		commitsPerDay[day]++
	}

	mostProductiveDay := findMaxKey(commitsPerDay)
	return mostProductiveDay, nil
}

func findMaxKey(data map[string]int) string {
	maxKey := ""
	maxVal := 0
	for key, val := range data {
		if val > maxVal {
			maxKey = key
			maxVal = val
		}
	}
	return maxKey
}

func calculateMostProductiveTime(ctx context.Context, client GitHubClient) (string, error) {
	user, _, err := client.GetUser(ctx, "")
	if err != nil {
		return "", fmt.Errorf("failed to get user details: %w", err)
	}

	commits, _, err := client.ListAllCommits(ctx, *user.Name, nil)
	if err != nil {
		return "", err
	}

	commitsPerHour := make(map[int]int)
	var wg sync.WaitGroup
	ch := make(chan int)

	// Process commits concurrently
	for _, commit := range commits {
		wg.Add(1)
		go func(c *github.RepositoryCommit) {
			defer wg.Done()
			hour := c.Commit.Author.Date.Hour()
			ch <- hour
		}(commit)
	}

	// Close the channel after all goroutines are done
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Collect results from goroutines
	for hour := range ch {
		commitsPerHour[hour]++
	}

	mostProductiveHour := findMaxKeyInt(commitsPerHour)

	// Determine the time range based on the most productive hour
	startHour := mostProductiveHour % 24
	endHour := (startHour + 1) % 24

	return fmt.Sprintf("%02d:00 - %02d:00", startHour, endHour), nil
}

func findMaxKeyInt(data map[int]int) int {
	maxKey := 0
	maxVal := 0
	for key, val := range data {
		if val > maxVal {
			maxKey = key
			maxVal = val
		}
	}
	return maxKey
}
