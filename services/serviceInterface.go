package services

import (
	"app3.1/database"
	"github.com/go-playground/validator/v10"
)

type UserServiceInterface interface {
	UserValidation(user database.User) error
	Create(user database.User) (int64, error)
	Get(userID int64) (*database.User, error)
	Edit(ID int64, user database.User) (int64, error)
	GetAll() (*[]database.User, error)
	Delete(userID int64) error
	GetPasswordByName(nickname string) (*database.User, error)
}

type UserService struct {
	DbUser   database.DbInterface
	validate *validator.Validate
}

// Not a conventional way of creating constructors, but did it to avoid calling "app3.1/database" in routes.go
func NewUserService() *UserService {
	return &UserService{DbUser: database.NewUserDatabase(), validate: validator.New()}
}
