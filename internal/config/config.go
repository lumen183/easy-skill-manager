package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Config represents the top-level configuration file structure.
type Config struct {
	Version      string          `json:"version"`
	DefaultStyle string          `json:"default_style,omitempty"`
	Repos        map[string]Repo `json:"repos"`
}

// Repo describes a repository entry.
type Repo struct {
	Path      string `json:"path"`
	CreatedAt string `json:"created_at"`
}

// getHomeDir is a thin wrapper around os.UserHomeDir so callers can
// get an error if the home directory cannot be determined.
func getHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if home == "" {
		return "", errors.New("user home directory is empty")
	}
	return home, nil
}

// GetConfigPath returns the absolute path to the config file (~/.skillmgr/config.json).
func GetConfigPath() (string, error) {
	home, err := getHomeDir()
	if err != nil {
		return "", err
	}
	cfgDir := filepath.Join(home, ".skillmgr")
	return filepath.Join(cfgDir, "config.json"), nil
}

// ensureConfigDir makes sure the directory for the config file exists.
func ensureConfigDir() (string, error) {
	home, err := getHomeDir()
	if err != nil {
		return "", err
	}
	cfgDir := filepath.Join(home, ".skillmgr")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		return "", err
	}
	return cfgDir, nil
}

// Load loads the configuration from disk. If the config directory or file
// does not exist, it will create the directory and return a default config.
func Load() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// Ensure directory exists
	if _, err := ensureConfigDir(); err != nil {
		return nil, err
	}

	// If file does not exist, return default config
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			def := &Config{Version: "1.0", DefaultStyle: "opencode", Repos: map[string]Repo{}}
			// Save default to disk so subsequent loads read it
			if serr := Save(def); serr != nil {
				return nil, serr
			}
			return def, nil
		}
		return nil, err
	}

	if fi.IsDir() {
		return nil, errors.New("config path is a directory")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Repos == nil {
		cfg.Repos = map[string]Repo{}
	}
	if cfg.Version == "" {
		cfg.Version = "1.0"
	}
	if cfg.DefaultStyle == "" {
		cfg.DefaultStyle = "opencode"
	}
	return &cfg, nil
}

// Save writes the provided configuration into the config file in JSON
// format with indentation. It ensures the config directory exists and sets
// file permissions to 0644.
func Save(cfg *Config) error {
	if cfg == nil {
		return errors.New("nil config")
	}
	cfgDir, err := ensureConfigDir()
	if err != nil {
		return err
	}
	path := filepath.Join(cfgDir, "config.json")

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	// Write temp file and rename for atomicity
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		// best-effort cleanup
		_ = os.Remove(tmp)
		return err
	}
	return nil
}
