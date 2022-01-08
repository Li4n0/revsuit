package database

import (
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func NewSqlite3(dsn string) (*gorm.DB, error) {
	if strings.Contains(dsn, "?") {
		dsn += "&_synchronous=1&_journal_mode=WAL"
	} else {
		dsn += "?_synchronous=1&_journal_mode=WAL"
	}

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.Exec("PRAGMA foreign_keys=ON")
	return db, nil
}
