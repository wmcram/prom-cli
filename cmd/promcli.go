package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/prometheus/common/expfmt"
	display "github.com/wmcram/prom-cli"
)



func main() {
	var nameFilter *string = flag.String("name", "", "Filter metrics by name: NAME1,NAME2,...")
	var labelFilter *string = flag.String("label", "", "Filter metrics by label: LABEL1=VAL1,LABEL2=VAL2,...")

	flag.Parse()

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
		display.DisplayMetrics(decoder, endpoint, processNameFilter(*nameFilter), processLabelFilter(*labelFilter))
	}
}

func processNameFilter(filter string) map[string]bool {
	if filter == "" {
		return nil
	}
	m := make(map[string]bool)
	for _, name := range strings.Split(filter, ",") {
		m[name] = true
	}
	return m
}

func processLabelFilter(filter string) map[string]string {
	if filter == "" {
		return nil
	}
	m := make(map[string]string)
	for pair := range strings.SplitSeq(filter, ", ") {
		parts := strings.Split(pair, "=")
		if len(parts) != 2 {
			continue
		}
		m[parts[0]] = parts[1]
	}
	return m
}

