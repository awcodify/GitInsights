package presentation_test

import (
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
		LastUpdated:        time.Date(2023, 11, 15, 12, 0, 0, 0, time.UTC),
	}

	gen := presentation.NewMarkdownGenerator()
	markdown := gen.Generate(stats)

	// Check if markdown contains expected sections
	if markdown == "" {
		t.Error("Expected markdown content, got empty string")
	}

	// Check for markers
	if !contains(markdown, "<!--START_SECTION:GitInsights-->") {
		t.Error("Expected start marker in markdown")
	}

	if !contains(markdown, "<!--END_SECTION:GitInsights-->") {
		t.Error("Expected end marker in markdown")
	}

	// Check for content
	if !contains(markdown, "Git Insight") {
		t.Error("Expected title in markdown")
	}

	if !contains(markdown, "Language Statistics:") {
		t.Error("Expected language statistics section")
	}

	if !contains(markdown, "Most Productive Day: Monday") {
		t.Error("Expected most productive day in markdown")
	}

	if !contains(markdown, "Most Productive Hour: 10:00 - 11:00") {
		t.Error("Expected most productive hour in markdown")
	}
}

func TestProgressBarGeneration(t *testing.T) {
	gen := presentation.NewMarkdownGenerator()
	stats := &domain.ProfileStats{
		Languages: []domain.LanguageStats{
			{Language: "Go", Bytes: 100, Percentage: 100.0},
		},
	}

	markdown := gen.Generate(stats)

	// Should contain filled progress bar for 100%
	if !contains(markdown, "████████████████████████████████████████") {
		t.Error("Expected full progress bar for 100%")
	}
}

// Helper function
func contains(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
