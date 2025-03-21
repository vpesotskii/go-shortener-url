package handlers

import (
	"encoding/base64"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/vpesotskii/go-shortener-url/cmd/config"
	"github.com/vpesotskii/go-shortener-url/internal/app/logger"
	"github.com/vpesotskii/go-shortener-url/internal/app/models"
	"github.com/vpesotskii/go-shortener-url/internal/app/storage"
	"go.uber.org/zap"
	"io"
	"net/http"
)

// func returns the original URL by short URL
func GetURL(db storage.Repository, res http.ResponseWriter, req *http.Request) {
	surl := chi.URLParam(req, "surl")
	if surl != "" {
		fromStorage, ok := db.GetByID(surl)
		if ok {
			res.Header().Set("Location", fromStorage.OriginalURL)
			res.WriteHeader(http.StatusTemporaryRedirect)
		}
	}
	res.Header().Set("Location", "URL not found")
	res.WriteHeader(http.StatusBadRequest)
}

// func encodes the URL from the request with json
func AddURLFromJSON(db storage.Repository, res http.ResponseWriter, req *http.Request) {

	var r models.Request
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&r); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	shortURL := base64.StdEncoding.EncodeToString([]byte(r.URL))
	logger.Log.Info("Body URL", zap.String("body", r.URL))
	url := models.NewURL(1, shortURL, r.URL)
	err := db.Create(url)
	if err != nil {
		logger.Log.Debug(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
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

// func encodes the URL from the request and put it into the map
func AddURL(db storage.Repository, res http.ResponseWriter, req *http.Request) {

	body, _ := io.ReadAll(req.Body)
	if string(body) == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("No URL in request"))
	}

	shortURL := base64.StdEncoding.EncodeToString(body)
	url := models.NewURL(1, shortURL, string(body))
	err := db.Create(url)
	if err != nil {
		logger.Log.Info(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(config.Options.BaseAddress + "/" + shortURL))
}

func Ping(db storage.Repository, res http.ResponseWriter, req *http.Request) {
	err := db.Ping()
	if err != nil {
		logger.Log.Info("Not Connected: ", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
	} else {
		res.WriteHeader(http.StatusOK)
	}
}

func Batch(db storage.Repository, res http.ResponseWriter, req *http.Request) {

	body, _ := io.ReadAll(req.Body)
	if string(body) == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Empty Batch"))
	}

	var r []models.BatchRequest
	err := json.Unmarshal(body, &r)
	if err != nil {
		logger.Log.Info("Error decoding JSON:", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	for i := range r {
		r[i].ShortURL = base64.StdEncoding.EncodeToString([]byte(r[i].OriginalURL))
	}

	err = db.InsertBatch(r)
	if err != nil {
		logger.Log.Debug("Cannot Insert Batch in Storage", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var modifiedResponse []map[string]interface{}
	for _, r := range r {
		modifiedResponse = append(modifiedResponse, map[string]interface{}{
			"correlation_id": r.CorrelationID,
			"original_url":   r.OriginalURL,
		})
	}

	// Send response with updated requests
	response, err := json.Marshal(modifiedResponse)
	if err != nil {
		logger.Log.Debug("Error encoding response:", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(response)
}
