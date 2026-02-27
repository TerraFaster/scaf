// Package config manages the global scaf configuration file.
//
// The config file lives at ~/.scaf/config.yaml and is created automatically
// on the first run if it does not exist. It is a single, user-wide file —
// one config for all projects.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const FileName = "config.yaml"

// Config represents the user-wide ~/.scaf/config.yaml configuration.
type Config struct {
	DefaultLicense string   `yaml:"default_license"`
	DefaultIgnore  []string `yaml:"default_ignore"`
	CacheTTLHours  int      `yaml:"cache_ttl_hours"`
	Interactive    bool     `yaml:"interactive"`
	DefaultAuthor  string   `yaml:"default_author"`
	DefaultYear    string   `yaml:"default_year"`
}

// Defaults returns a Config populated with sensible default values.
func Defaults() Config {
	return Config{
		CacheTTLHours: 168, // 7 days
		Interactive:   true,
	}
}

// Dir returns the path to the ~/.scaf directory.
func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ".scaf"), nil
}

// Path returns the absolute path to ~/.scaf/config.yaml.
func Path() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, FileName), nil
}

// Load reads ~/.scaf/config.yaml, creating it with defaults if it does not
// exist yet. It never returns an error just because the file is missing.
func Load() (Config, error) {
	cfg := Defaults()

	path, err := Path()
	if err != nil {
		return cfg, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		// First run — write defaults and return them.
		if writeErr := write(cfg, path); writeErr != nil {
			// Non-fatal: we can still run with defaults even if we cannot write.
			fmt.Fprintf(os.Stderr,
				"⚠  Could not create config file at %s: %v\n", path, writeErr)
		}
		return cfg, nil
	}
	if err != nil {
		return cfg, fmt.Errorf("failed to read %s: %w", path, err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse %s: %w", path, err)
	}

	return cfg, nil
}

// Save writes cfg back to ~/.scaf/config.yaml.
func Save(cfg Config) error {
	path, err := Path()
	if err != nil {
		return err
	}
	return write(cfg, path)
}

// write is the internal helper that serialises cfg to an arbitrary path.
func write(cfg Config, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("cannot create config directory: %w", err)
	}

	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("failed to serialise config: %w", err)
	}

	header := []byte(`# scaf configuration — https://github.com/TerraFaster/scaf
#
# default_license   : license key used when running "scaf license" with no argument
# default_ignore    : templates used when running "scaf ignore" with no argument
# cache_ttl_hours   : how long to keep cached templates before re-downloading (default 168 = 7 days)
# interactive       : always show interactive picker even when a default is configured
# default_author    : author name for LICENSE placeholders (fallback: git config user.name)
# default_year      : year for LICENSE placeholders (fallback: current year)

`)

	if err := os.WriteFile(path, append(header, data...), 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	return nil
}
