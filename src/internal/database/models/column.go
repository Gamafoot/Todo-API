package models

type Column struct {
	Id        uint `gorm:"primaryKey"`
	ProjectId uint
	Title     string `gorm:"type:varchar(50);not null"`
}

func (c Column) TableName() string {
	return "columns"
}
