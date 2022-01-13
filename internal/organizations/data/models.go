package data

import "convention.ninja/internal/data"

type Organization struct {
	data.SnowflakeModel
	Name    string
	OwnerId int64
}

func init() {
	_ = data.GetConn().AutoMigrate(&Organization{})
}
