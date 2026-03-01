package cmd

import (
	"github.com/spf13/cobra"
	"my_skill_manager/internal/status"
)

func init() {
	statusCmd := &cobra.Command{
		Use:   "status [path]",
		Short: "Show status of symlink or workspace",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var p string
			if len(args) == 1 {
				p = args[0]
			}
			return status.Status(p)
		},
	}
	addCmd(statusCmd)
}
