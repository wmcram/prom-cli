package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/common/expfmt"
	display "github.com/wmcram/prom-cli"
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
	display.DisplayMetrics(decoder)
}