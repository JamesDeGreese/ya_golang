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

	fullURL, err := h.Storage.GetURL(ID)
	if fullURL == "" || err != nil {
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

	userLinks := h.Storage.GetUserURLs(userIDDec)
	if len(userLinks) == 0 {
		c.JSON(http.StatusNoContent, "{}")
		return
	}

	res := make([]UserLinkItem, 0)
	for _, ul := range userLinks {
		res = append(res, UserLinkItem{
			fmt.Sprintf("%s/%s", h.Config.BaseURL, ul.ID),
			ul.OriginalURL,
		})
	}

	c.JSON(http.StatusOK, res)
}

func (h Handler) DBPingHandler(c *gin.Context) {
	_, err := h.Storage.GetURL("fake_id")
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.String(http.StatusOK, "")
}

func (h Handler) ShortenBatchHandler(c *gin.Context) {
	var req ShortenBatchRequest
	userID := c.GetString("user-id")

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	links := make([]storage.ShortLink, 0)

	for _, link := range req {
		links = append(links, storage.ShortLink{ID: link.ID, OriginalURL: link.URL, UserID: userID})
	}

	err = h.Storage.AddURLBatch(links)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	res := make([]BatchLinkItem, 0)

	for _, link := range links {
		res = append(res, BatchLinkItem{ID: link.ID, SortURL: fmt.Sprintf("%s/%s", h.Config.BaseURL, link.ID)})
	}

	c.JSON(http.StatusCreated, res)
}

func storeNewLink(h Handler, c *gin.Context, URL string) (string, error) {
	urlID := uuid.NewV4().String()
	userID := c.GetString("user-id")

	if userID == "" {
		err := h.Storage.AddURL(storage.ShortLink{ID: urlID, OriginalURL: URL, UserID: "no-user"})
		if err != nil {
			return "", err
		}
	} else {
		err := h.Storage.AddURL(storage.ShortLink{ID: urlID, OriginalURL: URL, UserID: userID})
		if err != nil {
			return "", err
		}
	}

	return urlID, nil
}
