package main

import (
	"github.com/satori/go.uuid"
	"io"
	"net/http"
	"strings"
)

type urlToShort struct {
	URL string `json:"url"`
}

// TODO: Переделать на другое хранилище
var localStore = make(map[string]string)

// TODO: Вынести реализацию из main.go в структуру проекта

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		PostHandler(w, r)
	} else if r.Method == http.MethodGet {
		GetHandler(w, r)
	} else {
		makeResponse(w, "", http.StatusMethodNotAllowed)
	}
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")

	if path == "" {
		makeResponse(w, "", http.StatusBadRequest)
		return
	}

	fullURL := localStore[path]
	if fullURL == "" {
		makeResponse(w, "", http.StatusNotFound)
		return
	}

	makeResponse(w, fullURL, http.StatusOK)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "text/plain" {
		makeResponse(w, "", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		makeResponse(w, "", http.StatusUnsupportedMediaType)
		return
	}

	urlID := uuid.NewV4().String()
	localStore[urlID] = string(body)

	makeResponse(w, urlID, http.StatusOK)
}

func makeResponse(w http.ResponseWriter, response string, httpStatusCode int) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(httpStatusCode)
	w.Write([]byte(response))
}

func main() {
	// маршрутизация запросов обработчику
	// TODO: Реализовать роутинг с вайлдкардами
	http.HandleFunc("/", RequestHandler)
	// запуск сервера с адресом localhost, порт 8080
	http.ListenAndServe(":8080", nil)
}
