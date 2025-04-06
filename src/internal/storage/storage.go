package storage

import (
	"root/internal/storage/postgres"

	"gorm.io/gorm"
)

type Storage struct {
	User    UserStorage
	Session SessionStorage
	Project ProjectStorage
	Column ColumnStorage
	Task    TaskStorage
}

func NewPostgresStorage(db *gorm.DB) *Storage {
	return &Storage{
		User:    postgres.NewUserStorage(db),
		Session: postgres.NewSessionStorage(db),
		Task:    postgres.NewTaskStorage(db),
		Project: postgres.NewProjectStorage(db),
		Column: postgres.NewColumnStorage(db),
	}
}
