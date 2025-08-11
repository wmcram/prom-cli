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
	endpointColor = color.New(color.FgHiMagenta, color.Bold).PrintfFunc()
	titleColor    = color.New(color.FgHiCyan).PrintfFunc()
	valueColor    = color.New(color.FgHiGreen).PrintfFunc()
	labelColor    = color.New(color.FgYellow).PrintfFunc()
)

// DisplayMetrics prints the metrics to the screen after filtering.
func DisplayMetrics(decoder expfmt.Decoder, endpoint string, filters *Filters) {
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

type Filters struct {
	nameFilter map[string]bool
	labelFilter map[string]string
	typeFilter map[string]bool
}

func NewFilters(nameFilter string, labelFilter string, typeFilter string) *Filters {
	return &Filters{
		nameFilter: processNameFilter(nameFilter),
		labelFilter: processLabelFilter(labelFilter),
		typeFilter: processTypeFilter(typeFilter),
	}
}

// MatchesMetricFamily determines whether this metric family passes the filter
func (f *Filters) MatchesMetricFamily(mf *dto.MetricFamily) bool {
	if f.nameFilter != nil && !f.nameFilter[mf.GetName()] {
		return false
	}
	if f.typeFilter != nil && !f.typeFilter[mf.GetType().String()] {
		return false
	}
	return true
}

// MatchesMetric determines whether this metric passes the filter
func (f *Filters) MatchesMetric(m *dto.Metric) bool {
	goal := len(f.labelFilter)
	for _, labelPair := range m.GetLabel() {
		if f.labelFilter[labelPair.GetName()] == labelPair.GetValue() {
			goal--
		}
	}
	return goal == 0
}

// processNameFilter parses a cli flag into a map of name filters
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

// processTypeFilter parses a cli flag into a map of type filters
func processTypeFilter(filter string) map[string]bool {
	if filter == "" {
		return nil
	}
	m := make(map[string]bool)
	for _, kind := range strings.Split(filter, ", ") {
		m[kind] = true
	}
	return m
}

// processLabelFilter parses a cli flag into a map of label filters
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
