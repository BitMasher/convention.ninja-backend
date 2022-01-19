package data

import (
	"gorm.io/gorm"
	"time"
)

type SnowflakeModel struct {
	ID        int64          `gorm:"primaryKey" json:"id,string"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}
