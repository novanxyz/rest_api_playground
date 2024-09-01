package config

import (
	"database/sql"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"fmt"

	"novanxyz/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	driver   = utils.Getenv("DB_TYPE", "mysql")
	host     = utils.Getenv("DB_HOST", "127.0.0.1")
	port     = utils.Getenv("DB_PORT", "3306")
	user     = utils.Getenv("DB_USER")
	password = utils.Getenv("DB_PASS")
	dbName   = utils.Getenv("DB_NAME")
)

func DatabaseConnection() *gorm.DB {
	port, err := strconv.Atoi(utils.Getenv("DB_PORT", "3306"))
	utils.ErrorPanic(err)

	conn, err := sql.Open(driver, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbName))
	utils.ErrorPanic(err)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: conn,
	}), &gorm.Config{})
	utils.ErrorPanic(err)

	return db
}
