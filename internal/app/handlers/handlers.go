package handlers

import (
	"github.com/JamesDeGreese/ya_golang/internal/app/storage"
	uuid "github.com/satori/go.uuid"
	"io"
	"net/http"
	"strings"
)

// TODO: Вынести в переменные окружения или конфиг
const hostName = "http://localhost:8080"

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetHandler(w, r)
		break
	case http.MethodPost:
		PostHandler(w, r)
		break
	default:
		makeResponse(w, "", http.StatusMethodNotAllowed)
	}
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")

	if path == "" {
		makeResponse(w, "", http.StatusBadRequest)
		return
	}

	s := storage.GetInstance()

	fullURL := s.Get(path)
	if fullURL == "" {
		makeResponse(w, "", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		makeResponse(w, "", http.StatusInternalServerError)
		return
	}

	s := storage.GetInstance()

	urlID := uuid.NewV4().String()
	s.Add(urlID, string(body))

	short := hostName + "/" + urlID

	makeResponse(w, short, http.StatusCreated)
}

func makeResponse(w http.ResponseWriter, response string, httpStatusCode int) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(httpStatusCode)
	w.Write([]byte(response))
}
