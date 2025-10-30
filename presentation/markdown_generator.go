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

	// Modern Header with Gradient Effect
	lines = append(lines, "<div align=\"center\">")
	lines = append(lines, "")
	lines = append(lines, "[![Profile Stats](https://img.shields.io/badge/Git-Insights-blueviolet?style=for-the-badge&logo=github)](https://github.com/awcodify/GitInsights)")
	lines = append(lines, "")
	lines = append(lines, "</div>")
	lines = append(lines, "")
	lines = append(lines, "---")
	lines = append(lines, "")

	// Stats Overview Cards
	lines = append(lines, "<div align=\"center\">")
	lines = append(lines, "")
	lines = append(lines, "## 🎯 Quick Stats")
	lines = append(lines, "")
	lines = append(lines, "</div>")
	lines = append(lines, "")

	// Create a visually appealing table for quick stats
	lines = append(lines, "<table align=\"center\">")
	lines = append(lines, "<tr>")
	lines = append(lines, "<td align=\"center\" width=\"200\">")
	lines = append(lines, "<img src=\"https://img.icons8.com/fluency/96/000000/resume.png\" width=\"48\"/>")
	lines = append(lines, "<br><strong>Account Age</strong>")
	lines = append(lines, "<br><code>"+stats.AccountAge+"</code>")
	lines = append(lines, "</td>")
	lines = append(lines, "<td align=\"center\" width=\"200\">")
	lines = append(lines, "<img src=\"https://img.icons8.com/fluency/96/000000/fire-element.png\" width=\"48\"/>")
	lines = append(lines, "<br><strong>Current Streak</strong>")
	lines = append(lines, "<br><code>"+fmt.Sprintf("%d days", stats.CurrentStreak)+"</code>")
	lines = append(lines, "</td>")
	lines = append(lines, "<td align=\"center\" width=\"200\">")
	lines = append(lines, "<img src=\"https://img.icons8.com/fluency/96/000000/trophy.png\" width=\"48\"/>")
	lines = append(lines, "<br><strong>Longest Streak</strong>")
	lines = append(lines, "<br><code>"+fmt.Sprintf("%d days", stats.LongestStreak)+"</code>")
	lines = append(lines, "</td>")
	lines = append(lines, "</tr>")
	lines = append(lines, "</table>")
	lines = append(lines, "")

	// Productivity Insights
	lines = append(lines, "<div align=\"center\">")
	lines = append(lines, "")
	lines = append(lines, "## ⚡ Productivity Insights")
	lines = append(lines, "")
	lines = append(lines, "</div>")
	lines = append(lines, "")

	lines = append(lines, "<table align=\"center\">")
	lines = append(lines, "<tr>")
	lines = append(lines, "<td align=\"center\">")
	lines = append(lines, "<img src=\"https://img.icons8.com/fluency/96/000000/calendar.png\" width=\"40\"/>")
	lines = append(lines, "<br><strong>Most Productive Day</strong>")
	lines = append(lines, "<br>"+m.getDayEmoji(stats.MostProductiveDay)+" <code>"+stats.MostProductiveDay+"</code>")
	lines = append(lines, "</td>")
	lines = append(lines, "<td align=\"center\">")
	lines = append(lines, "<img src=\"https://img.icons8.com/fluency/96/000000/clock.png\" width=\"40\"/>")
	lines = append(lines, "<br><strong>Peak Hours</strong>")
	lines = append(lines, "<br>⏰ <code>"+stats.MostProductiveHour+"</code>")
	lines = append(lines, "</td>")
	lines = append(lines, "</tr>")
	lines = append(lines, "</table>")
	lines = append(lines, "")

	// Weekly Activity Chart
	lines = append(lines, "<div align=\"center\">")
	lines = append(lines, "")
	lines = append(lines, "## 📈 Weekly Activity")
	lines = append(lines, "")
	lines = append(lines, "</div>")
	lines = append(lines, "")
	lines = append(lines, "```text")

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
		bar := m.generateModernCommitBar(count, maxCommits)
		emoji := m.getDayEmoji(day)
		lines = append(lines, fmt.Sprintf("%s %-9s %s %4d commits", emoji, day, bar, count))
	}

	lines = append(lines, "```")
	lines = append(lines, "")

	// Language Distribution
	lines = append(lines, "<div align=\"center\">")
	lines = append(lines, "")
	lines = append(lines, "## 💻 Language Distribution")
	lines = append(lines, "")
	lines = append(lines, "</div>")
	lines = append(lines, "")

	// Create visual language cards
	lines = append(lines, "<div align=\"center\">")
	lines = append(lines, "")

	// Generate language statistics with modern styling
	for i, lang := range stats.Languages {
		if i >= 5 { // Show top 5 in visual cards
			break
		}

		lines = append(lines, fmt.Sprintf("![%s](https://img.shields.io/badge/%s-%.1f%%25-blue?style=flat-square&logo=%s)",
			lang.Language,
			strings.ReplaceAll(lang.Language, " ", "_"),
			lang.Percentage,
			m.getLanguageLogo(lang.Language)))

		if (i+1)%3 == 0 {
			lines = append(lines, "")
		} else {
			lines[len(lines)-1] += " "
		}
	}

	lines = append(lines, "")
	lines = append(lines, "</div>")
	lines = append(lines, "")

	// Detailed breakdown
	lines = append(lines, "<details>")
	lines = append(lines, "<summary><b>📊 Detailed Breakdown</b></summary>")
	lines = append(lines, "")
	lines = append(lines, "```text")

	// Find max language name length for alignment
	maxLength := m.maxLanguageLength(stats.Languages)

	// Generate all language statistics
	for _, lang := range stats.Languages {
		progressBar := m.generateColoredProgressBar(lang.Percentage)
		langEmoji := m.getLanguageEmoji(lang.Language)
		lines = append(lines, fmt.Sprintf(
			"%s %-*s %s %6.2f%%",
			langEmoji,
			maxLength,
			lang.Language,
			progressBar,
			lang.Percentage,
		))
	}

	lines = append(lines, "```")
	lines = append(lines, "")
	lines = append(lines, "</details>")
	lines = append(lines, "")

	// Footer
	lines = append(lines, "---")
	lines = append(lines, "")
	lines = append(lines, "<div align=\"center\">")
	lines = append(lines, "")

	// Add last update and optional credit
	lines = append(lines, "<sub>📅 Last updated: "+stats.LastUpdated.Format("Monday, January 2, 2006 at 3:04 PM")+"</sub>")

	if m.showCredit {
		lines = append(lines, "")
		lines = append(lines, "<sub>⚡ Generated with [GitInsights](https://github.com/awcodify/GitInsights)</sub>")
	}

	lines = append(lines, "")
	lines = append(lines, "</div>")
	lines = append(lines, "")
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

