package presentation

import (
	"fmt"
	"strings"

	"GitInsights/domain"
)

// MarkdownGenerator generates markdown content for profile stats
type MarkdownGenerator struct {
	showCredit bool
}

// NewMarkdownGenerator creates a new markdown generator
func NewMarkdownGenerator(showCredit bool) *MarkdownGenerator {
	return &MarkdownGenerator{
		showCredit: showCredit,
	}
}

// Generate creates markdown content from profile stats
func (m *MarkdownGenerator) Generate(stats *domain.ProfileStats) string {
	var lines []string

	lines = append(lines, "<!--START_SECTION:GitInsights-->")
	lines = append(lines, "")
	lines = append(lines, "<div align=\"center\">")
	lines = append(lines, "")
	lines = append(lines, "# 📊 Git Insights")
	lines = append(lines, "")
	lines = append(lines, "</div>")
	lines = append(lines, "")
	lines = append(lines, "<div align=\"center\">")
	lines = append(lines, "")
	lines = append(lines, "## 📈 Profile Overview")
	lines = append(lines, "")

	// Create badge-like elements
	lines = append(lines, fmt.Sprintf("![Account Age](https://img.shields.io/badge/Account_Age-%s-blue?style=for-the-badge&logo=github)", m.urlEncode(stats.AccountAge)))
	lines = append(lines, fmt.Sprintf("![Current Streak](https://img.shields.io/badge/Current_Streak-%d_days-orange?style=for-the-badge&logo=fire)", stats.CurrentStreak))
	lines = append(lines, fmt.Sprintf("![Longest Streak](https://img.shields.io/badge/Longest_Streak-%d_days-red?style=for-the-badge&logo=trophy)", stats.LongestStreak))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("![Most Productive Day](https://img.shields.io/badge/Most_Productive_Day-%s-green?style=for-the-badge&logo=calendar)", m.urlEncode(stats.MostProductiveDay)))
	lines = append(lines, fmt.Sprintf("![Most Productive Hour](https://img.shields.io/badge/Most_Productive_Hour-%s-purple?style=for-the-badge&logo=clock)", m.urlEncode(stats.MostProductiveHour)))
	lines = append(lines, "")
	lines = append(lines, "</div>")
	lines = append(lines, "")
	lines = append(lines, "## 📊 Weekly Commit Distribution")
	lines = append(lines, "")
	lines = append(lines, "```text")

	// Order days from Monday to Sunday
	dayOrder := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	dayEmojis := map[string]string{
		"Monday":    "📅",
		"Tuesday":   "📅",
		"Wednesday": "📅",
		"Thursday":  "📅",
		"Friday":    "📅",
		"Saturday":  "🎉",
		"Sunday":    "🎉",
	}
	maxCommits := 0
	for _, count := range stats.WeeklyDistribution {
		if count > maxCommits {
			maxCommits = count
		}
	}

	for _, day := range dayOrder {
		count := stats.WeeklyDistribution[day]
		bar := m.generateCommitBar(count, maxCommits)
		emoji := dayEmojis[day]
		lines = append(lines, fmt.Sprintf("%s %-9s [%s] %d commits", emoji, day, bar, count))
	}

	lines = append(lines, "```")
	lines = append(lines, "")
	lines = append(lines, "## 💻 Language Statistics")
	lines = append(lines, "")
	lines = append(lines, "```text")

	// Find max language name length for alignment
	maxLength := m.maxLanguageLength(stats.Languages)

	// Generate language statistics with emojis
	for i, lang := range stats.Languages {
		progressBar := m.generateProgressBar(lang.Percentage)
		medal := ""
		if i == 0 {
			medal = "🥇 "
		} else if i == 1 {
			medal = "🥈 "
		} else if i == 2 {
			medal = "🥉 "
		} else {
			medal = "   "
		}
		lines = append(lines, fmt.Sprintf(
			"%s%-*s [%-30s] %5.2f%%",
			medal,
			maxLength,
			lang.Language,
			progressBar,
			lang.Percentage,
		))
	}

	lines = append(lines, "```")
	lines = append(lines, "")
	lines = append(lines, "---")
	lines = append(lines, "")

	// Add last update and optional credit
	lines = append(lines, "<div align=\"center\">")
	lines = append(lines, "")
	lastUpdateLine := "⏰ _Last updated: " + stats.LastUpdated.Format("2006-01-02 15:04:05") + "_"
	lines = append(lines, lastUpdateLine)

	if m.showCredit {
		lines = append(lines, "")
		lines = append(lines, "**✨ Generated with [GitInsights](https://github.com/awcodify/GitInsights) ✨**")
	}

	lines = append(lines, "")
	lines = append(lines, "</div>")
	lines = append(lines, "")
	lines = append(lines, "<!--END_SECTION:GitInsights-->")

	return strings.Join(lines, "\n")
}

// generateProgressBar creates a visual progress bar
func (m *MarkdownGenerator) generateProgressBar(percentage float64) string {
	const barWidth = 30
	numFilled := int(percentage / 100 * float64(barWidth))
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

// urlEncode encodes a string for use in URLs (simple encoding for spaces and special chars)
func (m *MarkdownGenerator) urlEncode(s string) string {
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "--")
	return s
}
