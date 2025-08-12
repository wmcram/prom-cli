package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:                  "promcli",
		Description:           "A command line tool for working with prometheus endpoints.",
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			getCommand,
			watchCommand,
			mockCommand,
			graphCommand,
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
