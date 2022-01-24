package data

import (
	"database/sql"
	"time"
)

type SnowflakeModel struct {
	ID        int64        `json:"id,string"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	DeletedAt sql.NullTime `json:"deletedAt"`
}
