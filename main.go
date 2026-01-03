package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"smart-organizer/organizer"
)

func main() {
	// Parse command line flags
	dirPath := flag.String("path", ".", "Path to the directory to organize")
	configPath := flag.String("config", "config.json", "Path to the configuration file")
	flag.Parse()

	// Validate directory path
	absPath, err := filepath.Abs(*dirPath)
	if err != nil {
		fmt.Printf("Error resolving path: %v\n", err)
		os.Exit(1)
	}

	info, err := os.Stat(absPath)
	if err != nil || !info.IsDir() {
		fmt.Printf("Invalid directory path: %s\n", absPath)
		os.Exit(1)
	}

	// Load configuration
	configData, err := os.ReadFile(*configPath)
	if err != nil {
		// If config file not found, try looking in the executable's directory
		exePath, _ := os.Executable()
		exeDir := filepath.Dir(exePath)
		configData, err = os.ReadFile(filepath.Join(exeDir, *configPath))
		if err != nil {
			fmt.Printf("Error reading config file: %v\n", err)
			os.Exit(1)
		}
	}

	var config organizer.Config
	if err := json.Unmarshal(configData, &config); err != nil {
		fmt.Printf("Error parsing config file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Organizing files in: %s\n", absPath)
	
	// Run organizer
	if err := organizer.Organize(absPath, config); err != nil {
		fmt.Printf("Error organizing files: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Organization complete!")
}
