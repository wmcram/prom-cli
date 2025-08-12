package main

import (
	"context"
	"errors"
	"time"

	"github.com/urfave/cli/v3"
	"github.com/wmcram/prom-cli/internal/display"
	"github.com/wmcram/prom-cli/internal/processing"
)

var graphCommand = &cli.Command{
	Name:    "graph",
	Aliases: []string{"g"},
	Usage:   "graph metrics from an endpoint. the metric must be uniquely determined by the filters!",
	Flags: append(filterFlags,
		&cli.DurationFlag{
			Name:    "interval",
			Aliases: []string{"i"},
			Usage:   "interval to watch at",
			Value:   5 * time.Second,
	}),
	Action: func(ctx context.Context, cmd *cli.Command) error {
		if cmd.NArg() < 1 {
			return errors.New("usage: promcli graph [FLAGS] endpoint")
		}
		endpoint := cmd.Args().Get(0)
		filters := processing.NewFilters(cmd.String("name"), cmd.String("labels"), cmd.String("type"))
		return display.GraphMetric(ctx, endpoint, filters, cmd.Duration("interval"))
	},
}
