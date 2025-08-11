package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/common/expfmt"
	"github.com/urfave/cli/v3"
	"github.com/wmcram/prom-cli/internal/display"
	"github.com/wmcram/prom-cli/internal/processing"
)

func main() {
	cmd := &cli.Command{
		Commands: []*cli.Command{
			// promcli get
			{
				Name: "get",
				Aliases: []string{"g"},
				Usage: "show metrics from an endpoint",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "name",
						Usage: "Filter metrics by name: NAME1,NAME2,...",
					},
					&cli.StringFlag{
						Name: "label",
						Usage: "Filter metrics by label: LABEL1=VAL1,LABEL2=VAL2,...",
					},
					&cli.StringFlag{
						Name: "type",
						Usage: "Filter metrics by type: TYPE1,TYPE2,...",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.NArg() != 1 {
						return errors.New("usage: promcli get [FLAGS] ENDPOINT")
					}
					endpoint := cmd.Args().Get(0)
					filters := processing.NewFilters(cmd.String("name"), cmd.String("label"), cmd.String("type"))
					return getEndpoint(endpoint, filters)
				},
			},
			// promcli watch
			{
				Name: "watch",
				Aliases: []string{"w"},
				Usage: "watch an endpoint for live metrics",
				Flags: []cli.Flag{

				},
			},
			// promcli mock
			{

			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func getEndpoint(endpoint string, filters *processing.Filters) error {
	resp, err := http.Get(endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	decoder := expfmt.NewDecoder(resp.Body, expfmt.NewFormat(expfmt.TypeTextPlain))
	return display.DisplayMetrics(decoder, filters)
}



