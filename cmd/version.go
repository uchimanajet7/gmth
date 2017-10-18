package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version number is set automatically at build time
var version = "0.0.1"

// revision is set automatically use "git rev-parse --short HEAD" at build time
var revision = "devel"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of gmth",
	Long:  `All software has versions. This is gmth's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gmth %s (-- HEAD rev: %s)\n", version, revision)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
