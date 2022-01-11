package data

import "convention.ninja/internal/data"

type User struct {
	data.SnowflakeModel
	Name       string
	Email      string
	FirebaseId string
}
