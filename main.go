package main

import (
	"log"

	"./cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalf("\n%v\n", err)
	}
}
