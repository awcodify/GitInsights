package infrastructure

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// FileManager implements domain.FileRepository
type FileManager struct {
	filePath string
}

// NewFileManager creates a new file manager
func NewFileManager(filePath string) *FileManager {
	return &FileManager{
		filePath: filePath,
	}
}

// UpdateReadme updates the README.md file with new content
func (f *FileManager) UpdateReadme(content string) error {
	file, err := os.OpenFile(f.filePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close file: %w", cerr)
		}
	}()

	// Read existing content
	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Find and replace the GitInsights section
	startMarker := []byte("<!--START_SECTION:GitInsights-->")
	endMarker := []byte("<!--END_SECTION:GitInsights-->")

	startIdx := bytes.Index(data, startMarker)
	endIdx := bytes.Index(data, endMarker)

	if startIdx == -1 || endIdx == -1 {
		return fmt.Errorf("GitInsights section markers not found in file")
	}

	// Build new content
	updatedContent := append(data[:startIdx], []byte(content)...)
	updatedContent = append(updatedContent, data[endIdx+len(endMarker):]...)

	// Write updated content
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}

	if _, err := file.WriteAt(updatedContent, 0); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
