package repo

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"my_skill_manager/internal/config"
)

// Add adds a repository with name and path to the config.
// It validates the path exists and is a directory, converts it to an
// absolute path and persists the change. If the name already exists
// an error is returned.
// Add adds a repository entry. If dryRun is true no changes are persisted
// and the function only prints the actions that would be taken.
func Add(name, p, style string, dryRun bool) error {
	if name == "" {
		return errors.New("repo name is required")
	}
	if p == "" {
		return errors.New("path is required")
	}
	if style == "" {
		style = "opencode"
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}
	fi, err := os.Stat(abs)
	if err != nil {
		return fmt.Errorf("path does not exist: %w", err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("path is not a directory: %s", abs)
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if _, ok := cfg.Repos[name]; ok {
		return fmt.Errorf("repo %q already exists", name)
	}
	if dryRun {
		fmt.Printf("Dry-run: would add repo %q -> %s with style %s\n", name, abs, style)
		return nil
	}
	cfg.Repos[name] = config.Repo{Path: abs, CreatedAt: time.Now().Format(time.RFC3339), Style: style}
	return config.Save(cfg)
}

// List returns a sorted list of repository names and the map of repos.
func List() ([]string, map[string]config.Repo, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, nil, err
	}
	names := make([]string, 0, len(cfg.Repos))
	for k := range cfg.Repos {
		names = append(names, k)
	}
	sort.Strings(names)
	return names, cfg.Repos, nil
}

// Remove deletes a repository entry by name and persists the change.
// Remove deletes a repository entry. When dryRun is true it only prints the
// action and does not modify the config file.
func Remove(name string, dryRun bool) error {
	if name == "" {
		return errors.New("repo name is required")
	}
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if _, ok := cfg.Repos[name]; !ok {
		return fmt.Errorf("repo %q not found", name)
	}
	if dryRun {
		fmt.Printf("Dry-run: would remove repo %q\n", name)
		return nil
	}
	delete(cfg.Repos, name)
	return config.Save(cfg)
}
