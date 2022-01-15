package data

import (
	"convention.ninja/internal/data"
	userData "convention.ninja/internal/users/data"
)

type Organization struct {
	data.SnowflakeModel
	Name           string
	NormalizedName string
	OwnerId        int64
	Owner          userData.User
}

func init() {
	_ = data.GetConn().AutoMigrate(&Organization{})
}
