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
	Subtask SubtaskStorage
}

func NewPostgresStorage(db *gorm.DB) *Storage {
	return &Storage{
		User:    postgres.NewUserStorage(db),
		Session: postgres.NewSessionStorage(db),
		Project: postgres.NewProjectStorage(db),
		Column: postgres.NewColumnStorage(db),
		Task:    postgres.NewTaskStorage(db),
		Subtask: postgres.NewSubtaskStorage(db),
	}
}
