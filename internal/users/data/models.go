package data

import "convention.ninja/internal/data"

type User struct {
	data.SnowflakeModel
	Name          string
	DisplayName   string
	Email         string
	EmailVerified bool
	FirebaseId    string
}

func Migrate() {
	_ = data.GetConn().AutoMigrate(&User{})
}
