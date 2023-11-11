package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"bytes"

	"github.com/google/go-github/v38/github"
	"golang.org/x/oauth2"
)

func main() {
	// Step 1: Get GitHub token from gh auth token
	token, err := getGitHubToken()
	if err != nil {
		log.Fatal("Error getting GitHub token:", err)
	}

	// Step 2: Initialize GitHub client
	ctx := context.Background()
	client := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)))

	// Step 3: Summarize GitHub profile
	stats, totalSize, err := summarizeGitHubProfile(ctx, client)
	if err != nil {
		log.Fatal("Error summarizing GitHub profile:", err)
	}

	// Step 4: Generate markdown content
	markdownContent := generateMarkdown(stats, totalSize)

	// Step 5: Update README.md
	err = updateReadme(markdownContent)
	if err != nil {
		log.Fatal("Error updating README.md:", err)
	}
}

func getGitHubToken() (string, error) {
	// Run "gh auth token" command to get the GitHub token
	cmd := exec.Command("gh", "auth", "token")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute 'gh auth token': %w", err)
	}

	// Trim any leading/trailing whitespaces from the token
	token := strings.TrimSpace(string(output))
	return token, nil
}

func summarizeGitHubProfile(ctx context.Context, client *github.Client) (map[string]int, int, error) {
	// Get the authenticated user
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user details: %w", err)
	}

	fmt.Printf("GitHub Profile Summary for %s (%s)\n", *user.Login, *user.Name)

	// Get the list of repositories for the authenticated user
	repos, _, err := client.Repositories.List(ctx, *user.Login, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user repositories: %w", err)
	}

	// Initialize a map to store language statistics
	languageStats := make(map[string]int)

	// Iterate through each repository and count the languages
	for _, repo := range repos {
		languages, _, err := client.Repositories.ListLanguages(ctx, *user.Login, *repo.Name)
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
	lines = append(lines, "\nLanguage Statistics:\n")

	// Sort languages by percentage
	sortedLanguages := sortLanguagesByPercentage(stats, totalSize)

	// Combine languages that make up less than 5% into "Other"
	combinedLanguages := combineLanguages(sortedLanguages, stats, totalSize)

	maxLangLength := maxLanguageLength(combinedLanguages)

	for _, lang := range combinedLanguages {
		percentage := float64(stats[lang]) / float64(totalSize) * 100
		lines = append(lines, fmt.Sprintf("%-*s %s%6.2f%%", maxLangLength, lang, generateProgressBar(percentage), percentage))
	}

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
