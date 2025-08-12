package processing

import (
	dto "github.com/prometheus/client_model/go"
	"strings"
)

// Filters contains information about what kinds of metrics should be included in a query.
type Filters struct {
	nameFilter  map[string]bool
	labelFilter map[string]string
	typeFilter  map[string]bool
}

// NewFilters creates a Filters struct from the given cli flags
func NewFilters(nameFilter string, labelFilter string, typeFilter string) *Filters {
	return &Filters{
		nameFilter:  processNameFilter(nameFilter),
		labelFilter: processLabelFilter(labelFilter),
		typeFilter:  processTypeFilter(typeFilter),
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
	for name := range strings.SplitSeq(filter, ",") {
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
	for kind := range strings.SplitSeq(filter, ", ") {
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
