package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/JamesDeGreese/ya_golang/internal/app"
	"github.com/JamesDeGreese/ya_golang/internal/app/storage"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type Handler struct {
	Config  app.Config
	Storage storage.Repository
}

func (h Handler) GetHandler(c *gin.Context) {
	ID := c.Param("ID")

	if ID == "" {
		c.String(http.StatusBadRequest, "")
		return
	}

	fullURL, err := h.Storage.Get(ID)
	if err != nil {
		c.String(http.StatusNotFound, "")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, fullURL)
}

func (h Handler) PostHandler(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	urlID := uuid.NewV4().String()
	err = h.Storage.Add(urlID, string(body))
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	short := fmt.Sprintf("%s:%d/%s", h.Config.Host, h.Config.Port, urlID)

	c.String(http.StatusCreated, "%s", short)
}

func (h Handler) PostHandlerJson(c *gin.Context) {
	var req PostJsonRequest

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	urlID := uuid.NewV4().String()
	err = h.Storage.Add(urlID, req.Url)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	res := PostJsonResponse{Result: fmt.Sprintf("%s:%d/%s", h.Config.Host, h.Config.Port, urlID)}

	c.JSON(http.StatusCreated, res)
}
