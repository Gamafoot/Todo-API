package model

import "time"

type Heatmap struct {
	Id     uint      `gorm:"primaryKey"`
	UserId uint      `gorm:"not null"`
	Date   time.Time `gorm:"type:date;not null;unique"`
	Count  int       `gorm:"default:1"`
	User   User      `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (h Heatmap) TableName() string {
	return "heatmaps"
}
