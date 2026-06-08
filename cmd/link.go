package cmd

import (
	"fmt"
	"strings"

	"my_skill_manager/internal/config"
	"my_skill_manager/internal/link"
	"my_skill_manager/internal/repo"

	"github.com/spf13/cobra"
)

func init() {
	var target string
	var dryRun bool
	var style string
	var copyFlag bool

	linkCmd := &cobra.Command{
		Use:   "link <repo> <skill-name>",
		Short: "Create a symlink for a skill from a repo",
		Args:  cobra.ExactArgs(2),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			// Arg 0: repo name
			if len(args) == 0 {
				names, _, err := repo.List()
				if err != nil {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return filterPrefix(names, toComplete), cobra.ShellCompDirectiveNoFileComp
			}
			// Arg 1: skill name in repo
			if len(args) == 1 {
				repoPath, err := repo.ResolveRepo(args[0])
				if err != nil {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				skills, err := repo.ListSkillsInRepo(repoPath)
				if err != nil {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return filterPrefix(skills, toComplete), cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			repo := args[0]
			skillName := args[1]
			if style == "" {
				cfg, err := config.Load()
				if err != nil {
					return fmt.Errorf("failed to load config: %w", err)
				}
				style = cfg.DefaultStyle
			}
			if err := link.Link(repo, skillName, target, style, dryRun, copyFlag); err != nil {
				return err
			}
			return nil
		},
	}
	linkCmd.Flags().StringVar(&target, "target", "", "Target directory to place the symlink (default: cwd)")
	linkCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be done without making changes")
	linkCmd.Flags().StringVar(&style, "style", "", "Style for the link path (default: from config)")
	linkCmd.Flags().BoolVar(&copyFlag, "copy", false, "Copy the skill instead of creating a symlink")

	addCmd(linkCmd)
}

func filterPrefix(items []string, prefix string) []string {
	if prefix == "" {
		return items
	}
	filtered := make([]string, 0, len(items))
	for _, item := range items {
		if strings.HasPrefix(item, prefix) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// compile-time check to avoid unused import when building separately
var _ = fmt.Printf
