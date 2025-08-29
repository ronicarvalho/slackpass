package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/slackpass/slackpass/internal/vm"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all virtual machine instances",
	Long: `List all virtual machine instances managed by Slackpass.

Shows the name, state, IPv4 address, image, and resource allocation
for each virtual machine.

Examples:
  slackpass list
  slackpass ls`,
	RunE: func(cmd *cobra.Command, args []string) error {
		manager := vm.NewManager()
		instances, err := manager.List()
		if err != nil {
			return fmt.Errorf("failed to list instances: %w", err)
		}

		if len(instances) == 0 {
			fmt.Println("No instances found.")
			return nil
		}

		// Create tabwriter for formatted output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Name\tState\tIPv4\tImage\tCPUs\tMemory\tDisk")

		for _, instance := range instances {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%s\t%s\n",
				instance.Name,
				instance.State,
				instance.IPv4,
				instance.Image,
				instance.CPUs,
				instance.Memory,
				instance.Disk,
			)
		}

		return w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}