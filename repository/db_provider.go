package repository

import (
	"context"

	"gorm.io/gorm"
)

type DBProvider interface {
	WithContext(ctx context.Context) *gorm.DB
}
