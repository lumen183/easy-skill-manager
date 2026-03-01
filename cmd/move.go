package cmd

import (
	"github.com/spf13/cobra"
	"my_skill_manager/internal/move"
)

func init() {
	var dryRun bool
	moveCmd := &cobra.Command{
		Use:   "move <source-path> <repo-name>",
		Short: "Move a skill directory into a repo and leave a symlink",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			src := args[0]
			repo := args[1]
			return move.Move(src, repo, dryRun)
		},
	}
	moveCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show actions without making changes")
	addCmd(moveCmd)
}
