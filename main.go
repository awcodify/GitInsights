package main

import (
	"context"
	"flag"
	"log"
	"os"

	"GitInsights/infrastructure"
	"GitInsights/presentation"
	"GitInsights/usecase"
)

func main() {
	// Parse command-line flags
	includeForks := flag.Bool("include-forks", false, "Include forked repositories in analysis")
	maxVisibleLanguages := flag.Int("max-visible-language", 10, "Maximum number of languages to display (rest grouped as 'Other')")
	flag.Parse()

	// Get GitHub token from environment
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable is not set")
	}

	// Initialize dependencies
	ctx := context.Background()
	githubClient := infrastructure.NewGitHubClient(token, *includeForks)
	fileManager := infrastructure.NewFileManager("README.md")

	// Initialize use case
	profileUseCase := usecase.NewProfileStatsUseCase(githubClient, *maxVisibleLanguages)

	// Initialize presentation layer
	markdownGen := presentation.NewMarkdownGenerator()

	// Execute business logic
	stats, err := profileUseCase.GetProfileStats(ctx)
	if err != nil {
		log.Fatalf("Failed to get profile stats: %v", err)
	}

	// Generate output
	markdown := markdownGen.Generate(stats)

	// Update README
	if err := fileManager.UpdateReadme(markdown); err != nil {
		log.Fatalf("Failed to update README: %v", err)
	}

	log.Println("✅ Successfully updated README.md with Git Insights!")
}
