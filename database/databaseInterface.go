package database

import "time"

type DbInterface interface {
	FindByID(ID int64) (*User, error)
	InsertUser(user User) (int64, error)
	UpdateUser(ID int64, user User) (int64, error)
	FindUsers() (*[]User, error)
	DeleteUserByID(ID int64) error
	FindByNicknameToGetUserPassword(nickname string) (*User, error)
	GetUserRating(ID int64) (*UserRating, error)
	VoteForUser(ID int64) error
	WriteUserVotes(voteTime map[string]time.Time, userVotes map[string][]int64, userName string) error
}
