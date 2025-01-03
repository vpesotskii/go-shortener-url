package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/vpesotskii/go-shortener-url/cmd/config"
	"github.com/vpesotskii/go-shortener-url/internal/app/compress"
	"github.com/vpesotskii/go-shortener-url/internal/app/logger"
	"github.com/vpesotskii/go-shortener-url/internal/app/models"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
)

var mapURLs map[string]string

// func encodes the URL from the request and put it into the map
func addURL(res http.ResponseWriter, req *http.Request) {

	body, _ := io.ReadAll(req.Body)
	if string(body) == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("No URL in request"))
	}
	shortURL := encodeURL(body)
	mapURLs[shortURL] = string(body)
	saveToFile(string(body), shortURL)
	logger.Log.Info("Body add", zap.String("body", string(body)))
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(config.Options.BaseAddress + "/" + shortURL))
}

// func encodes the URL from the request with json
func addURLFromJSON(res http.ResponseWriter, req *http.Request) {

	var r models.Request
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&r); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	shortURL := base64.StdEncoding.EncodeToString([]byte(r.URL))
	logger.Log.Info("Body URL", zap.String("body", r.URL))
	mapURLs[shortURL] = r.URL
	saveToFile(r.URL, shortURL)
	resp := models.Response{
		Result: config.Options.BaseAddress + "/" + shortURL,
	}

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(res)
	if err := encoder.Encode(resp); err != nil {
		logger.Log.Debug("cannot encode response", zap.Error(err))
		return
	}
}

// func returns the original URL by short URL
func getURL(res http.ResponseWriter, req *http.Request) {
	surl := chi.URLParam(req, "surl")
	if surl != "" {
		if originalURL, ok := mapURLs[chi.URLParam(req, "surl")]; ok {
			res.Header().Set("Location", originalURL)
			res.WriteHeader(http.StatusTemporaryRedirect)
		}
	}
	res.Header().Set("Location", "URL not found")
	res.WriteHeader(http.StatusBadRequest)
}

// func encodes string by base64
func encodeURL(url []byte) string {
	return base64.StdEncoding.EncodeToString(url)
}

func saveToFile(originalURL string, shortURL string) {
	file, err := os.OpenFile(config.Options.FileStorage, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		logger.Log.Debug("cannot open file", zap.Error(err))
		return
	}
	defer file.Close()

	rowNumber := 1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rowNumber++
	}
	if err := scanner.Err(); err != nil {
		logger.Log.Debug("Error Reading", zap.Error(err))
	}

	fileRecord := models.FileRecord{
		UUID:        rowNumber,
		OriginalURL: originalURL,
		ShortURL:    shortURL,
	}
	data, err := json.Marshal(fileRecord)
	if err != nil {
		return
	}

	_, err = file.Write(data)
	if err != nil {
		logger.Log.Debug("cannot write into file", zap.Error(err))
		return
	}
	_, err = file.WriteString("\n")
	if err != nil {
		return
	}
	logger.Log.Info("Data successfully appended to ", zap.String("file", config.Options.FileStorage))
}

func main() {
	mapURLs = make(map[string]string)
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/{surl}", logger.WithLogger(compress.GzipMiddleware(getURL)))
		r.Post("/", logger.WithLogger(compress.GzipMiddleware(addURL)))
		r.Post("/api/shorten", logger.WithLogger(compress.GzipMiddleware(addURLFromJSON)))
	})
	config.ParseFlags()
	err := logger.Initialize(config.Options.LogLevel)
	if err != nil {
		return
	}
	logger.Log.Info("Running server on", zap.String("server", config.Options.Server))
	logger.Log.Info("Base address", zap.String("base address", config.Options.BaseAddress))
	logger.Log.Info("File Storage Path", zap.String("file", config.Options.FileStorage))
	err = http.ListenAndServe(config.Options.Server, r)
	if err != nil {
		return
	}
}
