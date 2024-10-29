package main

import (
	"encoding/base64"
	"io"
	"net/http"
)

var mapURLs = map[string]string{}

// func handles the request
func handleMethod(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addURL(w, r)
	case http.MethodGet:
		getUrl(w, r)
	default:
		http.Error(w, "Method is allowed", http.StatusBadRequest)
	}

}

// func encodes the URL from the request and put it into the map
func addURL(res http.ResponseWriter, req *http.Request) {

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	shortURL := encodeURL(body)
	mapURLs[shortURL] = string(body)
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte("http://localhost:8080/" + shortURL))
}

// func returns the original URL by short URL
func getUrl(res http.ResponseWriter, req *http.Request) {
	shortURL := req.URL.String()[1:]
	if originalURL, ok := mapURLs[shortURL]; ok {
		res.Header().Set("Location", originalURL)
		res.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		res.Header().Set("Location", "URL not found")
		res.WriteHeader(http.StatusBadRequest)
	}
}

// func encodes string by base64
func encodeURL(url []byte) string {
	return base64.StdEncoding.EncodeToString(url)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleMethod)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return
	}
}
