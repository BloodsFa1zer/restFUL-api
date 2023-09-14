package service

import (
	"app3.1/database"
)

type UserServiceInterface interface {
	UserValidation(user database.User) error
	CreateUser(user database.User) (int64, error, int)
	GetUser(userID int64) (*database.User, error, int)
	EditUser(ID int64, user database.User) (int64, error, int)
	GetAllUsers() (*[]database.User, error, int)
	DeleteUser(userID int64) (error, int)
	GetToken(nickname, password string) (string, error, int)
	IsUserHavePermission(roleToCheck string, user interface{}) (bool, int)
}
