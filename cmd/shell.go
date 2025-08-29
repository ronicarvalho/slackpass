package cmd

import (
	"github.com/spf13/cobra"
	"github.com/slackpass/slackpass/internal/vm"
)

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:   "shell [name]",
	Short: "Open a shell session to a virtual machine",
	Long: `Open an interactive shell session to a virtual machine via SSH.

If no name is provided and only one instance is running, it will
connect to that instance automatically.

Examples:
  slackpass shell myvm
  slackpass shell          # Connect to the only running instance`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var name string
		if len(args) > 0 {
			name = args[0]
		}

		manager := vm.NewManager()
		return manager.Shell(name)
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}