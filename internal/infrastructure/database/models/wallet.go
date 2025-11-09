package models

import "time"

// Wallet represents the database model for wallets
type Wallet struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	AccountID string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	Type      string    `gorm:"type:varchar(20);not null"` // identified or unidentified
	Balance   int64     `gorm:"not null;default:0"`        // stored in minor units (dirams)
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (Wallet) TableName() string {
	return "wallets"
}
