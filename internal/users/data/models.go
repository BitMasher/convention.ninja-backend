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

// TODO this should not run in server mode, should be part of deployment
func init() {
	_ = data.GetConn().AutoMigrate(&User{})
}
