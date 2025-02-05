package config

import (
	"flag"
	"os"
)

var Options struct {
	Server      string
	BaseAddress string
	LogLevel    string
	FileStorage string
	DBUrl       string
}

func ParseFlags() {
	flag.StringVar(&Options.Server, "a", "localhost:8080", "address HTTP server")
	flag.StringVar(&Options.BaseAddress, "b", "http://localhost:8080", "Base address")
	flag.StringVar(&Options.LogLevel, "l", "info", "Log level")
	flag.StringVar(&Options.FileStorage, "f", "/tmp/short-url-db.json", "File storage location")
	flag.StringVar(&Options.DBUrl, "d", "host=localhost port=5432 user=postgres password=admin dbname=go sslmode=disable", "Database address")
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
	if envFileStorage := os.Getenv("FILE_STORAGE_PATH"); envFileStorage != "" {
		Options.FileStorage = envFileStorage
	}
	if envDBUrl := os.Getenv("DATABASE_DSN"); envDBUrl != "" {
		Options.DBUrl = envDBUrl
	}
}
