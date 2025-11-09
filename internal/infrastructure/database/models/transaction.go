package models

import "time"

// Transaction represents the database model for transactions
type Transaction struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	WalletID  int64     `gorm:"index;not null"`
	Type      string    `gorm:"type:varchar(20);not null"` // deposit, withdrawal, etc.
	Amount    int64     `gorm:"not null"`                  // stored in minor units (dirams)
	CreatedAt time.Time `gorm:"autoCreateTime;index"`
}

// TableName specifies the table name for GORM
func (Transaction) TableName() string {
	return "transactions"
}
