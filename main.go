package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/indeedsecurity/carbonbeat/v2/app"
)

func main() {
	err := beat.Run("carbonbeat", "", app.New)
	if err != nil {
		os.Exit(1)
	}
}
