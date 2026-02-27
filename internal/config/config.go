// Package config manages the global scaf configuration file.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const FileName = "config.yaml"

// GitProfile holds a named git identity.
type GitProfile struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

// Config represents the user-wide ~/.scaf/config.yaml configuration.
type Config struct {
	DefaultLicense string                `yaml:"default_license"`
	DefaultIgnore  []string              `yaml:"default_ignore"`
	CacheTTLHours  int                   `yaml:"cache_ttl_hours"`
	Interactive    bool                  `yaml:"interactive"`
	DefaultAuthor  string                `yaml:"default_author"`
	DefaultYear    string                `yaml:"default_year"`
	Profiles       map[string]GitProfile `yaml:"profiles,omitempty"`
}

// Defaults returns a Config populated with sensible default values.
func Defaults() Config {
	return Config{
		CacheTTLHours: 168,
		Interactive:   true,
		Profiles:      map[string]GitProfile{},
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

// TemplatesDir returns the path to ~/.scaf/templates.
func TemplatesDir() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "templates"), nil
}

// Load reads ~/.scaf/config.yaml, creating it with defaults if it does not exist.
func Load() (Config, error) {
	cfg := Defaults()

	path, err := Path()
	if err != nil {
		return cfg, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		if writeErr := write(cfg, path); writeErr != nil {
			fmt.Fprintf(os.Stderr, "⚠  Could not create config file at %s: %v\n", path, writeErr)
		}
		return cfg, nil
	}
	if err != nil {
		return cfg, fmt.Errorf("failed to read %s: %w", path, err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse %s: %w", path, err)
	}

	if cfg.Profiles == nil {
		cfg.Profiles = map[string]GitProfile{}
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

func write(cfg Config, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("cannot create config directory: %w", err)
	}

	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("failed to serialise config: %w", err)
	}

	header := []byte("# scaf configuration — https://github.com/TerraFaster/scaf\n" +
		"#\n" +
		"# default_license   : license key used when running \"scaf license\" with no argument\n" +
		"# default_ignore    : templates used when running \"scaf ignore\" with no argument\n" +
		"# cache_ttl_hours   : how long to keep cached templates (default 168 = 7 days)\n" +
		"# interactive       : show interactive picker when a default is configured\n" +
		"# default_author    : author name for LICENSE placeholders\n" +
		"# default_year      : year for LICENSE placeholders\n" +
		"# profiles          : named git profiles for 'scaf git profile <name>'\n" +
		"#   work:\n" +
		"#     name: Username\n" +
		"#     email: username@company.com\n\n")

	if err := os.WriteFile(path, append(header, data...), 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	return nil
}
