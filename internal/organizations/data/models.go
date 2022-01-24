package data

import (
	"convention.ninja/internal/data"
	userData "convention.ninja/internal/users/data"
)

type Organization struct {
	data.SnowflakeModel
	Name           string        `json:"name"`
	NormalizedName string        `json:"-"`
	OwnerId        int64         `json:"ownerId,string"`
	Owner          userData.User `json:"-"`
}
