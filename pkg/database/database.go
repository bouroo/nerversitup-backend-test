package database

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func BuildConnStr(driver string, dbName string, host string, port int, user string, passwd string) (dsn string, err error) {
	switch driver {
	case "sqlite":
		dsn = fmt.Sprintf("file:%s?cache=shared", dbName)
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8&parseTime=true", user, passwd, host, port, dbName)
	case "postgres":
		dsn = fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s", user, passwd, host, port, dbName)
	case "sqlserver":
		dsn = fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", user, passwd, host, port, dbName)
	default:
		err = errors.New("Not support DB_DRIVER")
		return
	}
	return
}

func Connect(driver string, dsn string) (dbConn *gorm.DB, err error) {
	switch driver {
	case "sqlite":
		dbConn, err = gorm.Open(sqlite.Open(dsn))
	case "mysql":
		dbConn, err = gorm.Open(mysql.Open(dsn))
	case "postgres":
		dbConn, err = gorm.Open(postgres.Open(dsn))
	case "sqlserver":
		dbConn, err = gorm.Open(sqlserver.Open(dsn))
	default:
		err = errors.New("Not support DB_DRIVER")
		return
	}
	err = applyDbPoolConfig(dbConn)
	return
}

func applyDbPoolConfig(dbConn *gorm.DB) (err error) {
	// Get generic database object sql.DB to use its functions
	sqlDB, err := dbConn.DB()
	if err != nil {
		return
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(viper.GetInt("db_min_conn"))

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(viper.GetInt("db_max_conn"))

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(viper.GetDuration("db_conn_time"))
	return
}
