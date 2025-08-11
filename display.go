package display

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"golang.org/x/term"
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
	endpointColor = color.New(color.FgMagenta, color.Bold).PrintfFunc()
	titleColor    = color.New(color.FgCyan).PrintfFunc()
	valueColor    = color.New(color.FgGreen).PrintfFunc()
	labelColor    = color.New(color.FgYellow).PrintfFunc()
)

// DisplayMetrics prints the metrics to the screen after filtering.
func DisplayMetrics(decoder expfmt.Decoder, endpoint string, nameFilter map[string]bool, labelFilter map[string]string) {
	endpointColor("Endpoint %s\n", endpoint)
	for {
		mf := dto.MetricFamily{}
		err := decoder.Decode(&mf)
		if err != nil {
			break
		}
		if nameFilter != nil && !nameFilter[mf.GetName()] {
			continue
		}

		metricType := mf.GetType()
		titleColor("%s (%s): %s\n", mf.GetName(), metricType, mf.GetHelp())

		for _, metric := range mf.Metric {
			if labelFilter != nil && !matchesLabelFilter(metric, labelFilter) {
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

// matchesLabelFilter determines whether this metric is matched by the given labelFilter
func matchesLabelFilter(metric *dto.Metric, labelFilter map[string]string) bool {
	goal := len(labelFilter)
	for _, labelPair := range metric.GetLabel() {
		if labelFilter[labelPair.GetName()] == labelPair.GetValue() {
			goal--
		}
	}
	return goal == 0
}
