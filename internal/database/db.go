package database

import "gorm.io/gorm"

var (
	DB     *gorm.DB
	Driver DriverType
)

type DriverType = string

const Sqlite = "sqlite"

func InitDB(driver DriverType, dsn string) (err error) {
	Driver = driver
	switch driver {
	case Sqlite:
		DB, err = NewSqlite3(dsn)
	}
	return err
}
