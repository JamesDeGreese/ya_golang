package main

import (
	"github.com/JamesDeGreese/ya_golang/internal/app/handlers"
	"net/http"
)

func main() {
	// маршрутизация запросов обработчику
	// TODO: Реализовать роутинг с вайлдкардами
	http.HandleFunc("/", handlers.RequestHandler)
	// запуск сервера с адресом localhost, порт 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
