package repo

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"my_skill_manager/internal/config"
)

// Add adds a repository with name and path to the config.
// It validates the path exists and is a directory, converts it to an
// absolute path and persists the change. If the name already exists
// an error is returned.
// Add adds a repository entry. If dryRun is true no changes are persisted
// and the function only prints the actions that would be taken.
func Add(name, p string, dryRun bool) error {
	if name == "" {
		return errors.New("repo name is required")
	}
	if p == "" {
		return errors.New("path is required")
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
		fmt.Printf("Dry-run: would add repo %q -> %s\n", name, abs)
		return nil
	}
	cfg.Repos[name] = config.Repo{Path: abs, CreatedAt: time.Now().Format(time.RFC3339)}
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

// ResolveRepo returns the absolute path for a registered repository name.
// It returns an error if the repository is not registered or the path
// does not exist on disk.
func ResolveRepo(name string) (string, error) {
	if name == "" {
		return "", errors.New("repo name is required")
	}
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}
	r, ok := cfg.Repos[name]
	if !ok {
		return "", fmt.Errorf("repo %q not found", name)
	}
	// Ensure the path exists and is a directory
	fi, err := os.Stat(r.Path)
	if err != nil {
		return "", fmt.Errorf("repo path not accessible: %w", err)
	}
	if !fi.IsDir() {
		return "", fmt.Errorf("repo path is not a directory: %s", r.Path)
	}
	abs, err := filepath.Abs(r.Path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve repo path: %w", err)
	}
	return abs, nil
}

// ListSkillsInRepo scans the provided path for immediate child directories
// that contain a file named exactly SKILL.md. It filters out .git,
// node_modules, vendor and hidden directories (starting with '.') and does
// not recurse. Returned list is sorted alphabetically case-insensitive.
func ListSkillsInRepo(path string) ([]string, error) {
	if path == "" {
		return nil, errors.New("path is required")
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}
	entries, err := os.ReadDir(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read repo directory: %w", err)
	}
	skills := make([]string, 0)
	for _, e := range entries {
		name := e.Name()
		// skip non-directories
		if !e.IsDir() {
			continue
		}
		// filter outs
		if name == ".git" || name == "node_modules" || name == "vendor" || strings.HasPrefix(name, ".") {
			continue
		}
		skillFile := filepath.Join(absPath, name, "SKILL.md")
		fi, statErr := os.Stat(skillFile)
		if statErr != nil {
			if os.IsNotExist(statErr) {
				// no SKILL.md, skip
				continue
			}
			// permission or other errors: log warning and skip
			fmt.Fprintf(os.Stderr, "warning: cannot access %s: %v\n", skillFile, statErr)
			continue
		}
		if fi.Mode().IsRegular() {
			skills = append(skills, name)
		}
	}
	sort.Slice(skills, func(i, j int) bool {
		return strings.ToLower(skills[i]) < strings.ToLower(skills[j])
	})
	return skills, nil
}
