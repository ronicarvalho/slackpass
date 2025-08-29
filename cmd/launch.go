package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/slackpass/slackpass/internal/vm"
)

// launchCmd represents the launch command
var launchCmd = &cobra.Command{
	Use:   "launch [image] [name]",
	Short: "Create and start a new virtual machine",
	Long: `Launch a new virtual machine instance from a specified image.

Supported distributions:
  - slackware (latest, 15.0, 14.2)
  - debian (bookworm, bullseye, buster)
  - fedora (39, 38, 37)
  - almalinux (9, 8)
  - rockylinux (9, 8)
  - centos (stream9, stream8, 7)
  - gentoo (latest)
  - opensuse (tumbleweed, leap)

Examples:
  slackpass launch                    # Launch default image with auto-generated name
  slackpass launch debian             # Launch latest Debian with auto-generated name
  slackpass launch debian:bookworm    # Launch Debian Bookworm
  slackpass launch debian myvm        # Launch Debian with name 'myvm'
  slackpass launch debian:bookworm myvm --cpus 2 --memory 2G --disk 20G`,
	Args: cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		image := "debian:bookworm" // default image
		name := ""

		if len(args) >= 1 {
			image = args[0]
		}
		if len(args) >= 2 {
			name = args[1]
		}

		// Get flags
		cpus, _ := cmd.Flags().GetInt("cpus")
		memory, _ := cmd.Flags().GetString("memory")
		disk, _ := cmd.Flags().GetString("disk")
		cloudInit, _ := cmd.Flags().GetString("cloud-init")

		// Validate image format
		if !isValidImage(image) {
			return fmt.Errorf("invalid image format: %s", image)
		}

		config := &vm.LaunchConfig{
			Image:     image,
			Name:      name,
			CPUs:      cpus,
			Memory:    memory,
			Disk:      disk,
			CloudInit: cloudInit,
		}

		manager := vm.NewManager()
		return manager.Launch(config)
	},
}

func init() {
	rootCmd.AddCommand(launchCmd)

	launchCmd.Flags().IntP("cpus", "c", 1, "Number of CPUs")
	launchCmd.Flags().StringP("memory", "m", "1G", "Amount of memory")
	launchCmd.Flags().StringP("disk", "d", "10G", "Disk size")
	launchCmd.Flags().String("cloud-init", "", "Path to cloud-init configuration file")
}

func isValidImage(image string) bool {
	validDistros := []string{
		"slackware", "debian", "fedora", "almalinux",
		"rockylinux", "centos", "gentoo", "opensuse",
	}

	parts := strings.Split(image, ":")
	distro := parts[0]

	for _, valid := range validDistros {
		if distro == valid {
			return true
		}
	}
	return false
}