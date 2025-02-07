package main

//
//import (
//	"bytes"
//	"context"
//	"encoding/json"
//	"github.com/go-chi/chi/v5"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"github.com/vpesotskii/go-shortener-url/cmd/config"
//	"github.com/vpesotskii/go-shortener-url/internal/app/models"
//	"net/http"
//	"net/http/httptest"
//	"strings"
//	"testing"
//)
//
//func Test_addURL(t *testing.T) {
//	mapURLs = make(map[string]string)
//	body := "https://practicum.yandex.ru/"
//
//	tests := []struct {
//		name         string
//		method       string
//		expectedCode int
//		expectedBody string
//	}{
//		{name: "Post with body", method: http.MethodPost, expectedCode: http.StatusCreated, expectedBody: body},
//		{name: "Post without body", method: http.MethodPost, expectedCode: http.StatusBadRequest, expectedBody: ""},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			request := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.expectedBody))
//			response := httptest.NewRecorder()
//			addURL(response, request)
//			assert.Equal(t, tt.expectedCode, response.Code)
//		})
//	}
//}
//
//func Test_encodeURL(t *testing.T) {
//	if shortURL := encodeURL([]byte("https://practicum.yandex.ru/")); shortURL != "aHR0cHM6Ly9wcmFjdGljdW0ueWFuZGV4LnJ1Lw==" {
//		t.Errorf("ShortURL is not correct; got %s", shortURL)
//	}
//}
//
//func Test_getURL(t *testing.T) {
//
//	tests := []struct {
//		name         string
//		expectedCode int
//		URL          string
//	}{
//		{name: "GET with existing URL", expectedCode: http.StatusTemporaryRedirect, URL: "https://practicum.yandex.ru/"},
//		{name: "Wrong GET", expectedCode: http.StatusBadRequest, URL: "/"},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			var addr string
//			if tt.URL == "/" {
//				addr = ""
//			} else {
//				for k, v := range mapURLs {
//					if v == tt.URL {
//						addr = k
//					}
//				}
//			}
//
//			req, err := http.NewRequest(http.MethodGet, "/", nil)
//			rctx := chi.NewRouteContext()
//			rctx.URLParams.Add("surl", addr)
//			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
//			require.NoError(t, err)
//
//			response := httptest.NewRecorder()
//			getURL(response, req)
//
//			require.NoError(t, err)
//			assert.Equal(t, tt.expectedCode, response.Code)
//		})
//	}
//}
//
//func Test_addURLFromJSON(t *testing.T) {
//
//	req := models.Request{
//		URL: "https://practicum.yandex.ru/",
//	}
//
//	resp := models.Response{
//		Result: config.Options.BaseAddress + "/" + "aHR0cHM6Ly9wcmFjdGljdW0ueWFuZGV4LnJ1Lw==",
//	}
//
//	tests := []struct {
//		name           string
//		method         string
//		expectedCode   int
//		expectedBody   models.Request
//		expectedResult models.Response
//	}{
//		{name: "Post with body", method: http.MethodPost, expectedCode: http.StatusCreated, expectedBody: req, expectedResult: resp},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			jsonBody, _ := json.Marshal(tt.expectedBody)
//			request := httptest.NewRequest(tt.method, "/api/shorten", bytes.NewBuffer(jsonBody))
//			request.Header.Set("Content-Type", "application/json")
//			response := httptest.NewRecorder()
//			addURLFromJSON(response, request)
//			assert.Equal(t, tt.expectedCode, response.Code)
//
//			responseBody := response.Body.String()
//			var respJSON models.Response
//			err := json.Unmarshal([]byte(responseBody), &respJSON)
//			assert.NoError(t, err, "Response JSON should be valid")
//
//			assert.Equal(t, tt.expectedResult, resp, "Short URL does not match expected value")
//			assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"), "Content-Type header mismatch")
//		})
//	}
//}
