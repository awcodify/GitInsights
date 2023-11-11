package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"sort"
	"strings"

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
	if err := summarizeGitHubProfile(ctx, client); err != nil {
		log.Fatal("Error summarizing GitHub profile:", err)
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

func summarizeGitHubProfile(ctx context.Context, client *github.Client) error {
	// Get the authenticated user
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to get user details: %w", err)
	}

	fmt.Printf("GitHub Profile Summary for %s (%s)\n", *user.Login, *user.Name)

	// Get the list of repositories for the authenticated user
	repos, _, err := client.Repositories.List(ctx, *user.Login, nil)
	if err != nil {
		return fmt.Errorf("failed to get user repositories: %w", err)
	}

	// Initialize a map to store language statistics
	languageStats := make(map[string]int)

	// Iterate through each repository and count the languages
	for _, repo := range repos {
		languages, _, err := client.Repositories.ListLanguages(ctx, *user.Login, *repo.Name)
		if err != nil {
			return fmt.Errorf("failed to get languages for repository %s: %w", *repo.Name, err)
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

	// Sort languages by percentage
	sortedLanguages := sortLanguagesByPercentage(languageStats, totalSize)

	// Combine languages that make up less than 5% into "Other"
	combinedLanguages := combineLanguages(sortedLanguages, languageStats, totalSize)

	// Display language statistics with percentage numbers and bars
	fmt.Println("\nLanguage Statistics:")
	for _, lang := range combinedLanguages {
		percentage := float64(languageStats[lang]) / float64(totalSize) * 100
		fmt.Printf("%-25s%10.2f%%   %s\n", lang, percentage, generateProgressBar(percentage))
	}

	return nil
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

