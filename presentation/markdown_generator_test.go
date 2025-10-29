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

	// Check for content
	if !strings.Contains(markdown, "# 📊 Git Insights") {
		t.Error("Expected title in markdown")
	}

	if !strings.Contains(markdown, "## 💻 Language Statistics") {
		t.Error("Expected language statistics section")
	}

	if !strings.Contains(markdown, "## 📈 Profile Overview") {
		t.Error("Expected profile overview section")
	}

	// Check for badge-like elements
	if !strings.Contains(markdown, "![Most Productive Day](https://img.shields.io/badge/Most_Productive_Day-Monday-green?style=for-the-badge&logo=calendar)") {
		t.Error("Expected most productive day badge in markdown")
	}

	if !strings.Contains(markdown, "![Most Productive Hour](https://img.shields.io/badge/Most_Productive_Hour-10:00_--_11:00-purple?style=for-the-badge&logo=clock)") {
		t.Error("Expected most productive hour badge in markdown")
	}

	if !strings.Contains(markdown, "![Account Age](https://img.shields.io/badge/Account_Age-5_years_9_months-blue?style=for-the-badge&logo=github)") {
		t.Error("Expected account age badge in markdown")
	}

	if !strings.Contains(markdown, "![Current Streak](https://img.shields.io/badge/Current_Streak-15_days-orange?style=for-the-badge&logo=fire)") {
		t.Error("Expected current streak badge in markdown")
	}

	if !strings.Contains(markdown, "![Longest Streak](https://img.shields.io/badge/Longest_Streak-45_days-red?style=for-the-badge&logo=trophy)") {
		t.Error("Expected longest streak badge in markdown")
	}

	if !strings.Contains(markdown, "## 📊 Weekly Commit Distribution") {
		t.Error("Expected weekly distribution section")
	}

	// Check for centered divs
	if !strings.Contains(markdown, "<div align=\"center\">") {
		t.Error("Expected centered div in markdown")
	}

	// Check for medals in language stats (top 3)
	if !strings.Contains(markdown, "🥇") {
		t.Error("Expected gold medal for top language")
	}

	if !strings.Contains(markdown, "🥈") {
		t.Error("Expected silver medal for second language")
	}
}

func TestProgressBarGeneration(t *testing.T) {
	gen := presentation.NewMarkdownGenerator(true)
	stats := &domain.ProfileStats{
		Languages: []domain.LanguageStats{
			{Language: "Go", Bytes: 100, Percentage: 100.0},
		},
		LastUpdated: time.Date(2023, 11, 15, 12, 0, 0, 0, time.UTC),
	}

	markdown := gen.Generate(stats)

	// Should contain filled progress bar for 100% (30 characters wide)
	if !strings.Contains(markdown, "██████████████████████████████") {
		t.Error("Expected full progress bar for 100%")
	}
}

func TestCreditGeneration(t *testing.T) {
	// Test with credit enabled
	gen := presentation.NewMarkdownGenerator(true)
	stats := &domain.ProfileStats{
		Languages: []domain.LanguageStats{
			{Language: "Go", Bytes: 100, Percentage: 100.0},
		},
		LastUpdated: time.Date(2023, 11, 15, 12, 0, 0, 0, time.UTC),
	}

	markdown := gen.Generate(stats)
	if !strings.Contains(markdown, "Generated with [GitInsights]") {
		t.Error("Expected credit line when showCredit is true")
	}

	// Test with credit disabled
	genNoCredit := presentation.NewMarkdownGenerator(false)
	markdownNoCredit := genNoCredit.Generate(stats)
	if strings.Contains(markdownNoCredit, "Generated with [GitInsights]") {
		t.Error("Expected no credit line when showCredit is false")
	}
}

func TestWeeklyDistributionEmojis(t *testing.T) {
	gen := presentation.NewMarkdownGenerator(false)
	stats := &domain.ProfileStats{
		Languages: []domain.LanguageStats{
			{Language: "Go", Bytes: 100, Percentage: 100.0},
		},
		WeeklyDistribution: map[string]int{
			"Monday":   10,
			"Saturday": 5,
			"Sunday":   3,
		},
		LastUpdated: time.Date(2023, 11, 15, 12, 0, 0, 0, time.UTC),
	}

	markdown := gen.Generate(stats)

	// Check for weekday and weekend emojis
	if !strings.Contains(markdown, "📅 Monday") {
		t.Error("Expected weekday emoji for Monday")
	}

	if !strings.Contains(markdown, "🎉 Saturday") {
		t.Error("Expected weekend emoji for Saturday")
	}

	if !strings.Contains(markdown, "🎉 Sunday") {
		t.Error("Expected weekend emoji for Sunday")
	}
}
