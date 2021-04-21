package database

import "gorm.io/gorm"

var DB *gorm.DB

func InitDB(driver, dsn string) (err error) {
	switch driver {
	case "sqlite":
		DB, err = NewSqlite3(dsn)
	}
	return err
}
