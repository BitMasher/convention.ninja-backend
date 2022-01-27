package data

import "convention.ninja/internal/data"

type User struct {
	data.SnowflakeModel
	Name          string `json:"name,omitempty"`
	DisplayName   string `json:"displayName"`
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"-"`
	FirebaseId    string `json:"-"`
}
