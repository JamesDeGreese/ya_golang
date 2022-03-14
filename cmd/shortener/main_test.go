package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/JamesDeGreese/ya_golang/internal/app"
	"github.com/JamesDeGreese/ya_golang/internal/app/handlers"
	"github.com/JamesDeGreese/ya_golang/internal/app/router"
	"github.com/JamesDeGreese/ya_golang/internal/app/storage"
	"github.com/bxcodec/faker/v3"
	"github.com/caarlos0/env/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetShortLink(t *testing.T) {
	c := app.Config{}
	err := env.Parse(&c)
	if err != nil {
		t.FailNow()
	}
	s := storage.InitStorage(c)
	URL := faker.URL()
	err = s.AddURL(storage.ShortLink{ID: "123", OriginalURL: URL, UserID: "12345"})
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
	s := storage.InitStorage(c)
	r := router.SetupRouter(c, s)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(faker.URL()))
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
	s := storage.InitStorage(c)
	r := router.SetupRouter(c, s)

	w := httptest.NewRecorder()
	rBody, _ := json.Marshal(handlers.PostJSONRequest{URL: faker.URL()})
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
	s := storage.InitStorage(c)
	ID := faker.DomainName()
	err = s.AddURL(storage.ShortLink{ID: ID, OriginalURL: faker.URL(), UserID: "12345"})
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
			request: fmt.Sprintf("/%s", ID),
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
	s := storage.InitStorage(c)
	r := router.SetupRouter(c, s)

	var b bytes.Buffer
	gzw := gzip.NewWriter(&b)
	gzw.Write([]byte(faker.URL()))
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
	s := storage.InitStorage(c)
	r := router.SetupRouter(c, s)

	w := httptest.NewRecorder()
	rBody, _ := json.Marshal(handlers.PostJSONRequest{URL: faker.URL()})
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

func TestGetUserLinks(t *testing.T) {
	c := app.Config{}
	err := env.Parse(&c)
	if err != nil {
		t.FailNow()
	}
	s := storage.InitStorage(c)
	userID := uuid.NewV4().String()
	userIDEnc, err := app.Encrypt(userID, c.AppKey)
	if err != nil {
		t.FailNow()
	}
	err = s.AddURL(storage.ShortLink{ID: faker.DomainName(), OriginalURL: "https://example.org", UserID: userID})
	if err != nil {
		t.FailNow()
	}

	r := router.SetupRouter(c, s)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/api/user/urls", nil)
	req.AddCookie(&http.Cookie{
		Name:     "user-id",
		Value:    url.QueryEscape(userIDEnc),
		MaxAge:   3600,
		Path:     "/",
		Domain:   c.Address,
		Secure:   false,
		HttpOnly: false,
	})
	r.ServeHTTP(w, req)

	res := w.Body.String()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, res)
}

func TestGetUserLinksEmpty(t *testing.T) {
	c := app.Config{}
	err := env.Parse(&c)
	if err != nil {
		t.FailNow()
	}
	s := storage.InitStorage(c)

	r := router.SetupRouter(c, s)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/api/user/urls", nil)
	r.ServeHTTP(w, req)

	res := w.Body.String()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, res)
}

func TestPingDBFail(t *testing.T) {
	c := app.Config{}
	err := env.Parse(&c)
	if err != nil {
		t.FailNow()
	}
	s := storage.InitStorage(c)
	s.CleanUp(c)

	r := router.SetupRouter(c, s)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/ping", nil)
	r.ServeHTTP(w, req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestBatchInsert(t *testing.T) {
	c := app.Config{}
	err := env.Parse(&c)
	if err != nil {
		t.FailNow()
	}
	s := storage.InitStorage(c)
	r := router.SetupRouter(c, s)

	w := httptest.NewRecorder()
	rBody, _ := json.Marshal(handlers.ShortenBatchRequest{
		{ID: "123", URL: faker.URL()},
		{ID: "321", URL: faker.URL()},
	})
	b := bytes.NewBuffer(rBody)
	req, err := http.NewRequest(http.MethodPost, "/api/shorten/batch", b)
	r.ServeHTTP(w, req)

	res := w.Body.String()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotEmpty(t, res)
}

func TestCreateShortLinkDuplicate(t *testing.T) {
	c := app.Config{}
	err := env.Parse(&c)
	if err != nil {
		t.FailNow()
	}
	s := storage.InitStorage(c)
	r := router.SetupRouter(c, s)
	ID := faker.DomainName()
	URL := faker.URL()
	err = s.AddURL(storage.ShortLink{ID: ID, OriginalURL: URL, UserID: "12345"})
	if err != nil {
		t.FailNow()
	}
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(URL))
	r.ServeHTTP(w, req)

	res := w.Body.String()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, w.Code)
	assert.NotEmpty(t, res)
}

func TestCreateShortLinkJSONDuplicate(t *testing.T) {
	c := app.Config{}
	err := env.Parse(&c)
	if err != nil {
		t.FailNow()
	}
	s := storage.InitStorage(c)
	r := router.SetupRouter(c, s)
	ID := faker.DomainName()
	URL := faker.URL()
	err = s.AddURL(storage.ShortLink{ID: ID, OriginalURL: URL, UserID: "12345"})
	if err != nil {
		t.FailNow()
	}
	w := httptest.NewRecorder()
	rBody, _ := json.Marshal(handlers.PostJSONRequest{URL: URL})
	b := bytes.NewBuffer(rBody)
	req, err := http.NewRequest(http.MethodPost, "/api/shorten", b)

	r.ServeHTTP(w, req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, w.Code)
}
