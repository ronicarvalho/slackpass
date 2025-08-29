package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/slackpass/slackpass/internal/images"
)

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find [image-name]",
	Short: "Display available images to launch",
	Long: `Display available cloud images that can be used to launch instances.

Shows information about supported Linux distributions and their versions.

Examples:
  slackpass find
  slackpass find debian
  slackpass find --format table`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")
		remoteOnly, _ := cmd.Flags().GetBool("remote-only")

		var filter string
		if len(args) > 0 {
			filter = args[0]
		}

		imageManager := images.NewManager()
		availableImages, err := imageManager.Find(filter, remoteOnly)
		if err != nil {
			return fmt.Errorf("failed to find images: %w", err)
		}

		if len(availableImages) == 0 {
			fmt.Println("No images found.")
			return nil
		}

		switch format {
		case "table":
			return displayImagesTable(availableImages)
		case "json":
			return displayImagesJSON(availableImages)
		default:
			return displayImagesTable(availableImages)
		}
	},
}

func init() {
	rootCmd.AddCommand(findCmd)

	findCmd.Flags().String("format", "table", "Output format: table, json")
	findCmd.Flags().Bool("remote-only", false, "Show only remote images")
}

func displayImagesTable(images []*images.ImageInfo) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Image\tAliases\tVersion\tDescription")

	for _, img := range images {
		fmt.Fprintf(w, "%s:%s\t%s\t%s\t%s\n",
			img.Distribution,
			img.Version,
			img.Aliases,
			img.Version,
			img.Description,
		)
	}

	return w.Flush()
}

func displayImagesJSON(images []*images.ImageInfo) error {
	// TODO: Implement JSON output
	fmt.Println("JSON output not implemented yet")
	return nil
}