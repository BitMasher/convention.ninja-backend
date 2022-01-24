package main

import (
	"convention.ninja/internal/data"
	"os"
)

func main() {
	switch os.Args[1] {
	case "migrate":
		migrateCommand()
		break
	}

}

func migrateCommand() {
	dsn := os.Getenv("SQL_DSN")
	err := data.Connect(dsn)
	if err != nil {
		panic(err)
	}

}
