package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSqlite3(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.Exec("PRAGMA foreign_keys=ON")
	return db, nil
}
