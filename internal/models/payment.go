package models

import "time"

type Payment struct {
	ID             uint   `gorm:"primaryKey"`
	TransactionID  uint   `gorm:"uniqueIndex;not null"`
	PaymentGateway string `gorm:"size:50;not null"` // midtrans
	PaymentType    string `gorm:"size:50"`          // qris, bank_transfer
	PaymentToken   string `gorm:"type:text"`        // snap token
	PaymentStatus  string `gorm:"size:20"`          // pending, settlement
	RawResponse    string `gorm:"type:text"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
