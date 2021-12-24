package database

import (
	"errors"
	"gorm.io/gorm"
	"strings"
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
	switch dbType(dsn) {
	case Sqlite:
		sqliteDsn := strings.TrimLeft(strings.TrimLeft(dsn, "sqlite3://"), "sqlite://")
		DB, err = NewSqlite3(sqliteDsn)
	case Mysql:
		mysqlDsn := strings.TrimLeft(dsn, "mysql://")
		DB, err = NewMysql(mysqlDsn)
	case Postgres:
		postgresDsn := strings.TrimLeft(dsn, "postgres://")
		DB, err = NewPostgres(postgresDsn)
	default:
		return errors.New("unsupported database")
	}
	return err
}

func dbType(dsn string) DriverType {

	if strings.HasSuffix(dsn, ".db") || strings.HasPrefix(dsn, "sqlite://") || strings.HasPrefix(dsn, "sqlite3://") {
		return Sqlite
	} else if strings.HasPrefix(dsn, "postgres://") {
		return Postgres
	} else if strings.HasPrefix(dsn, "mysql://") {
		return Mysql
	} else if strings.Contains(dsn, "@tcp") { //兼容前面的
		return Mysql
	}
	return UnknowDatabase
}
