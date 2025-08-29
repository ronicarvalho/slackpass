package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/slackpass/slackpass/internal/vm"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec [name] -- [command]",
	Short: "Execute a command on a virtual machine",
	Long: `Execute a command on a virtual machine via SSH.

The command and its arguments should be separated from the instance
name by '--'.

Examples:
  slackpass exec myvm -- ls -la
  slackpass exec myvm -- "echo 'Hello World'"
  slackpass exec myvm -- sudo apt update`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 || args[1] != "--" {
			return fmt.Errorf("command must be separated by '--'")
		}

		name := args[0]
		command := strings.Join(args[2:], " ")

		manager := vm.NewManager()
		return manager.Exec(name, command)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}