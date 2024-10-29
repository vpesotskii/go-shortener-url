package main

import (
	"encoding/base64"
	"io"
	"net/http"
)

var mapUrls = map[string]string{}

func handleMethod(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addUrl(w, r)
	case http.MethodGet:
		getUrl(w, r)
	default:
		http.Error(w, "Method is allowed", http.StatusBadRequest)
	}

}

func addUrl(res http.ResponseWriter, req *http.Request) {

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	shortUrl := encodeUrl(body)
	mapUrls[shortUrl] = string(body)
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortUrl))
}

func getUrl(res http.ResponseWriter, req *http.Request) {
	shortUrl := req.URL.String()[1:]
	if s, ok := mapUrls[shortUrl]; ok {
		res.Header().Set("Location", s)
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func encodeUrl(url []byte) string {
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
