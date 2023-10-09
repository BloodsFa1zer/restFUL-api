package redisDB

import "app3.1/database"

type ClientRedisInterface interface {
	GetUser(ID int64) (*database.User, error)
	SetUser(user database.User) error
	SetUsers(users []database.User) error
	GetUsers() (*[]database.User, error)
}
