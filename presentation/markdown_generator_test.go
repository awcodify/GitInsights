package presentation_test

import (
	"strings"
	"testing"
	"time"

	"GitInsights/domain"
	"GitInsights/presentation"
)

func TestMarkdownGeneration(t *testing.T) {
	stats := &domain.ProfileStats{
		Username:   "testuser",
		TotalBytes: 1500,
		Languages: []domain.LanguageStats{
			{Language: "Go", Bytes: 1000, Percentage: 66.67},
			{Language: "Java", Bytes: 500, Percentage: 33.33},
		},
		MostProductiveDay:  "Monday",
		MostProductiveHour: "10:00 - 11:00",
		AccountAge:         "5 years 9 months",
		CurrentStreak:      15,
		LongestStreak:      45,
		WeeklyDistribution: map[string]int{
			"Monday":    10,
			"Tuesday":   8,
			"Wednesday": 12,
			"Thursday":  7,
			"Friday":    9,
			"Saturday":  3,
			"Sunday":    2,
		},
		LastUpdated: time.Date(2023, 11, 15, 12, 0, 0, 0, time.UTC),
	}

	gen := presentation.NewMarkdownGenerator(true)
	markdown := gen.Generate(stats)

	// Check if markdown contains expected sections
	if markdown == "" {
		t.Error("Expected markdown content, got empty string")
	}

	// Check for markers
	if !strings.Contains(markdown, "<!--START_SECTION:GitInsights-->") {
		t.Error("Expected start marker in markdown")
	}

	if !strings.Contains(markdown, "<!--END_SECTION:GitInsights-->") {
		t.Error("Expected end marker in markdown")
	}

	if !strings.Contains(markdown, "Language Distribution") {
		t.Error("Expected language distribution section")
	}

	if !strings.Contains(markdown, "Most Productive Day") {
		t.Error("Expected most productive day in markdown")
	}

	if !strings.Contains(markdown, "Peak Hours") {
		t.Error("Expected peak hours in markdown")
	}

	if !strings.Contains(markdown, "5 years 9 months") {
		t.Error("Expected account age in markdown")
	}

	if !strings.Contains(markdown, "15 days") {
		t.Error("Expected current streak in markdown")
	}

	if !strings.Contains(markdown, "45 days") {
		t.Error("Expected longest streak in markdown")
	}

	if !strings.Contains(markdown, "Weekly Activity") {
		t.Error("Expected weekly activity section")
	}

	// Check for modern elements
	if !strings.Contains(markdown, "<div align=\"center\">") {
		t.Error("Expected centered div for modern layout")
	}

	if !strings.Contains(markdown, "<table") {
		t.Error("Expected table for stats cards")
	}

	if !strings.Contains(markdown, "<details>") {
		t.Error("Expected collapsible details section")
	}
}

func TestProgressBarGeneration(t *testing.T) {
	gen := presentation.NewMarkdownGenerator(true)
	stats := &domain.ProfileStats{
		Languages: []domain.LanguageStats{
			{Language: "Go", Bytes: 100, Percentage: 100.0},
		},
	}

	markdown := gen.Generate(stats)

	// Should contain filled progress bar for 100%
	if !strings.Contains(markdown, "████████████████████████████████████████") {
		t.Error("Expected full progress bar for 100%")
	}
}
