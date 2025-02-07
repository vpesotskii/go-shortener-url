package main

import (
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/vpesotskii/go-shortener-url/cmd/config"
	"github.com/vpesotskii/go-shortener-url/internal/app/compress"
	"github.com/vpesotskii/go-shortener-url/internal/app/handlers"
	"github.com/vpesotskii/go-shortener-url/internal/app/logger"
	"github.com/vpesotskii/go-shortener-url/internal/app/models"
	"github.com/vpesotskii/go-shortener-url/internal/app/storage"
	"go.uber.org/zap"
	"net/http"
	"os"
)

func main() {

	var db storage.Repository

	//Обертки для handlers, чтобы использовать их в роутере
	AddURLHandlerWrapper := func(res http.ResponseWriter, req *http.Request) {
		handlers.AddURL(db, res, req)
	}

	AddURLFromJSONHandlerWrapper := func(res http.ResponseWriter, req *http.Request) {
		handlers.AddURLFromJSON(db, res, req)
	}

	GetURLHandlerWrapper := func(res http.ResponseWriter, req *http.Request) {
		handlers.GetURL(db, res, req)
	}

	PingHandlerWrapper := func(res http.ResponseWriter, req *http.Request) {
		handlers.Ping(db, res, req)
	}

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/{surl}", logger.WithLogger(compress.GzipMiddleware(GetURLHandlerWrapper)))
		r.Post("/", logger.WithLogger(compress.GzipMiddleware(AddURLHandlerWrapper)))
		r.Post("/api/shorten", logger.WithLogger(compress.GzipMiddleware(AddURLFromJSONHandlerWrapper)))
		r.Get("/ping", logger.WithLogger(compress.GzipMiddleware(PingHandlerWrapper)))
	})
	config.ParseFlags()
	err := logger.Initialize(config.Options.LogLevel)
	if err != nil {
		return
	}
	logger.Log.Info("Running server on", zap.String("server", config.Options.Server))
	logger.Log.Info("Base address", zap.String("base address", config.Options.BaseAddress))
	logger.Log.Info("File Storage Path", zap.String("file", config.Options.FileStorage))
	logger.Log.Info("Database Connection", zap.String("db connection", config.Options.DBUrl))

	switch config.Options.DBUrl {
	case "":
		file, err := os.OpenFile(config.Options.FileStorage, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			logger.Log.Fatal("Error during opening file with shorten urls: %v", zap.Error(err))
		}
		database := storage.NewStorage(map[string]models.URL{})
		database.SetFile(file)

		err = database.FillFromFile(file)
		if err != nil {
			logger.Log.Info("Error during filling file with shorten urls: %v", zap.Error(err))
		}
		file.Close()

		db = database

	default:
		database, err := storage.NewDatabase(config.Options.DBUrl)
		if err != nil {
			logger.Log.Fatal("Error during creating database with shorten urls: %v", zap.Error(err))
		}
		defer database.DB.Close()

		db = database
	}

	err = http.ListenAndServe(config.Options.Server, r)
	if err != nil {
		return
	}
}
