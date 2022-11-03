package database

import (
	"fmt"

	"github.com/kangata/db-test-go/helpers"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	driver   = helpers.Env("DB_DRIVER", "mysql")
	host     = helpers.Env("DB_HOST", "127.0.0.1")
	port     = helpers.Env("DB_PORT", "3306")
	username = helpers.Env("DB_USERNAME", "root")
	password = helpers.Env("DB_PASSWORD", "")
	database = helpers.Env("DB_DATABASE", "test__database")
	sslmode  = helpers.Env("DB_SSLMODE", "disable")
	timezone = helpers.Env("DB_TIMEZONE", "Local")
)

func mysqlConn() *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		username,
		password,
		host,
		port,
		database,
		timezone,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Cannot connect to database: %v", err))
	}

	return db
}

func postgresConn() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host,
		username,
		password,
		database,
		port,
		sslmode,
		timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Cannot connect to database: %v", err))
	}

	return db
}

func New() *gorm.DB {
	if driver == "mysql" {
		return mysqlConn()
	}

	if driver == "postgres" {
		return postgresConn()
	}

	panic(fmt.Sprintf("Invalid %s driver", driver))
}
