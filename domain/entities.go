package domain

import "time"

// LanguageStats represents statistics about programming language usage
type LanguageStats struct {
	Language   string
	Bytes      int
	Percentage float64
}

// ProfileStats contains aggregated statistics about a GitHub profile
type ProfileStats struct {
	Username           string
	Languages          []LanguageStats
	TotalBytes         int
	MostProductiveDay  string
	MostProductiveHour string
	AccountAge         string
	CurrentStreak      int
	LongestStreak      int
	WeeklyDistribution map[string]int
	LastUpdated        time.Time
}

// Commit represents a simplified commit structure
type Commit struct {
	Date time.Time
}
