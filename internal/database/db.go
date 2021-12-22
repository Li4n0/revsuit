package database

import "gorm.io/gorm"

var (
	DB     *gorm.DB
	Driver DriverType
)

type DriverType = string

const Sqlite = "sqlite"
const Mysql = "mysql"

func InitDB(driver DriverType, dsn string) (err error) {
	Driver = driver
	switch driver {
	case Sqlite:
		DB, err = NewSqlite3(dsn)
	case Mysql:
		DB, err = NewMysql(dsn)
	}
	return err
}
