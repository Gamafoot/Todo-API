package models

type Project struct {
	Id     uint `gorm:"primaryKey"`
	UserId uint
	Name   string `gorm:"type:varchar(50);not null"`
}

func (p Project) TableName() string {
	return "projects"
}
