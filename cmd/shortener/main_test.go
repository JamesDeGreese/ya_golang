package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JamesDeGreese/ya_golang/internal/app"
	"github.com/JamesDeGreese/ya_golang/internal/app/handlers"
	"github.com/JamesDeGreese/ya_golang/internal/app/router"
	"github.com/JamesDeGreese/ya_golang/internal/app/storage"
	"github.com/caarlos0/env/v6"
	"github.com/stretchr/testify/assert"
)

func TestGetShortLink(t *testing.T) {
	c := app.Config{}
	err := env.Parse(&c)
	if err != nil {
		t.FailNow()
	}
	s := storage.ConstructStorage(c)
	err = s.Add("123", "https://example.org")
	if err != nil {
		t.FailNow()
	}

	tests := []struct {
		name    string
		method  string
		request string
		want    int
	}{
		{
			name:    "Test 200",
			request: "/123",
			want:    http.StatusTemporaryRedirect,
		},
		{
			name:    "Test 404",
			request: "/1234",
			want:    http.StatusNotFound,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := router.SetupRouter(c, s)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, testCase.request, nil)
			r.ServeHTTP(w, req)

			assert.NoError(t, err)
			assert.Equal(t, testCase.want, w.Code)
		})
	}
}

func TestCreateShortLink(t *testing.T) {
	c := app.Config{}
	err := env.Parse(&c)
	if err != nil {
		t.FailNow()
	}
	s := storage.ConstructStorage(c)
	r := router.SetupRouter(c, s)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader("https://youtube.com"))
	r.ServeHTTP(w, req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateShortLinkJSON(t *testing.T) {
	c := app.Config{}
	err := env.Parse(&c)
	if err != nil {
		t.FailNow()
	}
	s := storage.ConstructStorage(c)
	r := router.SetupRouter(c, s)

	w := httptest.NewRecorder()
	rBody, _ := json.Marshal(handlers.PostJSONRequest{URL: "https://youtube.com"})
	b := bytes.NewBuffer(rBody)
	req, err := http.NewRequest(http.MethodPost, "/api/shorten", b)
	r.ServeHTTP(w, req)

	res := w.Body.String()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotEmpty(t, res)
}

func TestGetShortLinkGzip(t *testing.T) {
	c := app.Config{}
	err := env.Parse(&c)
	if err != nil {
		t.FailNow()
	}
	s := storage.ConstructStorage(c)
	err = s.Add("123", "https://example.org")
	if err != nil {
		t.FailNow()
	}

	tests := []struct {
		name    string
		method  string
		request string
		want    int
	}{
		{
			name:    "Test 200",
			request: "/123",
			want:    http.StatusTemporaryRedirect,
		},
		{
			name:    "Test 404",
			request: "/1234",
			want:    http.StatusNotFound,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := router.SetupRouter(c, s)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, testCase.request, nil)
			req.Header.Set("Accept-Encoding", "gzip")
			r.ServeHTTP(w, req)

			assert.NoError(t, err)
			assert.Equal(t, testCase.want, w.Code)
		})
	}
}

func TestCreateShortLinkGzip(t *testing.T) {
	c := app.Config{}
	err := env.Parse(&c)
	if err != nil {
		return
	}
	s := storage.ConstructStorage(c)
	r := router.SetupRouter(c, s)

	var b bytes.Buffer
	gzw := gzip.NewWriter(&b)
	gzw.Write([]byte("https://youtube.com"))
	gzw.Close()

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(b.Bytes()))
	req.Header.Set("Content-Encoding", "gzip")
	r.ServeHTTP(w, req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotEmpty(t, w.Body.String())
}

func TestCreateShortLinkJSONGzip(t *testing.T) {
	c := app.Config{}
	err := env.Parse(&c)
	if err != nil {
		t.FailNow()
	}
	s := storage.ConstructStorage(c)
	r := router.SetupRouter(c, s)

	w := httptest.NewRecorder()
	rBody, _ := json.Marshal(handlers.PostJSONRequest{URL: "https://youtube.com"})
	var b bytes.Buffer
	gzw := gzip.NewWriter(&b)
	gzw.Write(rBody)
	gzw.Close()
	req, err := http.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(b.Bytes()))
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")
	r.ServeHTTP(w, req)

	res := w.Body.String()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotEmpty(t, res)
}
