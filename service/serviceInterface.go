package service

import (
	"app3.1/database"
	"time"
)

type UserServiceInterface interface {
	UserValidation(user database.User) error
	CreateUser(user database.User) (int64, error, int)
	GetUser(userID int64) (*database.User, error, int)
	EditUser(ID int64, user database.User) (int64, error, int)
	GetAllUsers() (*[]database.User, error, int)
	DeleteUser(userID int64) (error, int)
	CreateToken(user database.User) (string, error, int)
	GetUserNameViaToken(user interface{}) string
	IsUserHavePermission(roleToCheck string, user interface{}) (bool, int)
	isUserAllowedToVote(userName string, voteID int) (map[string][]int, map[string]time.Time, error)
	Vote(userID int64, userName string) (error, int)
	GetUserRate(ID int64) (*database.UserRating, error, int)
	GetAllUsersRate() (*[]database.UserRating, error, int)
	GetUserRateModerator(ID int64) (*database.UserRating, error, int)
}
