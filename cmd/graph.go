package main

import (
	"context"
	"errors"
	"time"

	"github.com/urfave/cli/v3"
	"github.com/wmcram/prom-cli/internal/display"
)

var graphCommand = &cli.Command{
	Name:    "graph",
	Aliases: []string{"g"},
	Usage:   "graph metrics from an endpoint",
	Flags: []cli.Flag{
		&cli.DurationFlag{
			Name:    "interval",
			Aliases: []string{"i"},
			Usage:   "interval to watch at",
			Value:   5 * time.Second,
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		if cmd.NArg() != 2 {
			return errors.New("usage: promcli graph ENDPOINT METRIC_NAME")
		}
		endpoint := cmd.Args().Get(0)
		metricName := cmd.Args().Get(1)
		return display.GraphMetric(ctx, endpoint, metricName, cmd.Duration("interval"))
	},
}
