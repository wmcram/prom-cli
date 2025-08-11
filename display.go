package display

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

const TermWidth = 80

func DisplayMetrics(decoder expfmt.Decoder) {
	for {
		mf := dto.MetricFamily{}
		err := decoder.Decode(&mf)
		if err != nil {
			break
		}
		metricType := mf.GetType()
		color.Cyan("%s (%s): %s", mf.GetName(), metricType, mf.GetHelp())
		for _, metric := range mf.Metric {
			color.Cyan(strings.Repeat("-", TermWidth))
			var labelStrings = []string{}
			for _, labelPair := range metric.GetLabel() {
				labelStrings = append(labelStrings, fmt.Sprintf("%s=\"%s\"", labelPair.GetName(), labelPair.GetValue()))
			}
			labels := strings.Join(labelStrings, ",")

			switch metricType {
			// only support counter and gauge for now
			case dto.MetricType_COUNTER:
				color.Green("> %.0f", metric.GetCounter().GetValue())
			case dto.MetricType_GAUGE:
				color.Green("> %.2f", metric.GetGauge().GetValue())
			}
			color.Yellow(labels)
		}
		color.Cyan(strings.Repeat("-", TermWidth))
	}
}
