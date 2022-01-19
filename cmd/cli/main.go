package main

import (
	"convention.ninja/internal/data"
	data4 "convention.ninja/internal/inventory/data"
	data3 "convention.ninja/internal/organizations/data"
	data2 "convention.ninja/internal/users/data"
	"os"
)

func main() {
	dsn := os.Getenv("SQL_DSN")
	err := data.Connect(dsn)
	if err != nil {
		panic(err)
	}
	data2.Migrate()
	data3.Migrate()
	data4.Migrate()
}
