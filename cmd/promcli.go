package main

import (
	"context"
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
			{
				Name: "get",
				Aliases: []string{"g"},
				Usage: "get metrics from an endpoint",
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
					endpoint := cmd.Args().Get(0)
					if endpoint == "" {
						return fmt.Errorf("endpoint is required")
					}
					filters := processing.NewFilters(cmd.String("name"), cmd.String("label"), cmd.String("type"))
					return getEndpoint(endpoint, filters)
				},
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

	decoder := expfmt.NewDecoder(resp.Body, expfmt.NewFormat(expfmt.TypeTextPlain))
	display.DisplayMetrics(decoder, filters)
	return nil
}



