package models

import "time"

type Transaction struct {
	ID          uint   `gorm:"primaryKey"`
	OrderID     string `gorm:"size:50;uniqueIndex;not null"`
	GameID      uint   `gorm:"not null"`
	PackageID   uint   `gorm:"not null"`
	GameUserID  string `gorm:"size:100;not null"`
	Email       string `gorm:"size:100;not null"`
	TotalAmount int64  `gorm:"not null"`
	Status      string `gorm:"size:20;not null"` // PENDING, PAID, FAILED
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Game    Game         `gorm:"foreignKey:GameID"`
	Package TopupPackage `gorm:"foreignKey:PackageID"`
	Payment Payment
}
