package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDatabase() *sql.DB {
	dsn := "root:2007hadi@tcp(localhost:3306)/angkringan"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	DB = db
	return db
}
