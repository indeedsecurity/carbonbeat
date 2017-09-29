package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/indeedsecurity/carbonbeat/beater"
)

func main() {
	err := beat.Run("carbonbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
