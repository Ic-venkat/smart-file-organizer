package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config maps folder names to a list of file extensions
// This matches the structure expected by the organizer package
type Config map[string][]string

// LoadConfig reads configuration from file or environment
func LoadConfig(path string) (Config, error) {
	v := viper.New()

	if path != "" {
		v.SetConfigFile(path)
	} else {
		// Defaults
		v.SetConfigName("config")
		v.SetConfigType("json")
		v.AddConfigPath(".")
		
		// Also look in executable directory
		exePath, err := os.Executable()
		if err == nil {
			v.AddConfigPath(filepath.Dir(exePath))
		}
	}

	// Read config
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if we strictly didn't ask for a path, 
			// BUT for now, let's error out or return default/empty if that's safer.
			// The original code errored if config wasn't found. Let's keep that behavior mostly,
			// or better, return a default config if none found?
			// For this refactor, let's treat "no config found" as an error if user expected one,
			// but if we are just looking in defaults, maybe we should warn.
			// Let's stick to returning error for safety as the tool needs rules to work.
			return nil, fmt.Errorf("config file not found")
		}
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return cfg, nil
}
