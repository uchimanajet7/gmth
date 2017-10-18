package cmd

import (
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Markdown Convert the file to HTML file once",
	Long:  `Execute conversion of Markdown file to HTML file once according to setting file.`,
	RunE:  execRunCmd,
}

func init() {
	RootCmd.AddCommand(runCmd)
}

func execRunCmd(cmd *cobra.Command, args []string) error {
	config, err := getConfig()
	if err != nil {
		return err
	}
	return convertFile(rootMarkdownFile, config)
}
