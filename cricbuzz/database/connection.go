package database

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type DBConnection struct {
	conn *sql.DB
	err  error
}

var once sync.Once
var instance *DBConnection

func NewDBConnection() *DBConnection {
	conn, err := sql.Open("sqlite3", "./cricbuzz.db")
	return &DBConnection{
		conn: conn,
		err:  err,
	}
}

func GetDBConnection() *DBConnection {
	once.Do(func() {
		instance = NewDBConnection()
	})
	return instance
}
