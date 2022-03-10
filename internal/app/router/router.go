package router

import (
	"github.com/JamesDeGreese/ya_golang/internal/app"
	"github.com/JamesDeGreese/ya_golang/internal/app/handlers"
	"github.com/JamesDeGreese/ya_golang/internal/app/storage"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func SetupRouter(c app.Config, s storage.Repository) *gin.Engine {
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.BestSpeed, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	r.Use(app.AuthCookieMiddleware(c))
	h := handlers.Handler{
		Config:  c,
		Storage: s,
	}
	r.GET("/:ID", h.GetHandler)
	r.POST("/", h.PostHandler)
	r.POST("/api/shorten", h.PostHandlerJSON)
	r.GET("/api/user/urls", h.UserURLsHandler)
	r.GET("/ping", h.DBPingHandler)
	return r
}
