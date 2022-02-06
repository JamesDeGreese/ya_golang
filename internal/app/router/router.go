package router

import (
	"github.com/JamesDeGreese/ya_golang/internal/app/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/:ID", handlers.GetHandler)
	r.POST("/", handlers.PostHandler)
	return r
}
