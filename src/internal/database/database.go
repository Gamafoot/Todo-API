package database

import (
	"log"
	"root/internal/database/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewConnect(dbURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateAllTables(db *gorm.DB) error {
	tables := []interface{}{
		models.User{},
		models.Session{},
		models.Project{},
		models.Column{},
		models.Task{},
		models.Subtask{},
	}

	if err := db.AutoMigrate(tables...); err != nil {
		log.Fatalf("failed to migrate models: %v", err)
	}

	return nil
}
