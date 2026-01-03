package scanner

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
)

type Analytics struct {
	TotalFiles int
	TotalSize  int64
	Extensions map[string]int
	Sizes      map[string]int64
}

// Scan analyzes the directory and returns analytics
func Scan(path string) (*Analytics, error) {
	analytics := &Analytics{
		Extensions: make(map[string]int),
		Sizes:      make(map[string]int64),
	}

	fmt.Println("Calculating total files...")
	
	// Pass 1: Count total files for progress bar
	totalFiles := 0
	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		if !d.IsDir() {
			totalFiles++
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("counting files failed: %w", err)
	}
	analytics.TotalFiles = totalFiles

	// Pass 2: Analyze types and sizes
	bar := progressbar.Default(int64(totalFiles), "Scanning files")
	
	err = filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		
		if d.IsDir() {
			return nil
		}

		bar.Add(1)
		
		info, err := d.Info()
		if err != nil {
			return nil
		}
		
		size := info.Size()
		analytics.TotalSize += size

		ext := strings.ToLower(filepath.Ext(d.Name()))
		if ext == "" {
			ext = "(no extension)"
		} else {
			ext = strings.TrimPrefix(ext, ".")
		}
		
		analytics.Extensions[ext]++
		analytics.Sizes[ext] += size
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("analysis failed: %w", err)
	}

	return analytics, nil
}
