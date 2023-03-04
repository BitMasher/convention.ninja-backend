package data

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var db *sql.DB

func Connect(dsn string) error {
	db_, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	db = db_
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(100)
	return nil
}

func GetConn() *sql.DB {
	return db
}

func CloseConn() {
	if db != nil {
		_ = db.Close()
	}
}
