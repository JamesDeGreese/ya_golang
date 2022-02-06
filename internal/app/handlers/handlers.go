package handlers

import (
	"github.com/JamesDeGreese/ya_golang/internal/app/storage"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
)

// TODO: Вынести в переменные окружения или конфиг
const hostName = "http://localhost:8080"

func GetHandler(c *gin.Context) {
	ID := c.Param("ID")

	if ID == "" {
		MakeResponse(c, "", http.StatusBadRequest)
		return
	}

	s := storage.GetInstance()

	fullURL := s.Get(ID)
	if fullURL == "" {
		MakeResponse(c, "", http.StatusNotFound)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, fullURL)
}

func PostHandler(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		MakeResponse(c, "", http.StatusInternalServerError)
		return
	}

	s := storage.GetInstance()

	urlID := uuid.NewV4().String()
	s.Add(urlID, string(body))

	short := hostName + "/" + urlID

	MakeResponse(c, short, http.StatusCreated)
}

func MakeResponse(c *gin.Context, response string, httpStatusCode int) {
	c.String(httpStatusCode, "%s", response)
}
