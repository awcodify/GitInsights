package usecase

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"GitInsights/domain"
)

// ProfileStatsUseCase orchestrates business logic for profile statistics
type ProfileStatsUseCase struct {
	githubRepo          domain.GitHubRepository
	maxVisibleLanguages int
	excludeLanguages    []string
}

// NewProfileStatsUseCase creates a new instance
func NewProfileStatsUseCase(githubRepo domain.GitHubRepository, maxVisibleLanguages int, excludeLanguagesStr string) *ProfileStatsUseCase {
	var excludeLanguages []string
	if excludeLanguagesStr != "" {
		// Split by comma and trim spaces, convert to lowercase for case-insensitive matching
		langs := strings.Split(excludeLanguagesStr, ",")
		for _, lang := range langs {
			trimmed := strings.TrimSpace(lang)
			if trimmed != "" {
				excludeLanguages = append(excludeLanguages, strings.ToLower(trimmed))
			}
		}
	}
	return &ProfileStatsUseCase{
		githubRepo:          githubRepo,
		maxVisibleLanguages: maxVisibleLanguages,
		excludeLanguages:    excludeLanguages,
	}
}

// GetProfileStats retrieves and calculates all profile statistics
func (uc *ProfileStatsUseCase) GetProfileStats(ctx context.Context) (*domain.ProfileStats, error) {
	// Get username
	username, err := uc.githubRepo.GetUsername(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get username: %w", err)
	}

	// Get user profile for account age
	userProfile, err := uc.githubRepo.GetUserProfile(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	// Calculate account age
	accountAge := uc.calculateAccountAge(userProfile.CreatedAt)

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

	// Calculate streaks
	currentStreak, longestStreak := uc.calculateStreaks(commits)

	// Calculate weekly distribution
	weeklyDistribution := uc.calculateWeeklyDistribution(commits)

	return &domain.ProfileStats{
		Username:           username,
		Languages:          languages,
		TotalBytes:         totalBytes,
		MostProductiveDay:  mostProductiveDay,
		MostProductiveHour: mostProductiveHour,
		AccountAge:         accountAge,
		CurrentStreak:      currentStreak,
		LongestStreak:      longestStreak,
		WeeklyDistribution: weeklyDistribution,
		LastUpdated:        time.Now(),
	}, nil
}

// processLanguages sorts and combines languages, showing top N languages
func (uc *ProfileStatsUseCase) processLanguages(languageMap map[string]int, totalBytes int) []domain.LanguageStats {
	// Guard against division by zero
	if totalBytes == 0 {
		return []domain.LanguageStats{}
	}

	maxVisibleLanguages := uc.maxVisibleLanguages

	// Filter out excluded languages and recalculate total bytes
	filteredLanguageMap := make(map[string]int)
	filteredTotalBytes := 0
	for lang, bytes := range languageMap {
		// Check if language should be excluded (case-insensitive)
		excluded := false
		for _, excludedLang := range uc.excludeLanguages {
			if strings.ToLower(lang) == excludedLang {
				excluded = true
				break
			}
		}
		if !excluded {
			filteredLanguageMap[lang] = bytes
			filteredTotalBytes += bytes
		}
	}

	// If all languages are excluded, return empty
	if filteredTotalBytes == 0 {
		return []domain.LanguageStats{}
	}

	// Sort languages by bytes
	type langPair struct {
		name  string
		bytes int
	}

	var pairs []langPair
	for lang, bytes := range filteredLanguageMap {
		pairs = append(pairs, langPair{lang, bytes})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].bytes > pairs[j].bytes
	})

	// If total languages <= max count, show all
	if len(pairs) <= maxVisibleLanguages {
		var result []domain.LanguageStats
		for _, pair := range pairs {
			percentage := float64(pair.bytes) / float64(filteredTotalBytes) * 100
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
		percentage := float64(pair.bytes) / float64(filteredTotalBytes) * 100
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
			Percentage: float64(otherBytes) / float64(filteredTotalBytes) * 100,
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
	if mostProductiveHour == -1 {
		return "N/A"
	}

	startHour := mostProductiveHour % 24
	endHour := (startHour + 1) % 24

	return fmt.Sprintf("%02d:00 - %02d:00", startHour, endHour)
}

// Helper functions
func findMaxKey(data map[string]int) string {
	if len(data) == 0 {
		return "N/A"
	}

	maxKey := ""
	maxVal := -1 // Use -1 to handle cases where all values are 0
	for key, val := range data {
		if val > maxVal {
			maxKey = key
			maxVal = val
		}
	}

	// If all values are 0 or negative, return N/A
	if maxVal <= 0 {
		return "N/A"
	}

	return maxKey
}

func findMaxKeyInt(data map[int]int) int {
	if len(data) == 0 {
		return -1 // Return -1 to indicate no valid hour found
	}

	maxKey := -1
	maxVal := -1 // Use -1 to handle cases where all values are 0
	for key, val := range data {
		if val > maxVal {
			maxKey = key
			maxVal = val
		}
	}

	// If all values are 0 or negative, return -1 to indicate no valid hour
	if maxVal <= 0 {
		return -1
	}

	return maxKey
}

// calculateStreaks calculates current and longest commit streaks
func (uc *ProfileStatsUseCase) calculateStreaks(commits []domain.Commit) (int, int) {
	if len(commits) == 0 {
		return 0, 0
	}

	// Sort commits by date
	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Date.Before(commits[j].Date)
	})

	// Get unique days with commits
	uniqueDays := make(map[string]bool)
	for _, commit := range commits {
		day := commit.Date.Truncate(24 * time.Hour).Format("2006-01-02")
		uniqueDays[day] = true
	}

	// Convert to sorted slice
	var days []time.Time
	for dayStr := range uniqueDays {
		day, _ := time.Parse("2006-01-02", dayStr)
		days = append(days, day)
	}
	sort.Slice(days, func(i, j int) bool {
		return days[i].Before(days[j])
	})

	if len(days) == 0 {
		return 0, 0
	}

	// Calculate streaks
	currentStreak := 0
	longestStreak := 0
	tempStreak := 1

	for i := 0; i < len(days); i++ {
		if i > 0 {
			diff := days[i].Sub(days[i-1]).Hours() / 24
			if diff == 1 {
				tempStreak++
			} else {
				if tempStreak > longestStreak {
					longestStreak = tempStreak
				}
				tempStreak = 1
			}
		}
	}

	// Check last streak
	if tempStreak > longestStreak {
		longestStreak = tempStreak
	}

	// Calculate current streak (from most recent commit)
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)
	mostRecentDay := days[len(days)-1]

	// If last commit was today or yesterday, start counting backwards
	if mostRecentDay.Equal(today) || mostRecentDay.Equal(yesterday) {
		currentStreak = 1
		for i := len(days) - 2; i >= 0; i-- {
			diff := days[i+1].Sub(days[i]).Hours() / 24
			if diff == 1 {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0
	}

	return currentStreak, longestStreak
}

// calculateWeeklyDistribution returns commit counts for each day of the week
func (uc *ProfileStatsUseCase) calculateWeeklyDistribution(commits []domain.Commit) map[string]int {
	distribution := map[string]int{
		"Monday":    0,
		"Tuesday":   0,
		"Wednesday": 0,
		"Thursday":  0,
		"Friday":    0,
		"Saturday":  0,
		"Sunday":    0,
	}

	for _, commit := range commits {
		day := commit.Date.Weekday().String()
		distribution[day]++
	}

	return distribution
}

// calculateAccountAge calculates how long the account has been active
func (uc *ProfileStatsUseCase) calculateAccountAge(createdAt time.Time) string {
	now := time.Now()
	years := now.Year() - createdAt.Year()
	months := int(now.Month()) - int(createdAt.Month())

	// Adjust months if the current day is before the creation day
	if now.Day() < createdAt.Day() {
		months--
	}

	// If months negative, adjust years and months
	if months < 0 {
		years--
		months += 12
	}

	if years > 0 {
		if months > 0 {
			return fmt.Sprintf("%d years %d months", years, months)
		}
		return fmt.Sprintf("%d years", years)
	}

	if months > 0 {
		return fmt.Sprintf("%d months", months)
	}

	days := int(now.Sub(createdAt).Hours() / 24)
	return fmt.Sprintf("%d days", days)
}
