package main

import "github.com/JamesDeGreese/ya_golang/internal/app/router"

func main() {
	r := router.SetupRouter()
	r.Run(":8080")
}
