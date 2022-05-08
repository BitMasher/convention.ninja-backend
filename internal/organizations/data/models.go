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

type Venue struct {
	data.SnowflakeModel
	Name           string      `json:"name"`
	OrganizationId int64       `json:"organizationId,string"`
	Address        string      `json:"address"`
	Rooms          []VenueRoom `json:"rooms,omitempty"`
}

type VenueRoom struct {
	data.SnowflakeModel
	Name           string `json:"name"`
	VenueId        int64  `json:"venueId,string"`
	OrganizationId int64  `json:"organizationId,string"`
}
