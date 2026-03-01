package cmd

import (
	"fmt"
	"my_skill_manager/internal/config"
	"my_skill_manager/internal/link"

	"github.com/spf13/cobra"
)

func init() {
	var target string
	var dryRun bool
	var style string

	linkCmd := &cobra.Command{
		Use:   "link <repo> <skill-name>",
		Short: "Create a symlink for a skill from a repo",
		Args:  cobra.ExactArgs(2),
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
			if err := link.Link(repo, skillName, target, style, dryRun); err != nil {
				return err
			}
			return nil
		},
	}
	linkCmd.Flags().StringVar(&target, "target", "", "Target directory to place the symlink (default: cwd)")
	linkCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be done without making changes")
	linkCmd.Flags().StringVar(&style, "style", "", "Style for the link path (default: from config)")

	addCmd(linkCmd)
}

// compile-time check to avoid unused import when building separately
var _ = fmt.Printf
