package main

import (
	"log"

	"github.com/ouijan/aether/pkg/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		log.Fatal(err)
	}
}
