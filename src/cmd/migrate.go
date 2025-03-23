package main

import (
	"log"
	"root/internal/config"
	"root/internal/database"
)

func main() {
	log.Println("start migrating...")

	cfg := config.GetConfig()

	db, err := database.NewConnect(cfg.Database.URL)
	if err != nil {
		log.Fatalf("failed to connect to database: %+v", err)
	}

	err = database.CreateAllTables(db)
	if err != nil {
		log.Fatalf("failed to create all tables: %+v\n", err)
	}

	log.Println("migrate was successed")
}
