package main

import (
	"log"
	"root/internal/app"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