// generateColoredProgressBar creates a colorful progress bar
func (m *MarkdownGenerator) generateColoredProgressBar(percentage float64) string {
	const barWidth = 40
	numFilled := int(percentage / 100 * barWidth)

	// Use different characters for visual interest
	var filled string
	if percentage >= 50 {
		filled = strings.Repeat("█", numFilled)
	} else if percentage >= 20 {
		filled = strings.Repeat("▓", numFilled)
	} else if percentage >= 5 {
		filled = strings.Repeat("▒", numFilled)
	} else {
		filled = strings.Repeat("░", numFilled)
	}

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

// generateModernCommitBar creates a modern visual bar for commit counts
func (m *MarkdownGenerator) generateModernCommitBar(count, maxCount int) string {
	const barWidth = 30
	if maxCount == 0 {
		return strings.Repeat("░", barWidth)
	}

	numFilled := int(float64(count) / float64(maxCount) * float64(barWidth))

	// Use gradient-like characters
	var filled string
	percentage := float64(count) / float64(maxCount) * 100

	if percentage >= 75 {
		filled = strings.Repeat("█", numFilled)
	} else if percentage >= 50 {
		filled = strings.Repeat("▓", numFilled)
	} else if percentage >= 25 {
		filled = strings.Repeat("▒", numFilled)
	} else {
		filled = strings.Repeat("░", numFilled)
	}

	empty := strings.Repeat("░", barWidth-numFilled)
	return filled + empty
}

// getDayEmoji returns an emoji for each day of the week
func (m *MarkdownGenerator) getDayEmoji(day string) string {
	emojiMap := map[string]string{
		"Monday":    "🌙",
		"Tuesday":   "🔥",
		"Wednesday": "💎",
		"Thursday":  "⚡",
		"Friday":    "🎉",
		"Saturday":  "🌟",
		"Sunday":    "☀️",
	}

	if emoji, ok := emojiMap[day]; ok {
		return emoji
	}
	return "📅"
}

// getLanguageEmoji returns an emoji for programming languages
func (m *MarkdownGenerator) getLanguageEmoji(language string) string {
	emojiMap := map[string]string{
		"Go":          "🔵",
		"JavaScript":  "🟨",
		"TypeScript":  "🔷",
		"Python":      "🐍",
		"Java":        "☕",
		"Ruby":        "💎",
		"PHP":         "🐘",
		"C":           "⚙️",
		"C++":         "⚡",
		"C#":          "💜",
		"Rust":        "🦀",
		"Swift":       "🍎",
		"Kotlin":      "🟣",
		"Scala":       "🔴",
		"Elixir":      "💧",
		"HTML":        "🌐",
		"CSS":         "🎨",
		"SCSS":        "🎨",
		"Shell":       "🐚",
		"Vim Script":  "🟢",
		"Lua":         "🌙",
		"Dart":        "🎯",
		"R":           "📊",
		"Julia":       "🟣",
		"Haskell":     "🔮",
		"Perl":        "🐪",
		"Objective-C": "🍎",
		"Matlab":      "📈",
	}

	if emoji, ok := emojiMap[language]; ok {
		return emoji
	}
	return "💻"
}

// getLanguageLogo returns a logo identifier for shields.io badges
func (m *MarkdownGenerator) getLanguageLogo(language string) string {
	logoMap := map[string]string{
		"Go":          "go",
		"JavaScript":  "javascript",
		"TypeScript":  "typescript",
		"Python":      "python",
		"Java":        "java",
		"Ruby":        "ruby",
		"PHP":         "php",
		"C":           "c",
		"C++":         "cplusplus",
		"C#":          "csharp",
		"Rust":        "rust",
		"Swift":       "swift",
		"Kotlin":      "kotlin",
		"Scala":       "scala",
		"Elixir":      "elixir",
		"HTML":        "html5",
		"CSS":         "css3",
		"SCSS":        "sass",
		"Shell":       "gnubash",
		"Vim Script":  "vim",
		"Lua":         "lua",
		"Dart":        "dart",
		"R":           "r",
		"Julia":       "julia",
		"Haskell":     "haskell",
		"Perl":        "perl",
		"Objective-C": "apple",
		"Matlab":      "mathworks",
	}

	if logo, ok := logoMap[language]; ok {
		return logo
	}
	return "code"
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
