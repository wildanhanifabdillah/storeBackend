package models

import "time"

type Admin struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"size:100;not null"`
	Email        string `gorm:"size:100;unique;not null"`
	PasswordHash string `gorm:"type:text;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
