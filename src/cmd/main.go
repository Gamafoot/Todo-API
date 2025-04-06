package main

import (
	"log"
	_ "root/docs"
	"root/internal/app"
)

// @title Todo API
// @version 1.0
// @description API для управления задачами
// @host localhost:8000
// @BasePath /api/v1
func main() {
	err := app.Run()
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
