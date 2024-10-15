package pkg

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Database struct {
	handler *sql.DB
}

func NewDatabase(dbFile string) *Database {
	handler, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalln(err)
	}
	db := Database{handler: handler}
	return &db
}
