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
	WeeklyProgress     WeeklyProgress
	OvertimeActivity   OvertimeActivity
	LastUpdated        time.Time
}

// WeeklyProgress tracks commits made in the current week vs last week
type WeeklyProgress struct {
	CurrentWeek   int
	LastWeek      int
	ChangeAmount  int
	ChangePercent float64
}

// OvertimeActivity tracks monthly commit trends over the last 6 months
type OvertimeActivity struct {
	MonthlyData []MonthlyCommits
}

// MonthlyCommits represents commit count for a specific month
type MonthlyCommits struct {
	Month   string // Format: "Jan 2024"
	Commits int
}

// Commit represents a simplified commit structure
type Commit struct {
	Date time.Time
}
