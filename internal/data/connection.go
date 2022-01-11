package data

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect(dsn string) error {
	locDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:       true,
		AllowGlobalUpdate: false,
	})
	if err != nil {
		return err
	}
	sqlDb, err := locDb.DB()
	if err != nil {
		return err
	}
	sqlDb.SetMaxIdleConns(5)
	sqlDb.SetMaxOpenConns(15)
	db = locDb
	return nil
}

func GetConn() *gorm.DB {
	return db
}
