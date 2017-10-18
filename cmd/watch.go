package cmd

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch Markdown file and converts it to HTML file",
	Long:  `Watch Markdown file according to the settings and convert it to HTML file if there is a change.`,
	RunE:  execWatchCmd,
}

func init() {
	RootCmd.AddCommand(watchCmd)
}

func execWatchCmd(cmd *cobra.Command, args []string) error {
	config, err := getConfig()
	if err != nil {
		return err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		log.Printf("\nChange detection started: %+v\nPlease press 'Ctrl + C' to stop", rootMarkdownFile)
		for {
			select {
			case event := <-watcher.Events:
				log.Printf("\nEvent: %+v\n", event)

				switch {
				case event.Op&fsnotify.Write == fsnotify.Write:
					err := convertFile(rootMarkdownFile, config)
					if err != nil {
						log.Fatalf("%v\n", err)
					}
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					fmt.Println("Since file was deleted, wtach was aborted.")
					done <- true
				case event.Op&fsnotify.Rename == fsnotify.Rename:
					fmt.Println("Since file name was changed, wtach was aborted.")
					done <- true
				}
			case err := <-watcher.Errors:
				log.Printf("\nEvent: %+v\n", err)
				done <- true
			}
		}
	}()

	err = watcher.Add(rootMarkdownFile)
	if err != nil {
		return err
	}
	<-done

	return nil
}
