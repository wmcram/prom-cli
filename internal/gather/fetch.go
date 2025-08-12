package gather

import (
	"fmt"
	"net/http"

	"github.com/prometheus/common/expfmt"
)

// DecoderFromEndpoint turns the response from an endpoint into a prometheus metric decoder.
func DecoderFromEndpoint(endpoint string) (expfmt.Decoder, error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}
	decoder := expfmt.NewDecoder(resp.Body, expfmt.NewFormat(expfmt.TypeTextPlain))
	return decoder, nil
}
