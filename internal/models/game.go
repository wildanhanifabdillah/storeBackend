package models

import "time"

type Game struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:100;not null"`
	Code      string `gorm:"size:50;unique;not null"`
	IsActive  bool   `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Packages []TopupPackage `gorm:"foreignKey:GameID"`
}
