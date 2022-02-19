package main

import (
	"github.com/JamesDeGreese/ya_golang/internal/app"
	"github.com/JamesDeGreese/ya_golang/internal/app/router"
	"github.com/JamesDeGreese/ya_golang/internal/app/storage"
	"github.com/caarlos0/env/v6"
)

func main() {
	c := app.Config{}
	err := env.Parse(&c)
	if err != nil {
		return
	}
	s := storage.ConstructStorage()
	r := router.SetupRouter(c, s)
	err = r.Run(c.Address)
	if err != nil {
		return
	}
}
