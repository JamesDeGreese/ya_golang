package handlers

import (
	"github.com/JamesDeGreese/ya_golang/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetHandler(t *testing.T) {
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
			request := httptest.NewRequest(http.MethodGet, testCase.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(GetHandler)
			h.ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, testCase.want, result.StatusCode)
		})
	}
}

func TestPostHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://youtube.com"))
	w := httptest.NewRecorder()
	h := http.HandlerFunc(PostHandler)
	h.ServeHTTP(w, request)
	result := w.Result()
	defer result.Body.Close()

	assert.Equal(t, http.StatusCreated, result.StatusCode)
}

func TestRequestHandler(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		request string
		want    int
	}{
		{
			name:    "Test GET",
			method:  http.MethodGet,
			request: "/foobarbaz",
			want:    http.StatusNotFound,
		},
		{
			name:    "Test POST",
			method:  http.MethodPost,
			request: "https://example.org",
			want:    http.StatusCreated,
		},
		{
			name:    "Test PATCH",
			method:  http.MethodPatch,
			request: "https://example.org",
			want:    http.StatusMethodNotAllowed,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			request := httptest.NewRequest(testCase.method, testCase.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(RequestHandler)
			h.ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, testCase.want, result.StatusCode)
		})
	}
}
