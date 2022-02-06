package main

import (
	"github.com/JamesDeGreese/ya_golang/internal/app/router"
	"github.com/JamesDeGreese/ya_golang/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetShortLink(t *testing.T) {
	s := storage.GetInstance()
	s.Add("123", "https://example.org")

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
			r := router.SetupRouter()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, testCase.request, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.want, w.Code)
		})
	}
}

func TestCreateShortLink(t *testing.T) {
	r := router.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader("https://youtube.com"))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
