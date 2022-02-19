package main

import (
	"fmt"

	"github.com/JamesDeGreese/ya_golang/internal/app"
	"github.com/JamesDeGreese/ya_golang/internal/app/router"
	"github.com/JamesDeGreese/ya_golang/internal/app/storage"
)

func main() {
	c := app.Config{
		Host: "http://localhost",
		Port: 8080,
	}
	s := storage.ConstructStorage()
	r := router.SetupRouter(c, s)
	r.Run(fmt.Sprintf(":%d", c.Port))
}
