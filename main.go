package main

import (
	"log"

	"github.com/uchimanajet7/gmth/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalf("\n%v\n", err)
	}
}
