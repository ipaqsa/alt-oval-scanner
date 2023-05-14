package main

import (
	"alt-oval-scanner/pkg/command"
	"log"
)

func main() {
	err := command.RootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
