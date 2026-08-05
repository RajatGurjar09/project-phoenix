package config

import "os"

const (
	DefaultAddress = ":8080"
	AddressEnv     = "PHOENIX_HTTP_ADDR"
)

type Config struct {
	Address string
}

func Load() Config {
	address := os.Getenv(AddressEnv)
	if address == "" {
		address = DefaultAddress
	}

	return Config{Address: address}
}
