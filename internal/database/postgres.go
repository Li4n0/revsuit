package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"vitess.io/vitess/go/vt/log"
)

func NewPostgres(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("连接数据库异常: %v", err.Error())
		os.Exit(0)
	}
	return db, nil
}
