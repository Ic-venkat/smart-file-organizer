package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"smart-organizer/internal/organizer"
	"smart-organizer/internal/config"
)

var (
	rootCmd = &cobra.Command{
		Use:   "sfo",
		Short: "Smart File Organizer CLI",
		Long:  `A powerful CLI tool to organize your files based on smart rules.`,
	}
)

var (
	dirPath    string
	configPath string
)

func init() {
	// Add commands
	rootCmd.AddCommand(organizeCmd)

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Path to configuration file")
}

var organizeCmd = &cobra.Command{
	Use:   "organize [path]",
	Short: "Organize files in the specified directory",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetPath := "."
		if len(args) > 0 {
			targetPath = args[0]
		}

		// Load Config
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Organizing files in: %s\n", targetPath)
		if err := organizer.Organize(targetPath, cfg); err != nil {
			fmt.Printf("Error organizing files: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Organization complete!")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
