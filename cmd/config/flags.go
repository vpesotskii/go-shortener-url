package config

import (
	"flag"
)

var Options struct {
	Server      string
	BaseAddress string
}

func ParseFlags() {
	flag.StringVar(&Options.Server, "a", "localhost:8888", "address HTTP server")
	flag.StringVar(&Options.BaseAddress, "b", "http://localhost:8000", "Base address")
	flag.Parse()
}