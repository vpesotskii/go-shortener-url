package config

import (
	"flag"
	"os"
)

var Options struct {
	Server      string
	BaseAddress string
}

func ParseFlags() {
	flag.StringVar(&Options.Server, "a", "localhost:8080", "address HTTP server")
	flag.StringVar(&Options.BaseAddress, "b", "http://localhost:8080", "Base address")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		Options.Server = envRunAddr
	}
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		Options.BaseAddress = envBaseAddr
	}
}
