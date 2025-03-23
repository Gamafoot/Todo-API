package models

type Project struct {
	Id     uint   `gorm:"primaryKey"`
	Name   string `gorm:"type:varchar(50);not null"`
	UserId uint
}

func (p Project) TableName() string {
	return "porjects"
}
