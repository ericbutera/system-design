package env

import (
	"github.com/caarlos0/env/v11"
)

// Unmarshal environment into a struct using caarlos0/env.
// Does not require an init() function to set defaults.
func New[T any]() (*T, error) {
	var config T
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	cfg, err := env.ParseAs[T]()
	return &cfg, err
}
