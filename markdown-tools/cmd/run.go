package cmd

import (
	converter "github.com/lubieniebieski/scripts/markdown_links_as_references/pkg"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Replace all inline links in a Markdown file",
	Long:  `It can change either one file or many`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		converter.Run(args[0])
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
