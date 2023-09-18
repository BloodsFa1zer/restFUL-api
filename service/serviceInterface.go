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
	CreateToken(nickname, password string) (string, error, int)
	GetUserNameViaToken(user interface{}) string
	IsUserHavePermission(roleToCheck string, user interface{}) (bool, int)
	Registration(username, firstName, surName, password string) (int, error, int)
	Vote(userID int64, userName string) (error, int)
	GetUserRate(ID int64) (*database.UserRating, error, int)
	//	isUserAllowedToVoteAgain(voteTime time.Time) bool
	//	isUserAllowedToVoteForThatCandidate(userVote map[string][]int64, userName string, desiredID int64) bool
}
