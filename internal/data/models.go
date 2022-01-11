package data

import (
	"gorm.io/gorm"
	"time"
)

type SnowflakeModel struct {
	ID        int64 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
