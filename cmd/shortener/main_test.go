package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_addURL(t *testing.T) {
	mapURLs = make(map[string]string)
	body := "https://practicum.yandex.ru/"

	tests := []struct {
		name         string
		method       string
		expectedCode int
		expectedBody string
	}{
		{name: "Post with body", method: http.MethodPost, expectedCode: http.StatusCreated, expectedBody: body},
		{name: "Post without body", method: http.MethodPost, expectedCode: http.StatusBadRequest, expectedBody: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.expectedBody))
			response := httptest.NewRecorder()
			addURL(response, request)
			assert.Equal(t, tt.expectedCode, response.Code)
		})
	}
}

func Test_encodeURL(t *testing.T) {
	if shortURL := encodeURL([]byte("https://practicum.yandex.ru/")); shortURL != "aHR0cHM6Ly9wcmFjdGljdW0ueWFuZGV4LnJ1Lw==" {
		t.Errorf("ShortURL is not correct; got %s", shortURL)
	}
}

func Test_getURL(t *testing.T) {

	tests := []struct {
		name         string
		expectedCode int
		URL          string
	}{
		{name: "GET with existing URL", expectedCode: http.StatusTemporaryRedirect, URL: "https://practicum.yandex.ru/"},
		{name: "Wrong GET", expectedCode: http.StatusBadRequest, URL: "/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var addr string
			if tt.URL == "/" {
				addr = ""
			} else {
				for k, v := range mapURLs {
					if v == tt.URL {
						addr = k
					}
				}
			}
			request := httptest.NewRequest(http.MethodGet, "/"+addr, nil)
			response := httptest.NewRecorder()
			getURL(response, request)
			assert.Equal(t, tt.expectedCode, response.Code)
		})
	}
}
