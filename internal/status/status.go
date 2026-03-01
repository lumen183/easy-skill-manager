package status

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"my_skill_manager/internal/config"
)

// Status reports on a single path when provided, otherwise inspects the
// current working directory for symlinks and prints a summary for each.
func Status(path string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if path != "" {
		return singlePathStatus(path)
	}
	return cwdStatus(cfg)
}

func singlePathStatus(p string) error {
	fi, err := os.Lstat(p)
	if err != nil {
		return err
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		// print exact string required by task
		fmt.Println("不是符号链接")
		return nil
	}
	target, err := os.Readlink(p)
	if err != nil {
		return err
	}
	// normalize to absolute
	if !filepath.IsAbs(target) {
		target = filepath.Join(filepath.Dir(p), target)
	}
	abs, err := filepath.Abs(target)
	if err == nil {
		target = abs
	}
	fmt.Println(target)
	return nil
}

func cwdStatus(cfg *config.Config) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(wd)
	if err != nil {
		return err
	}
	for _, e := range entries {
		name := e.Name()
		// use Lstat to detect symlink
		fi, err := os.Lstat(name)
		if err != nil {
			// skip unreadable entries
			continue
		}
		if fi.Mode()&os.ModeSymlink == 0 {
			continue
		}
		linkTarget, err := os.Readlink(name)
		if err != nil {
			fmt.Printf("%s -> (readlink error: %v)\n", name, err)
			continue
		}
		if !filepath.IsAbs(linkTarget) {
			linkTarget = filepath.Join(wd, linkTarget)
		}
		absTarget, err := filepath.Abs(linkTarget)
		if err == nil {
			linkTarget = absTarget
		}
		// check existence
		_, statErr := os.Stat(linkTarget)
		state := "Valid"
		if statErr != nil {
			state = "Broken"
		}

		// find repo that contains this target
		matched := false
		for repoName, r := range cfg.Repos {
			rel, relErr := filepath.Rel(r.Path, linkTarget)
			if relErr != nil {
				continue
			}
			// if rel begins with .. then it's outside repo
			if strings.HasPrefix(rel, "..") || rel == ".." {
				continue
			}
			// matched
			skill := rel
			if skill == "." {
				skill = "."
			}
			fmt.Printf("%s -> %s:%s (%s)\n", name, repoName, skill, state)
			matched = true
			break
		}
		if !matched {
			fmt.Printf("%s -> %s (%s)\n", name, linkTarget, state)
		}
	}
	return nil
}
