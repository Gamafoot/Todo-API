package main

import (
	"fmt"
	"log"
	"os"
	"root/internal/config"
	"root/internal/database"
	"strings"

	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

func main() {
	log.Println("Start migrating...")

	cfg := config.GetConfig()

	db, err := database.NewConnect(cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %+v", err)
	}

	err = database.CreateAllTables(db)
	if err != nil {
		log.Fatalf("Failed to create all tables: %+v\n", err)
	}

	scripts, err := getSqlScripts("functions")
	if err != nil {
		log.Fatalf("Failed get function scripts: %+v\n", err)
	}

	err = executeScripts(db, scripts)
	if err != nil {
		log.Fatalf("Failed to execute functions: %+v\n", err)
	}

	scripts, err = getSqlScripts("triggers")
	if err != nil {
		log.Fatalf("Failed get trigger scripts: %+v\n", err)
	}

	err = executeScripts(db, scripts)
	if err != nil {
		log.Fatalf("Failed to execute triggers: %+v\n", err)
	}

	log.Println("Migrate was successed")
}

func getSqlScripts(folderName string) (map[string]string, error) {
	rootDir := fmt.Sprintf("assets/migration/%s/", folderName)

	scripts := make(map[string]string)

	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	for _, entries := range entries {
		if !strings.HasSuffix(entries.Name(), ".sql") {
			continue
		}

		filepath := rootDir + entries.Name()

		content, err := os.ReadFile(filepath)
		if err != nil {
			return nil, pkgErrors.WithStack(err)
		}

		scripts[entries.Name()] = string(content)
	}

	return scripts, nil
}

func executeScripts(db *gorm.DB, scripts map[string]string) error {
	tx := db.Begin()

	for filename, script := range scripts {
		if err := tx.Exec(script).Error; err != nil {
			tx.Rollback()
			return pkgErrors.Errorf("%s: %v", filename, err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return pkgErrors.WithStack(err)
	}

	return nil
}
