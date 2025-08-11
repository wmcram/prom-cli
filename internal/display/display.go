package display

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"golang.org/x/term"
	"github.com/wmcram/prom-cli/internal/processing"
)

var TermWidth, TermHeight int

// init gets the terminal dimensions.
func init() {
	var err error
	TermWidth, TermHeight, err = term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		TermWidth = 80
		TermHeight = 24
	}
}

var (
	endpointColor = color.New(color.FgHiMagenta, color.Bold).PrintfFunc()
	titleColor    = color.New(color.FgHiCyan).PrintfFunc()
	valueColor    = color.New(color.FgHiGreen).PrintfFunc()
	labelColor    = color.New(color.FgYellow).PrintfFunc()
)

// DisplayMetrics prints the metrics to the screen after filtering.
func DisplayMetrics(decoder expfmt.Decoder, endpoint string, filters *processing.Filters) {
	endpointColor("Endpoint %s\n", endpoint)
	for {
		mf := &dto.MetricFamily{}
		err := decoder.Decode(mf)
		if err != nil {
			break
		}
		if !filters.MatchesMetricFamily(mf) {
			continue
		}

		metricType := mf.GetType()
		titleColor("%s (%s): %s\n", mf.GetName(), metricType, mf.GetHelp())
		for _, metric := range mf.Metric {
			if !filters.MatchesMetric(metric) {
				continue
			}

			var labelStrings = []string{}
			for _, labelPair := range metric.GetLabel() {
				labelStrings = append(labelStrings, fmt.Sprintf("%s=\"%s\"", labelPair.GetName(), labelPair.GetValue()))
			}
			labels := strings.Join(labelStrings, ",")
			labelColor("%s", labels)

			switch metricType {
			// only support counter and gauge for now
			case dto.MetricType_COUNTER:
				valueColor(" %.0f\n", metric.GetCounter().GetValue())
			case dto.MetricType_GAUGE:
				valueColor(" %.2f\n", metric.GetGauge().GetValue())
			}
			
		}
	}
}
