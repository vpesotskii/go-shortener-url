package config

import (
	"flag"
	"os"
)

var Options struct {
	Server      string
	BaseAddress string
	LogLevel    string
}

func ParseFlags() {
	flag.StringVar(&Options.Server, "a", "localhost:8080", "address HTTP server")
	flag.StringVar(&Options.BaseAddress, "b", "http://localhost:8080", "Base address")
	flag.StringVar(&Options.LogLevel, "l", "info", "Log level")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		Options.Server = envRunAddr
	}
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		Options.BaseAddress = envBaseAddr
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		Options.LogLevel = envLogLevel
	}
}
