package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/dereking/grest/config"

	_ "github.com/go-sql-driver/mysql"
)

func GetReadConn() (*sql.DB, error) {
	var dbName string

	if host, found := config.AppConfig.String("db.mysql.hostRead"); found {
		//host = 127.0.0.1:3306
		user := config.AppConfig.StringDefault("db.mysql.hostRead.user", "root")
		psw := config.AppConfig.StringDefault("db.mysql.hostRead.psw", "")
		if dbName, found = config.AppConfig.String("db.mysql.hostRead.dbName"); !found {
			panic("db.mysql.hostWrite assigned ,but db.mysql.hostRead.dbName not found!")
		}
		maxOpenConns := config.AppConfig.IntDefault("db.mysql.hostRead.maxOpenConns", 200)
		maxIdleConns := config.AppConfig.IntDefault("db.mysql.hostRead.maxIdleConns", 100)

		return NewMysqlPool(host, user, psw, dbName, maxOpenConns, maxIdleConns)
	}
	return nil, errors.New("db.mysql.hostRead not found in config file.")
}
func GetWriteConn() (*sql.DB, error) {

	var dbName string
	if host, found := config.AppConfig.String("db.mysql.hostWrite"); found {
		//host = 127.0.0.1:3306
		user := config.AppConfig.StringDefault("db.mysql.hostWrite.user", "root")
		psw := config.AppConfig.StringDefault("db.mysql.hostWrite.psw", "")
		if dbName, found = config.AppConfig.String("db.mysql.hostWrite.dbName"); !found {
			panic("db.mysql.hostWrite assigned ,but db.mysql.hostWrite.dbName not found!")
		}
		maxOpenConns := config.AppConfig.IntDefault("db.mysql.hostWrite.maxOpenConns", 200)
		maxIdleConns := config.AppConfig.IntDefault("db.mysql.hostWrite.maxIdleConns", 100)

		return NewMysqlPool(host, user, psw, dbName, maxOpenConns, maxIdleConns)
	}
	return nil, errors.New("db.mysql.hostWrite not found in config file.")
}

func NewMysqlPool(host string, user string, psw string, dbName string,
	MaxOpenConns int, MaxIdleConns int) (*sql.DB, error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		user, psw, host, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
		return nil, err
	}
	db.SetMaxOpenConns(MaxOpenConns)
	db.SetMaxIdleConns(MaxIdleConns)
	db.Ping()

	return db, nil
}
