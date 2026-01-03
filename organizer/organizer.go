package organizer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config maps folder names to a list of file extensions
type Config map[string][]string

// Organize scans the source directory and moves files based on the config
func Organize(sourceDir string, config Config) error {
	// Create a reverse map for faster lookup: extension -> folder
	extToFolder := make(map[string]string)
	for folder, exts := range config {
		for _, ext := range exts {
			// Normalize extensions to lowercase and remove leading dots if present
			cleanExt := strings.ToLower(strings.TrimPrefix(ext, "."))
			extToFolder[cleanExt] = folder
		}
	}

	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue // Skip directories
		}

		info, err := entry.Info()
		if err != nil {
			fmt.Printf("Skipping file %s: %v\n", entry.Name(), err)
			continue
		}

		// Skip hidden files
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(entry.Name()), "."))
		targetFolder, exists := extToFolder[ext]

		if !exists {
			// Optional: Move to "Others" or just skip
			// For now, let's skip files that don't match any rule
			continue
		}

		// Construct target path
		targetPath := filepath.Join(sourceDir, targetFolder)
		if err := os.MkdirAll(targetPath, 0755); err != nil {
			fmt.Printf("Failed to create directory %s: %v\n", targetPath, err)
			continue
		}

		oldPath := filepath.Join(sourceDir, entry.Name())
		newPath := filepath.Join(targetPath, entry.Name())

		// Handle duplicate filenames
		if _, err := os.Stat(newPath); err == nil {
			// File exists, rename the new file
			baseName := strings.TrimSuffix(entry.Name(), "."+ext)
			newPath = filepath.Join(targetPath, fmt.Sprintf("%s_%d.%s", baseName, info.ModTime().Unix(), ext))
		}

		fmt.Printf("Moving %s to %s\n", entry.Name(), targetFolder)
		if err := os.Rename(oldPath, newPath); err != nil {
			fmt.Printf("Failed to move file %s: %v\n", entry.Name(), err)
		}
	}

	return nil
}
