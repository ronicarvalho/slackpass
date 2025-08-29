package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/slackpass/slackpass/internal/vm"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [name...]",
	Short: "Delete virtual machine instances",
	Long: `Delete one or more virtual machine instances.

This will stop the instances if they are running and remove all
associated files including disk images.

Examples:
  slackpass delete myvm
  slackpass delete vm1 vm2 vm3
  slackpass delete --purge myvm    # Delete and purge immediately`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		purge, _ := cmd.Flags().GetBool("purge")
		force, _ := cmd.Flags().GetBool("force")

		manager := vm.NewManager()
		for _, name := range args {
			if err := manager.Delete(name, purge, force); err != nil {
				return fmt.Errorf("failed to delete %s: %w", name, err)
			}
			fmt.Printf("Deleted: %s\n", name)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolP("purge", "p", false, "Purge the instance immediately")
	deleteCmd.Flags().BoolP("force", "f", false, "Force deletion without confirmation")
}