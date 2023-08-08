package cmd

import (
	converter "github.com/lubieniebieski/markdown-tools/pkg"

	"github.com/spf13/cobra"
)

var createBackup bool
var verbose bool

var linksAsReferencesCmd = &cobra.Command{
	Use:   "links_as_references",
	Short: "Replace all inline links in a Markdown file(s)",
	Long:  `It can change either one file or many, you can provide a single file name or entire directory - it will process all files with .md extension`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		converter.ConvertFilesInPath(args[0], createBackup, verbose)
	},
}

func init() {
	linksAsReferencesCmd.Flags().BoolVarP(&createBackup, "backup", "b", false, "Create backup file(s)")
	linksAsReferencesCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print verbose output")

	rootCmd.AddCommand(linksAsReferencesCmd)
}
