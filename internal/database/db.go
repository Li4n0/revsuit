package database

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	Driver DriverType
)

type DriverType = string

const Sqlite = "sqlite"
const Mysql = "mysql"
const Postgres = "postgres"
const UnknowDatabase = "unknown database"

func InitDB(dsn string) (err error) {
	dbName, dbDsn := dbType(dsn)

	switch dbName {
	case Sqlite:
		DB, err = NewSqlite3(dbDsn)
	case Mysql:
		DB, err = NewMysql(dbDsn)
	case Postgres:
		DB, err = NewPostgres(dbDsn)
	default:
		err = errors.New("unsupported database")
	}
	return err
}

func dbType(dsn string) (DriverType, string) {
	if strings.HasSuffix(dsn, ".db") || strings.HasPrefix(dsn, "sqlite3://") {
		return Sqlite, strings.TrimPrefix(dsn, "sqlite3://")
	} else if strings.HasPrefix(dsn, "postgres://") {
		return Postgres, strings.TrimPrefix(dsn, "postgres://")
	} else if strings.HasPrefix(dsn, "mysql://") {
		return Mysql, strings.TrimPrefix(dsn, "mysql://")
	} else if strings.Contains(dsn, "@tcp") { //兼容上个版本的写法
		return Mysql, dsn
	}
	return UnknowDatabase, dsn
}
