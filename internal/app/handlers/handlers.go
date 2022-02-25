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

	short := fmt.Sprintf("%s/%s", h.Config.BaseURL, urlID)

	c.String(http.StatusCreated, "%s", short)
}

func (h Handler) PostHandlerJSON(c *gin.Context) {
	var req PostJSONRequest

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	urlID := uuid.NewV4().String()
	err = h.Storage.Add(urlID, req.URL)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	res := PostJSONResponse{Result: fmt.Sprintf("%s/%s", h.Config.BaseURL, urlID)}

	c.JSON(http.StatusCreated, res)
}
