package main

import (
	"context"
	"errors"
	"time"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/wmcram/prom-cli/internal/display"
	"github.com/wmcram/prom-cli/internal/processing"
)

var graphCommand = &cli.Command{
	Name:    "graph",
	Aliases: []string{"g"},
	Usage:   "graph a metric from an endpoint. the metric must be uniquely determined by the filters!",
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
		filters := processing.NewFilters(cmd.String("name"), cmd.String("labels"), cmd.String("type"))
		return display.GraphMetric(ctx, endpoint, filters, cmd.Duration("interval"))
	},
}
