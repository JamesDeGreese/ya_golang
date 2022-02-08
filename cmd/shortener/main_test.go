package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JamesDeGreese/ya_golang/internal/app"
	"github.com/JamesDeGreese/ya_golang/internal/app/router"
	"github.com/JamesDeGreese/ya_golang/internal/app/storage"
	"github.com/stretchr/testify/assert"
)

func TestGetShortLink(t *testing.T) {
	c := app.Config{
		Host: "http://localhost",
		Port: 8080,
	}
	s := storage.ConstructStorage()
	err := s.Add("123", "https://example.org")
	if err != nil {
		return
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

			assert.Equal(t, testCase.want, w.Code)
			assert.NoError(t, err)
		})
	}
}

func TestCreateShortLink(t *testing.T) {
	c := app.Config{
		Host: "http://localhost",
		Port: 8080,
	}
	s := storage.ConstructStorage()
	r := router.SetupRouter(c, s)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader("https://youtube.com"))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NoError(t, err)
}
