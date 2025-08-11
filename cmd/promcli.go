package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/prometheus/common/expfmt"
	display "github.com/wmcram/prom-cli"
)

func main() {
	var nameFilter *string = flag.String("name", "", "Filter metrics by name: NAME1,NAME2,...")
	var labelFilter *string = flag.String("label", "", "Filter metrics by label: LABEL1=VAL1,LABEL2=VAL2,...")
	var typeFilter *string = flag.String("type", "", "Filter metrics by type: TYPE1,TYPE2,...")

	flag.Parse()
	filters := display.NewFilters(*nameFilter, *labelFilter, *typeFilter)

	endpoints := flag.Args()
	if len(endpoints) == 0 {
		fmt.Println("Usage: promcli <endpoints>")
		return
	}

	for _, endpoint := range endpoints {
		resp, err := http.Get(endpoint)
		if err != nil {
			fmt.Println("Error reaching endpoint:", err)
			return
		}
		defer resp.Body.Close()

		decoder := expfmt.NewDecoder(resp.Body, expfmt.NewFormat(expfmt.TypeTextPlain))
		display.DisplayMetrics(decoder, endpoint, filters)
	}
}



