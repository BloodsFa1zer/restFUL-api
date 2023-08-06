package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

type Database struct {
	connection *sql.DB
}

func Connection() *Database {
	db, err := sql.Open("sqlite3", "database/test.db")
	if err != nil {
		log.Warn().Err(err).Msg(" can`t connect to SQLite")
		return nil
	}
	database := Database{
		connection: db,
	}
	return &database
}
