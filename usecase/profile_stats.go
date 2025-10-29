package usecase

import (
	"context"
	"fmt"
	"sort"
	"time"

	"GitInsights/domain"
)

// ProfileStatsUseCase orchestrates business logic for profile statistics
type ProfileStatsUseCase struct {
	githubRepo          domain.GitHubRepository
	maxVisibleLanguages int
}

// NewProfileStatsUseCase creates a new instance
func NewProfileStatsUseCase(githubRepo domain.GitHubRepository, maxVisibleLanguages int) *ProfileStatsUseCase {
	return &ProfileStatsUseCase{
		githubRepo:          githubRepo,
		maxVisibleLanguages: maxVisibleLanguages,
	}
}

// GetProfileStats retrieves and calculates all profile statistics
func (uc *ProfileStatsUseCase) GetProfileStats(ctx context.Context) (*domain.ProfileStats, error) {
	// Get username
	username, err := uc.githubRepo.GetUsername(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get username: %w", err)
	}

	// Get language statistics
	languageMap, err := uc.githubRepo.GetLanguageStats(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get language stats: %w", err)
	}

	// Calculate total bytes and prepare language stats
	totalBytes := 0
	for _, bytes := range languageMap {
		totalBytes += bytes
	}

	languages := uc.processLanguages(languageMap, totalBytes)

	// Get commits
	commits, err := uc.githubRepo.GetAllCommits(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits: %w", err)
	}

	// Calculate productivity metrics
	mostProductiveDay := uc.calculateMostProductiveDay(commits)
	mostProductiveHour := uc.calculateMostProductiveTime(commits)

	return &domain.ProfileStats{
		Username:           username,
		Languages:          languages,
		TotalBytes:         totalBytes,
		MostProductiveDay:  mostProductiveDay,
		MostProductiveHour: mostProductiveHour,
		LastUpdated:        time.Now(),
	}, nil
}

// processLanguages sorts and combines languages, showing top N languages
func (uc *ProfileStatsUseCase) processLanguages(languageMap map[string]int, totalBytes int) []domain.LanguageStats {
	maxVisibleLanguages := uc.maxVisibleLanguages

	// Sort languages by bytes
	type langPair struct {
		name  string
		bytes int
	}

	var pairs []langPair
	for lang, bytes := range languageMap {
		pairs = append(pairs, langPair{lang, bytes})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].bytes > pairs[j].bytes
	})

	// If total languages <= max count, show all
	if len(pairs) <= maxVisibleLanguages {
		var result []domain.LanguageStats
		for _, pair := range pairs {
			percentage := float64(pair.bytes) / float64(totalBytes) * 100
			result = append(result, domain.LanguageStats{
				Language:   pair.name,
				Bytes:      pair.bytes,
				Percentage: percentage,
			})
		}
		return result
	}

	// Otherwise, show top N languages and group rest into "Other"
	var result []domain.LanguageStats
	otherBytes := 0

	for i, pair := range pairs {
		percentage := float64(pair.bytes) / float64(totalBytes) * 100
		if i < maxVisibleLanguages {
			result = append(result, domain.LanguageStats{
				Language:   pair.name,
				Bytes:      pair.bytes,
				Percentage: percentage,
			})
		} else {
			otherBytes += pair.bytes
		}
	}

	// Add "Other" category
	if otherBytes > 0 {
		result = append(result, domain.LanguageStats{
			Language:   "Other",
			Bytes:      otherBytes,
			Percentage: float64(otherBytes) / float64(totalBytes) * 100,
		})
	}

	return result
}

// calculateMostProductiveDay finds the weekday with most commits
func (uc *ProfileStatsUseCase) calculateMostProductiveDay(commits []domain.Commit) string {
	if len(commits) == 0 {
		return "N/A"
	}

	dayCount := make(map[string]int)
	for _, commit := range commits {
		day := commit.Date.Weekday().String()
		dayCount[day]++
	}

	return findMaxKey(dayCount)
}

// calculateMostProductiveTime finds the hour with most commits
func (uc *ProfileStatsUseCase) calculateMostProductiveTime(commits []domain.Commit) string {
	if len(commits) == 0 {
		return "N/A"
	}

	hourCount := make(map[int]int)
	for _, commit := range commits {
		hour := commit.Date.Hour()
		hourCount[hour]++
	}

	mostProductiveHour := findMaxKeyInt(hourCount)
	startHour := mostProductiveHour % 24
	endHour := (startHour + 1) % 24

	return fmt.Sprintf("%02d:00 - %02d:00", startHour, endHour)
}

// Helper functions
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
