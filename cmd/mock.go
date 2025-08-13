package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/urfave/cli/v3"
)

// The subcommand `promcli mock`. It takes a file containing some metrics in plaintext format,
// and serves them on a specified port.
var mockCommand = &cli.Command{
	Name:    "mock",
	Aliases: []string{"m"},
	Usage:   "mock a prometheus endpoint for testing",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Usage:   "port to serve on",
			Value:   8080,
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		if cmd.NArg() != 1 {
			return errors.New("usage: promcli mock [FLAGS] FILE")
		}
		file := cmd.Args().First()
		dat, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain; version=0.0.4")
			w.WriteHeader(http.StatusOK)
			w.Write(dat)
		})

		port := cmd.Int("port")
		fmt.Println("Serving mock metrics on port", port)
		return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	},
}
