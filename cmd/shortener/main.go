package main

import (
	"os"
	"os/signal"
	"syscall"

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
	s := storage.ConstructStorage(c)
	r := router.SetupRouter(c, s)

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-ch
		storage.DestructStorage(c, s)
		os.Exit(0)
	}()

	err = r.Run(c.Address)
	if err != nil {
		return
	}
}
