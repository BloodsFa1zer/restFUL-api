package serviceLayer

import (
	"app3.1/database"
)

type UserServiceInterface interface {
	UserValidation(user database.User) error
	CreateUser(user database.User) (int64, error)
	GetUser(userID int64) (*database.User, error)
	EditUser(ID int64, user database.User) (int64, error)
	GetAllUsers() (*[]database.User, error)
	DeleteUser(userID int64) error
	GetUserByName(nickname, password string) (string, error)
	IsUserHavePermission(roleToCheck string, user interface{}) bool
}
