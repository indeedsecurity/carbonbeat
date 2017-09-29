// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

// Config has custom options for Carbonbeat
type Config struct {
	Period time.Duration `config:"period"`
	Type   string

	API struct {
		ID   string
		Key  string
		Host string
	}
}

// DefaultConfig contains defaults for custom options
var DefaultConfig = Config{
	Period: 5 * time.Minute,
	Type:   "cb",
}
