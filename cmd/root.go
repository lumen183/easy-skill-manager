package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "skillmgr",
	Short: "skillmgr manages local skill repositories",
	Long:  "skillmgr is a small tool to register and manage local skill repositories.",
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func addCmd(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}
