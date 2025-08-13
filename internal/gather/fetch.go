package gather

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/wmcram/prom-cli/internal/processing"
)

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

// DecoderFromEndpoint makes an http request to an endpoint, turning the response into a Prometheus Decoder.
func DecoderFromEndpoint(endpoint string) (expfmt.Decoder, error) {
	url, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	request := http.Request{
		Method: http.MethodGet,
		URL:    url,
		Header: make(map[string][]string),
	}
	// attempt to negotiate protobufs with endpoint
	request.Header.Set("Accept", "application/vnd.google.protobuf; proto=io.prometheus.client.MetricFamily; encoding=delimited")

	resp, err := httpClient.Do(&request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	decoder := expfmt.NewDecoder(resp.Body, expfmt.ResponseFormat(resp.Header))
	return decoder, nil
}

// GetMetricValue gets the value of a uniquely-identified metric from an endpoint.
// It throws an error if the metric could not be found, or if the query did not
// return a unique metric.
func GetMetricValue(endpoint string, filters *processing.Filters) (float64, string, error) {
	decoder, err := DecoderFromEndpoint(endpoint)
	if err != nil {
		return 0, "", err
	}

	for {
		mf := &dto.MetricFamily{}
		err := decoder.Decode(mf)
		// We don't care about whether there was an actual error or EOF here, either way we couldn't get the metric.
		if err != nil {
			return 0, "", err
		}
		if !filters.MatchesMetricFamily(mf) {
			continue
		}

		// Beyond this point, we know that we will be returning in this iteration.
		metricType := mf.GetType()
		metricName := mf.GetName()
		res := []float64{}
		for _, metric := range mf.Metric {
			if !filters.MatchesMetric(metric) {
				continue
			}
			switch metricType {
			// only support counter and gauge for now
			case dto.MetricType_COUNTER:
				res = append(res, metric.GetCounter().GetValue())
			case dto.MetricType_GAUGE:
				res = append(res, metric.GetGauge().GetValue())
			}
		}
		if len(res) != 1 {
			return 0, "", fmt.Errorf("query did not return a unique metric")
		}
		return res[0], metricName, nil
	}
}
