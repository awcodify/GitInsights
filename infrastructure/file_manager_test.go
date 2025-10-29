package infrastructure

import (
	"bytes"
	"os"
	"testing"
)

func TestFileManagerUpdateReadme(t *testing.T) {
	// Create a temporary test file
	content := `# Test README

Some content here

<!--START_SECTION:GitInsights-->
Old content that should be replaced
<!--END_SECTION:GitInsights-->

More content after
`

	tmpFile, err := os.CreateTemp("", "test_readme_*.md")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test update
	fm := NewFileManager(tmpFile.Name())
	newContent := `<!--START_SECTION:GitInsights-->
New content from test
<!--END_SECTION:GitInsights-->`

	err = fm.UpdateReadme(newContent)
	if err != nil {
		t.Errorf("UpdateReadme failed: %v", err)
	}

	// Verify the update
	updatedData, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read updated file: %v", err)
	}

	if !bytes.Contains(updatedData, []byte("New content from test")) {
		t.Error("File was not updated with new content")
	}

	if bytes.Contains(updatedData, []byte("Old content that should be replaced")) {
		t.Error("Old content still exists, should have been replaced")
	}

	if !bytes.Contains(updatedData, []byte("Some content here")) {
		t.Error("Content before GitInsights section was removed")
	}

	if !bytes.Contains(updatedData, []byte("More content after")) {
		t.Error("Content after GitInsights section was removed")
	}
}

func TestFileManagerMissingMarkers(t *testing.T) {
	// Create a file without markers
	content := `# Test README
No markers here
`

	tmpFile, err := os.CreateTemp("", "test_readme_*.md")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test should fail with missing markers
	fm := NewFileManager(tmpFile.Name())
	err = fm.UpdateReadme("test content")

	if err == nil {
		t.Error("Expected error for missing markers, got nil")
	}
}
