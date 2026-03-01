package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"my_skill_manager/internal/repo"

	"github.com/spf13/cobra"
)

func init() {
	// repo command
	repoCmd := &cobra.Command{
		Use:   "repo",
		Short: "Manage repositories",
	}

	add := &cobra.Command{
		Use:   "add <name> <path>",
		Short: "Add a repository",
		Args:  cobra.ExactArgs(2),
	}
	var addDryRun bool
	var addStyle string
	add.RunE = func(cmd *cobra.Command, args []string) error {
		name := args[0]
		path := args[1]
		return repo.Add(name, path, addStyle, addDryRun)
	}
	add.Flags().BoolVar(&addDryRun, "dry-run", false, "Show actions without making changes")
	add.Flags().StringVar(&addStyle, "style", "opencode", "Style for the repository")

	list := &cobra.Command{
		Use:   "list",
		Short: "List repositories",
		RunE: func(cmd *cobra.Command, args []string) error {
			names, repos, err := repo.List()
			if err != nil {
				return err
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tPATH\tSTYLE\tCREATED")
			for _, n := range names {
				r := repos[n]
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", n, r.Path, r.Style, r.CreatedAt)
			}
			return w.Flush()
		},
	}

	remove := &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a repository",
		Args:  cobra.ExactArgs(1),
	}
	var removeDryRun bool
	remove.RunE = func(cmd *cobra.Command, args []string) error {
		return repo.Remove(args[0], removeDryRun)
	}
	remove.Flags().BoolVar(&removeDryRun, "dry-run", false, "Show actions without making changes")

	repoCmd.AddCommand(add, list, remove)
	addCmd(repoCmd)
}
