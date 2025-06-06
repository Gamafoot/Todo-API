package model

type Column struct {
	Id        uint   `gorm:"primaryKey"`
	ProjectId uint   `gorm:"not null"`
	Name      string `gorm:"type:varchar(50);not null"`
	Position  int
	Project   Project `gorm:"foreignKey:ProjectId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (c Column) TableName() string {
	return "columns"
}
