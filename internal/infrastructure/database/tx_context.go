package database

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

// InjectTx injects a transaction into the context
func InjectTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// ExtractTx extracts a transaction from the context
func ExtractTx(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return nil
}

// GetDB returns the transaction from context if exists, otherwise returns the provided db
func GetDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx := ExtractTx(ctx); tx != nil {
		return tx
	}
	return db
}
