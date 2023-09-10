package database

import (
	"app3.1/config"
	"database/sql"
	"github.com/rs/zerolog/log"
)

type DbInterface interface {
	FindByID(ID int64) (*User, error)
	InsertUser(user User) (int64, error)
	UpdateUser(ID int64, user User) (int64, error)
	FindUsers() (*[]User, error)
	DeleteUserByID(ID int64) error
	FindByNicknameToGetUserPassword(nickname string) (*User, error)
}

type UserDatabase struct {
	Connection *sql.DB
}

func NewUserDatabase() *UserDatabase {
	cfg := config.LoadENV("config/.env")
	cfg.ParseENV()

	db, err := sql.Open(cfg.DbName, cfg.DbPath)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t connect to SQLite")
		return nil
	}
	log.Info().Msg("successfully connected to SQLite")

	return &UserDatabase{Connection: db}
}
