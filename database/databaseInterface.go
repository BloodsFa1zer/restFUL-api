package database

type DbInterface interface {
	FindByID(ID int64) (*User, error)
	InsertUser(user User) (int64, error)
	UpdateUser(ID int64, user User) (int64, error)
	FindUsers() (*[]User, error)
	DeleteUserByID(ID int64) error
	FindByNicknameToGetUserPassword(nickname string) (*User, error)
	CountUserRate(userID int64) (int, error)
	WriteUserVotes(userID, voterID int) error
	GetUserVotes(userID, voterID int64) (string, error)
}
