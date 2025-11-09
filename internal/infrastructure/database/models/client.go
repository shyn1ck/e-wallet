package models

import "time"

// APIClient represents the database model for API clients
type APIClient struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	UserID    string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	SecretKey string    `gorm:"type:varchar(255);not null"`
	IsActive  bool      `gorm:"not null;default:true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (APIClient) TableName() string {
	return "api_clients"
}
