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
		Description: "A command line tool for working with prometheus endpoints.",
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			// promcli get
			{
				Name: "get",
				Aliases: []string{"g"},
				Usage: "show metrics from an endpoint",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "name",
						Aliases: []string{"n"},
						Usage: "Filter metrics by name: NAME1,NAME2,...",
					},
					&cli.StringFlag{
						Name: "label",
						Aliases: []string{"l"},
						Usage: "Filter metrics by label: LABEL1=VAL1,LABEL2=VAL2,...",
					},
					&cli.StringFlag{
						Name: "type",
						Aliases: []string{"t"},
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
				Name: "mock",
				Aliases: []string{"m"},
				Usage: "mock a prometheus endpoint for testing",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name: "port",
						Aliases: []string{"p"},
						Usage: "port to serve on",
						Value: 8080,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.NArg() != 1 {
						return errors.New("usage: promcli mock [FLAGS] FILE")
					}
					return mockFile(cmd.Args().Get(0), cmd.Int("port"))
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

// getEndpoint pretty-prints the metrics from a given endpoint, filtered by the given filters.
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

// mockFile runs a prometheus endpoint serving metrics from the given file.
func mockFile(file string, port int) error {
	dat, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	
	http.HandleFunc("/metrics", func (w http.ResponseWriter, r *http.Request)  {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		w.WriteHeader(http.StatusOK)
		w.Write(dat)
	})
	fmt.Println("Serving mock metrics on port", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}



