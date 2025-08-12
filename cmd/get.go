package main

import (
	"context"
	"errors"

	"github.com/urfave/cli/v3"
	"github.com/wmcram/prom-cli/internal/display"
	"github.com/wmcram/prom-cli/internal/gather"
	"github.com/wmcram/prom-cli/internal/processing"
)

// A set of command-line flags for filtering on a list of Prometheus metrics.
// By filtering, we mean keeping only those metrics which match the flags.
var filterFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "name",
		Aliases: []string{"n"},
		Usage:   "Filter metrics by name: NAME1,NAME2,...",
	},
	&cli.StringFlag{
		Name:    "label",
		Aliases: []string{"l"},
		Usage:   "Filter metrics by label: LABEL1=VAL1,LABEL2=VAL2,...",
	},
	&cli.StringFlag{
		Name:    "type",
		Aliases: []string{"t"},
		Usage:   "Filter metrics by type: TYPE1,TYPE2,...",
	},
}

// The subcommand `promcli get`. It queries a prometheus endpoint,
// pretty-printing the metrics after they are filtered by the command
// flags.
var getCommand = &cli.Command{
	Name:    "get",
	Aliases: []string{"g"},
	Usage:   "show metrics from an endpoint",
	Flags:   filterFlags,
	Action: func(ctx context.Context, cmd *cli.Command) error {
		if cmd.NArg() != 1 {
			return errors.New("usage: promcli get [FLAGS] ENDPOINT")
		}
		endpoint := cmd.Args().Get(0)
		filters := processing.NewFilters(cmd.String("name"), cmd.String("label"), cmd.String("type"))
		return getAndDisplayMetrics(endpoint, filters)
	},
}

// getAndDisplayMetrics is a helper function for shared behavior between 
// `promcli get` and `promcli watch`. It queries a metric endpoint, applies
// the filters to the decoded metrics, and prints them to the terminal.
func getAndDisplayMetrics(endpoint string, filters *processing.Filters) error{
	decoder, err := gather.DecoderFromEndpoint(endpoint)
		if err != nil {
			return err
		}
	return display.DisplayMetrics(decoder, filters)
}
