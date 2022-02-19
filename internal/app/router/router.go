package router

import (
	"github.com/JamesDeGreese/ya_golang/internal/app"
	"github.com/JamesDeGreese/ya_golang/internal/app/handlers"
	"github.com/JamesDeGreese/ya_golang/internal/app/storage"
	"github.com/gin-gonic/gin"
)

func SetupRouter(c app.Config, s *storage.Storage) *gin.Engine {
	r := gin.Default()
	h := handlers.Handler{
		Config:  c,
		Storage: s,
	}
	r.GET("/:ID", h.GetHandler)
	r.POST("/", h.PostHandler)
	r.POST("/api/shorten", h.PostHandlerJSON)
	return r
}
