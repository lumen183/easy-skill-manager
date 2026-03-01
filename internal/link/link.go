package link

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"my_skill_manager/internal/config"
)

// Link creates a symlink for a skill from a named repo into targetDir (or cwd when empty).
// repoName: name of repository as stored in config
// skillName: directory or file name inside the repository
// targetDir: where to create the symlink; when empty use current working directory
// style: style for the link path, default "opencode"
// dryRun: when true, only print what would be done
// Link creates a symlink for a skill from a named repo into targetDir (or cwd when empty).
// It supports dryRun mode where no filesystem changes or config saves are performed.
func Link(repoName, skillName, targetDir, style string, dryRun bool) error {
	if repoName == "" {
		return errors.New("repo name is required")
	}
	if skillName == "" {
		return errors.New("skill name is required")
	}
	if style == "" {
		style = "opencode"
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	r, ok := cfg.Repos[repoName]
	if !ok {
		return fmt.Errorf("repo %q not found in config", repoName)
	}

	// source inside the repo
	source := filepath.Join(r.Path, skillName)
	sfi, err := os.Stat(source)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("skill %q not found in repo %q", skillName, repoName)
		}
		return fmt.Errorf("failed to stat source %s: %w", source, err)
	}
	// ensure target dir
	if targetDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working dir: %w", err)
		}
		targetDir = filepath.Join(wd, "."+style, "skills")
	}
	// ensure targetDir exists and is directory
	tfi, err := os.Stat(targetDir)
	if err != nil {
		if os.IsNotExist(err) {
			if !dryRun {
				if err := os.MkdirAll(targetDir, 0o755); err != nil {
					return fmt.Errorf("failed to create target directory %s: %w", targetDir, err)
				}
			}
		} else {
			return fmt.Errorf("failed to stat target dir %s: %w", targetDir, err)
		}
	} else if !tfi.IsDir() {
		return fmt.Errorf("target %s is not a directory", targetDir)
	}

	// use absolute paths for both ends
	absSource, err := filepath.Abs(source)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute source: %w", err)
	}
	targetPath := filepath.Join(targetDir, filepath.Base(skillName))
	absTarget, err := filepath.Abs(targetPath)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute target: %w", err)
	}

	// check existence
	if _, err := os.Lstat(absTarget); err == nil {
		return fmt.Errorf("target %s already exists", absTarget)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat target %s: %w", absTarget, err)
	}

	if dryRun {
		fmt.Printf("Dry-run: will create symlink from %s to %s\n", absSource, absTarget)
		// also print if source is file or dir for clarity
		if sfi.IsDir() {
			fmt.Printf("Source is directory\n")
		} else {
			fmt.Printf("Source is file\n")
		}
		return nil
	}

	// perform symlink creation
	if err := os.Symlink(absSource, absTarget); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}
	fmt.Printf("Created symlink %s -> %s\n", absTarget, absSource)
	return nil
}
