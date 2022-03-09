package handlers

import (
	"context"
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

	fullURL, err := h.Storage.GetURL(ID)
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

	urlID, err := storeNewLink(h, c, string(body))
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

	urlID, err := storeNewLink(h, c, req.URL)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	res := PostJSONResponse{Result: fmt.Sprintf("%s/%s", h.Config.BaseURL, urlID)}

	c.JSON(http.StatusCreated, res)
}

func (h Handler) UserURLsHandler(c *gin.Context) {
	userIDEnc, err := c.Cookie("user-id")
	if err != nil {
		c.JSON(http.StatusNoContent, "{}")
		return
	}
	userIDDec, err := app.Decrypt([]byte(userIDEnc), h.Config.AppKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	userURLs := h.Storage.GetUserURLs(userIDDec)
	if len(userURLs) == 0 {
		c.JSON(http.StatusNoContent, "{}")
		return
	}

	var res []UserLinkItem

	for _, shortID := range userURLs {
		URL, _ := h.Storage.GetURL(shortID)
		res = append(res, UserLinkItem{
			fmt.Sprintf("%s/%s", h.Config.BaseURL, shortID),
			URL,
		})
	}

	c.JSON(http.StatusOK, res)
}

func (h Handler) DBPingHandler(c *gin.Context) {
	var res string
	err := h.Storage.DB().QueryRow(context.Background(), "select 'Hello, world!'").Scan(&res)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.String(http.StatusOK, "")
}

func storeNewLink(h Handler, c *gin.Context, URL string) (string, error) {
	userIDEnc, _ := c.Cookie("user-id")
	urlID := uuid.NewV4().String()

	if userIDEnc == "" {
		err := h.Storage.AddURL(urlID, URL, "no-user")
		if err != nil {
			return "", err
		}
	} else {
		userIDDec, err := app.Decrypt([]byte(userIDEnc), h.Config.AppKey)
		if err != nil {
			return "", err
		}
		err = h.Storage.AddURL(urlID, URL, userIDDec)
		if err != nil {
			return "", err
		}
	}

	return urlID, nil
}
