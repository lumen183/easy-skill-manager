package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"my_skill_manager/internal/link"
)

func init() {
	var target string
	var dryRun bool

	linkCmd := &cobra.Command{
		Use:   "link <repo> <skill-name>",
		Short: "Create a symlink for a skill from a repo",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo := args[0]
			skillName := args[1]
			if err := link.Link(repo, skillName, target, dryRun); err != nil {
				return err
			}
			return nil
		},
	}
	linkCmd.Flags().StringVar(&target, "target", "", "Target directory to place the symlink (default: cwd)")
	linkCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be done without making changes")

	addCmd(linkCmd)
}

// compile-time check to avoid unused import when building separately
var _ = fmt.Printf
