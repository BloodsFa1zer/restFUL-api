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
	CreateToken(user database.User) (string, error, int)
	GetUserIDViaToken(user interface{}) (int64, error)
	IsUserHavePermission(roleToCheck string, user interface{}) (bool, int)
	isUserAllowedToVote(voterID, userID int) (bool, error)
	PostVoteFor(userID, voterID int) (error, int)
	PostVoteAgainst(userID, voterID int) (error, int)
	DeleteVote(userID, voterID int) (error, int)
	ChangeVote(userID, voterID int) (error, int)
}
