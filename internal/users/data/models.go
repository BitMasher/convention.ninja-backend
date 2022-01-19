package data

import "convention.ninja/internal/data"

type User struct {
	data.SnowflakeModel
	Name          string `json:"-"`
	DisplayName   string `json:"displayName"`
	Email         string `json:"-"`
	EmailVerified bool   `json:"-"`
	FirebaseId    string `json:"-"`
}

func Migrate() {
	_ = data.GetConn().AutoMigrate(&User{})
}
