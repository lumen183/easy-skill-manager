package move

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"my_skill_manager/internal/config"
)

// Move moves a directory into a named repo and leaves a symlink at the
// original location. If a .git directory exists it will be moved as well
// (see notes). dryRun will only print planned actions.
func Move(sourcePath, repoName string, dryRun bool) error {
	if sourcePath == "" {
		return errors.New("source path is required")
	}
	// stat source
	sfi, err := os.Stat(sourcePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("source path %s does not exist", sourcePath)
		}
		return fmt.Errorf("failed to stat source: %w", err)
	}
	if !sfi.IsDir() {
		return fmt.Errorf("source %s is not a directory", sourcePath)
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	r, ok := cfg.Repos[repoName]
	if !ok {
		return fmt.Errorf("repo %q not found in config", repoName)
	}

	absSource, err := filepath.Abs(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute source: %w", err)
	}
	repoPath := r.Path
	absRepoPath, err := filepath.Abs(repoPath)
	if err != nil {
		return fmt.Errorf("failed to resolve repo path: %w", err)
	}

	targetPath := filepath.Join(absRepoPath, filepath.Base(absSource))

	// ensure target does not exist
	if _, err := os.Lstat(targetPath); err == nil {
		return fmt.Errorf("target %s already exists in repo", targetPath)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat target: %w", err)
	}

	// detect .git
	gitPath := filepath.Join(absSource, ".git")
	hasGit := false
	if gi, err := os.Stat(gitPath); err == nil && gi.IsDir() {
		hasGit = true
	}

	if dryRun {
		fmt.Printf("Dry-run: move %s -> %s\n", absSource, targetPath)
		if hasGit {
			fmt.Printf("Detected .git at %s; would move contents (and .git) to target\n", gitPath)
		} else {
			fmt.Printf("No .git detected; would move entire directory\n")
		}
		fmt.Printf("Dry-run: would create symlink %s -> %s\n", absSource, targetPath)
		return nil
	}

	// create target dir parent (repoPath exists, but ensure)
	if err := os.MkdirAll(absRepoPath, 0o755); err != nil {
		return fmt.Errorf("failed to ensure repo path: %w", err)
	}

	if hasGit {
		// create target dir
		if err := os.MkdirAll(targetPath, 0o755); err != nil {
			return fmt.Errorf("failed to create target dir: %w", err)
		}
		// move every entry from source into target except .git
		entries, err := os.ReadDir(absSource)
		if err != nil {
			return fmt.Errorf("failed to read source entries: %w", err)
		}
		for _, e := range entries {
			if e.Name() == ".git" {
				// leave .git in place
				continue
			}
			src := filepath.Join(absSource, e.Name())
			dst := filepath.Join(targetPath, e.Name())
			if err := os.Rename(src, dst); err != nil {
				// attempt copy+remove for robustness
				if err := copyPath(src, dst); err != nil {
					return fmt.Errorf("failed to move %s -> %s: %w", src, dst, err)
				}
				// remove original
				if err := os.RemoveAll(src); err != nil {
					return fmt.Errorf("moved but failed to remove original %s: %w", src, err)
				}
			}
			// create symlink at original entry location pointing to the new location
			if err := os.Symlink(dst, src); err != nil {
				return fmt.Errorf("failed to create symlink for %s -> %s: %w", src, dst, err)
			}
		}
		// leave .git in place; do not replace the source dir itself
		fmt.Printf("Moved contents (excluding .git) to %s, left .git in place and created symlinks for moved entries\n", targetPath)
		return nil
	}

	// no .git: move entire directory
	if err := os.Rename(absSource, targetPath); err != nil {
		// fallback to copy+remove
		if err := copyPath(absSource, targetPath); err != nil {
			return fmt.Errorf("failed to move directory: %w", err)
		}
		if err := os.RemoveAll(absSource); err != nil {
			return fmt.Errorf("moved but failed to remove original: %w", err)
		}
	}
	// create symlink at original location pointing to absolute target
	if err := os.Symlink(targetPath, absSource); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}
	fmt.Printf("Moved %s -> %s and created symlink %s -> %s\n", absSource, targetPath, absSource, targetPath)
	return nil
}

// copyPath copies files or directories from src to dst.
func copyPath(src, dst string) error {
	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if si.IsDir() {
		if err := os.MkdirAll(dst, 0o755); err != nil {
			return err
		}
		entries, err := os.ReadDir(src)
		if err != nil {
			return err
		}
		for _, e := range entries {
			s := filepath.Join(src, e.Name())
			d := filepath.Join(dst, e.Name())
			if err := copyPath(s, d); err != nil {
				return err
			}
		}
		return nil
	}
	// file copy
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	// try to copy mode
	if err := os.Chmod(dst, si.Mode()); err != nil {
		// non-fatal
	}
	return nil
}
