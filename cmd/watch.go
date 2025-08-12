package main

import (
	"context"
	"errors"
	"time"
	"os"

	"github.com/guptarohit/asciigraph"
	"github.com/urfave/cli/v3"
	"github.com/wmcram/prom-cli/internal/processing"
)

// The subcommand `promcli watch`. It queries a metrics endpoint on
// a fixed interval, outputting metrics with identical behavior to
// `promcli get`.
var watchCommand = &cli.Command{
	Name:    "watch",
	Aliases: []string{"w"},
	Usage:   "watch an endpoint for live metrics",
	Flags: append(filterFlags,
		&cli.DurationFlag{
			Name:    "interval",
			Aliases: []string{"i"},
			Usage:   "interval to watch at",
			Value:   5 * time.Second,
		}),
	Action: func(ctx context.Context, cmd *cli.Command) error {
		var endpoint string
		if os.Getenv(endpointEnv) != "" {
			endpoint = os.Getenv(endpointEnv)
		} else if cmd.NArg() == 1 {
			endpoint = cmd.Args().Get(0)
		} else {
			return errors.New("usage: promcli get [FLAGS] ENDPOINT")
		}
		filters := processing.NewFilters(cmd.String("name"), cmd.String("label"), cmd.String("type"))
		ticker := time.NewTicker(cmd.Duration("interval"))
		defer ticker.Stop()

		asciigraph.Clear()
		if err := getAndDisplayMetrics(endpoint, filters); err != nil {
			return err
		}

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				asciigraph.Clear()
				if err := getAndDisplayMetrics(endpoint, filters); err != nil {
					return err
				}
			}
		}
	},
}
