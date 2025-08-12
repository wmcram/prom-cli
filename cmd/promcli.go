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
		Description:           `promcli is an easy-to-use command-line tool for interacting with Prometheus endpoints. Its output is explicitly designed to be human-readable and easier to filter than grep'ing and awk'ing your way through curl output.`,
		Usage:                 "a command line tool for working with prometheus endpoints.",
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
