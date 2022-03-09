package main

import (
	"flag"
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

	flag.StringVar(&c.Address, "a", c.Address, "a 127.0.0.1:8080")
	flag.StringVar(&c.BaseURL, "b", c.BaseURL, "b https://example.org")
	flag.StringVar(&c.FileStoragePath, "f", c.FileStoragePath, "f /tmp/storage")
	flag.StringVar(&c.FileStoragePath, "d", c.DatabaseDSN, "f postgres://username:password@host:port/database_name")
	flag.Parse()

	s := storage.InitStorage(c)
	r := router.SetupRouter(c, s)

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-ch
		storage.CleanupStorage(c, s)
		os.Exit(0)
	}()

	err = r.Run(c.Address)
	if err != nil {
		return
	}
}
