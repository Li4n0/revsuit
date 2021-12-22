package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysql(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
