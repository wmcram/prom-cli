package main

import (
	"os"
	"fmt"
	"net/http"
	"github.com/prometheus/common/expfmt"
	dto "github.com/prometheus/client_model/go"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: promcli <endpoint>")
		return
	}
	endpoint := os.Args[1]

	resp, err := http.Get(endpoint)
	if err != nil {
		fmt.Println("Error reaching endpoint:", err)
		return
	}
	defer resp.Body.Close()

	decoder := expfmt.NewDecoder(resp.Body, expfmt.NewFormat(expfmt.TypeTextPlain))
	mf := dto.MetricFamily{}

	for {
		err := decoder.Decode(&mf)
		if err != nil {
			break
		}
		metricType := mf.GetType()
		fmt.Printf("%s\n", mf.GetName())
		for _, metric := range mf.Metric {
			switch metricType {
			// only support counter and gauge for now
			case dto.MetricType_COUNTER:
				fmt.Printf("\t%f\n", metric.GetCounter().GetValue())
			case dto.MetricType_GAUGE:
				fmt.Printf("\t%f\n", metric.GetGauge().GetValue())
			}
		}
	}
}