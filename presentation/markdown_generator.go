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
	lines = append(lines, "\nLanguage Statistics:")
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
	lines = append(lines, "\n📅 Most Productive Day: "+stats.MostProductiveDay)
	lines = append(lines, "\n⌚️ Most Productive Hour: "+stats.MostProductiveHour)
	lines = append(lines, "\n _Last update: "+stats.LastUpdated.Format("2006-01-02 15:04:05")+"_")
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
