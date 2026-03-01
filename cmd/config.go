package cmd

import (
	"fmt"
	"my_skill_manager/internal/config"

	"github.com/spf13/cobra"
)

func init() {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
	}

	setDefaultStyle := &cobra.Command{
		Use:   "set-default-style <style>",
		Short: "Set the default style for links",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			style := args[0]
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			cfg.DefaultStyle = style
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
			fmt.Printf("Default style set to %s\n", style)
			return nil
		},
	}

	getDefaultStyle := &cobra.Command{
		Use:   "get-default-style",
		Short: "Get the current default style",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			fmt.Printf("Default style: %s\n", cfg.DefaultStyle)
			return nil
		},
	}

	configCmd.AddCommand(setDefaultStyle, getDefaultStyle)
	addCmd(configCmd)
}
