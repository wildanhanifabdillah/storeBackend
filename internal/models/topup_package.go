package models

import "time"

type TopupPackage struct {
	ID        uint   `gorm:"primaryKey"`
	GameID    uint   `gorm:"not null"`
	Name      string `gorm:"size:100;not null"`
	Amount    int    `gorm:"not null"`
	Price     int64  `gorm:"not null"`
	IsActive  bool   `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Game Game `gorm:"foreignKey:GameID"`
}
