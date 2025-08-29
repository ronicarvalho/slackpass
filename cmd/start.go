package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/slackpass/slackpass/internal/vm"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start [name...]",
	Short: "Start virtual machine instances",
	Long: `Start one or more virtual machine instances.

Examples:
  slackpass start myvm
  slackpass start vm1 vm2 vm3`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		manager := vm.NewManager()
		for _, name := range args {
			if err := manager.Start(name); err != nil {
				return fmt.Errorf("failed to start %s: %w", name, err)
			}
			fmt.Printf("Started: %s\n", name)
		}

		return nil
	},
}

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop [name...]",
	Short: "Stop virtual machine instances",
	Long: `Stop one or more virtual machine instances.

Examples:
  slackpass stop myvm
  slackpass stop vm1 vm2 vm3
  slackpass stop --force myvm    # Force stop without graceful shutdown`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

		manager := vm.NewManager()
		for _, name := range args {
			if err := manager.Stop(name, force); err != nil {
				return fmt.Errorf("failed to stop %s: %w", name, err)
			}
			fmt.Printf("Stopped: %s\n", name)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)

	stopCmd.Flags().BoolP("force", "f", false, "Force stop without graceful shutdown")
}