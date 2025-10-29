package presentation

import (
	"fmt"
	"strings"

	"GitInsights/domain"
)

// MarkdownGenerator generates markdown content for profile stats
type MarkdownGenerator struct{}

// NewMarkdownGenerator creates a new markdown generator
func NewMarkdownGenerator() *MarkdownGenerator {
	return &MarkdownGenerator{}
}

// Generate creates markdown content from profile stats
func (m *MarkdownGenerator) Generate(stats *domain.ProfileStats) string {
	var lines []string

	lines = append(lines, "<!--START_SECTION:GitInsights-->")
	lines = append(lines, "### Git Insight")
	lines = append(lines, "")
	lines = append(lines, "👤 **Account Age:** "+stats.AccountAge)
	lines = append(lines, "")
	lines = append(lines, "🔥 **Current Streak:** "+fmt.Sprintf("%d days", stats.CurrentStreak))
	lines = append(lines, "")
	lines = append(lines, "🏆 **Longest Streak:** "+fmt.Sprintf("%d days", stats.LongestStreak))
	lines = append(lines, "")
	lines = append(lines, "📅 **Most Productive Day:** "+stats.MostProductiveDay)
	lines = append(lines, "")
	lines = append(lines, "⌚️ **Most Productive Hour:** "+stats.MostProductiveHour)
	lines = append(lines, "")
	lines = append(lines, "📊 **Weekly Commit Distribution:**")
	lines = append(lines, "```")

	// Order days from Monday to Sunday
	dayOrder := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	maxCommits := 0
	for _, count := range stats.WeeklyDistribution {
		if count > maxCommits {
			maxCommits = count
		}
	}

	for _, day := range dayOrder {
		count := stats.WeeklyDistribution[day]
		bar := m.generateCommitBar(count, maxCommits)
		lines = append(lines, fmt.Sprintf("%-9s [%s] %d commits", day, bar, count))
	}

	lines = append(lines, "```")
	lines = append(lines, "")
	lines = append(lines, "💻 **Language Statistics:**")
	lines = append(lines, "```")

	// Find max language name length for alignment
	maxLength := m.maxLanguageLength(stats.Languages)

	// Generate language statistics
	for _, lang := range stats.Languages {
		progressBar := m.generateProgressBar(lang.Percentage)
		lines = append(lines, fmt.Sprintf(
			"%-*s [%-30s] %5.2f%%",
			maxLength,
			lang.Language,
			progressBar,
			lang.Percentage,
		))
	}

	lines = append(lines, "```")
	lines = append(lines, "")
	lines = append(lines, " _Last update: "+stats.LastUpdated.Format("2006-01-02 15:04:05")+"_")
	lines = append(lines, "<!--END_SECTION:GitInsights-->")

	return strings.Join(lines, "\n")
}

// generateProgressBar creates a visual progress bar
func (m *MarkdownGenerator) generateProgressBar(percentage float64) string {
	const barWidth = 40
	numFilled := int(percentage / 100 * barWidth)
	filled := strings.Repeat("█", numFilled)
	empty := strings.Repeat("░", barWidth-numFilled)
	return filled + empty
}

// generateCommitBar creates a visual bar for commit counts
func (m *MarkdownGenerator) generateCommitBar(count, maxCount int) string {
	const barWidth = 30
	if maxCount == 0 {
		return strings.Repeat("░", barWidth)
	}
	numFilled := int(float64(count) / float64(maxCount) * float64(barWidth))
	filled := strings.Repeat("█", numFilled)
	empty := strings.Repeat("░", barWidth-numFilled)
	return filled + empty
}

// maxLanguageLength finds the longest language name for alignment
func (m *MarkdownGenerator) maxLanguageLength(languages []domain.LanguageStats) int {
	maxLength := 0
	for _, lang := range languages {
		if len(lang.Language) > maxLength {
			maxLength = len(lang.Language)
		}
	}
	return maxLength
}
