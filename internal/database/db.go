package database

import (
	"errors"
	"strings"
	"sync"

	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	Driver DriverType
	Locker sync.Mutex
)

type DriverType = string

const Sqlite = "sqlite"
const Mysql = "mysql"
const Postgres = "postgres"
const UnknownDatabase = "unknown database"

func InitDB(dsn string) (err error) {
	dbName, dbDsn := dbType(dsn)

	switch dbName {
	case Sqlite:
		DB, err = NewSqlite3(dbDsn)
		Driver = Sqlite
	case Mysql:
		DB, err = NewMysql(dbDsn)
		Driver = Mysql
	case Postgres:
		DB, err = NewPostgres(dbDsn)
		Driver = Postgres
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
	return UnknownDatabase, dsn
}
