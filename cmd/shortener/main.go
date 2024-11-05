package main

import (
	"encoding/base64"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vpesotskii/go-shortener-url/cmd/config"
	"io"
	"net/http"
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
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte("http://localhost:8080/" + shortURL))
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

func main() {
	mapURLs = make(map[string]string)
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/{surl}", getURL)
		r.Post("/", addURL)
	})
	config.ParseFlags()
	fmt.Println("Running server on", config.Options.Server)
	fmt.Println("Base Address server: ", config.Options.BaseAddress)
	err := http.ListenAndServe(config.Options.Server, r)
	if err != nil {
		return
	}
}
