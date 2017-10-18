package cmd

import (
	"errors"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// holds value of "file" flag
var rootMarkdownFile string

// RootCmd defines 'Golang Markdown to HTML(gmth)' command
var RootCmd = &cobra.Command{
	Use:           "gmth",
	Short:         "Golang Markdown to HTML(gmth) command",
	Long:          `A command line tool converts Markdown file to HTML file`,
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cmdName := strings.ToLower(cmd.Name())
		if cmdName == "help" || cmdName == "version" {
			return nil
		}

		if rootMarkdownFile == "" {
			// cmd.Help()
			return errors.New("Please specify flag [--file] required markdown file.\n")
		}

		abspath, err := filepath.Abs(rootMarkdownFile)
		if err != nil {
			return err
		}
		rootMarkdownFile = abspath

		return nil
	},
}

func init() {
	cobra.OnInitialize(initRootConfig)
	RootCmd.PersistentFlags().StringVarP(&rootMarkdownFile, "file", "f", "", "Specify markdown file. [required]")
}

func initRootConfig() {
	if _, err := getConfig(); err != nil {
		log.Fatalf("\nCan't read config file: %s\n\n", err)
	}
}
