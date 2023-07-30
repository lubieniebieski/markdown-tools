package cmd

import (
	converter "github.com/lubieniebieski/scripts/markdown-tools/pkg"

	"github.com/spf13/cobra"
)

var linksAsReferencesCmd = &cobra.Command{
	Use:   "links_as_references",
	Short: "Replace all inline links in a Markdown file",
	Long:  `It can change either one file or many`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		converter.Run(args[0])
	},
}

func init() {
	rootCmd.AddCommand(linksAsReferencesCmd)
}
