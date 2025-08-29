package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/slackpass/slackpass/internal/vm"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info [name...]",
	Short: "Display detailed information about virtual machine instances",
	Long: `Display detailed information about one or more virtual machine instances.

Shows comprehensive details including state, resource allocation,
network configuration, and system information.

Examples:
  slackpass info myvm
  slackpass info vm1 vm2
  slackpass info --all        # Show info for all instances`,
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")

		manager := vm.NewManager()

		if all {
			instances, err := manager.List()
			if err != nil {
				return fmt.Errorf("failed to list instances: %w", err)
			}

			for i, instance := range instances {
				if i > 0 {
					fmt.Println()
				}
				if err := manager.Info(instance.Name); err != nil {
					return fmt.Errorf("failed to get info for %s: %w", instance.Name, err)
				}
			}
			return nil
		}

		if len(args) == 0 {
			return fmt.Errorf("instance name required (or use --all flag)")
		}

		for i, name := range args {
			if i > 0 {
				fmt.Println()
			}
			if err := manager.Info(name); err != nil {
				return fmt.Errorf("failed to get info for %s: %w", name, err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	infoCmd.Flags().Bool("all", false, "Show info for all instances")
}