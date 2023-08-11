package database

import (
	"app3.1/config"
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"time"
)

type User struct {
	ID        int64   `db:"ID" json:"ID"`
	Nickname  string  `db:"NickName" json:"Nickname" validate:"required"`
	FirstName string  `db:"FirstName" json:"FirstName" validate:"required"`
	LastName  string  `db:"LastName" json:"LastName" validate:"required"`
	Password  string  `db:"Password" json:"Password" validate:"required"`
	CreatedAt string  `db:"created_at" json:"created_at"`
	UpdatedAt *string `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt *string `db:"deleted_at" json:"deleted_at,omitempty"`
}

type DbInterface interface {
	FindByUserName(userName string) (*User, error)
	IsUserDeleted(userName string) bool
	InsertUser(user User) (string, error)
	UpdateUser(userName string, user User) (string, error)
	FindUsers() (*[]User, error)
	DeleteUserByNick(userName string) (string, error)
}

type UserDatabase struct {
	Connection *sql.DB
}

func NewDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "database/test.db")
	if err != nil {
		log.Warn().Err(err).Msg(" can`t connect to SQLite")
		return nil
	}

	return db
}

func (db *UserDatabase) FindByUserName(userName string) (*User, error) {
	sqlSelect := `SELECT * FROM User WHERE NickName = ?`
	var selectedUser User

	row := db.Connection.QueryRow(sqlSelect, userName)
	err := row.Scan(&selectedUser.ID, &selectedUser.Nickname, &selectedUser.FirstName,
		&selectedUser.LastName, &selectedUser.Password, &selectedUser.CreatedAt,
		&selectedUser.UpdatedAt, &selectedUser.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle "not found" scenario
			return nil, errors.New("user not found")
		}
		log.Warn().Err(err).Msg(" can`t find user")
		return nil, err
	}

	return &selectedUser, nil
}

func (db *UserDatabase) IsUserDeleted(userName string) bool {
	sqlSelect := `SELECT NickName FROM User WHERE NickName = ? AND DeletedAt IS NOT NULL`
	var checkedUser string
	row := db.Connection.QueryRow(sqlSelect, userName)
	err := row.Scan(&checkedUser)
	if err != nil {
		if err == sql.ErrNoRows {
			return true
		}
	}
	return false

}

func (db *UserDatabase) InsertUser(user User) (string, error) {
	formattedTime := time.Now().Format("2006.01.02 15:04")

	sqlInsert := "INSERT INTO User (NickName, FirstName, LastName, Password, CreatedAt) VALUES (?, ?, ?, ?, ?)"

	hashedPassword, err := config.Hash(user.Password)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t hashed user`s password")
	}

	_, err = db.Connection.Exec(sqlInsert, user.Nickname, user.FirstName, user.LastName, hashedPassword, formattedTime)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t insert user")
		return "", err
	}

	return user.Nickname, nil

}

func (db *UserDatabase) UpdateUser(userName string, user User) (string, error) {

	hashedPassword, err := config.Hash(user.Password)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t hashed user`s password")
	}
	formattedTime := time.Now().Format("2006.01.02 15:04")
	sqlUpdate := "UPDATE User SET NickName = ?, FirstName = ?, LastName = ?, Password = ?, UpdatedAt = ? WHERE Nickname = ?"

	_, err = db.Connection.Exec(sqlUpdate, user.Nickname, user.FirstName, user.LastName, hashedPassword, formattedTime, userName)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t update user`s data")
		return "", err
	}

	return user.Nickname, nil

}

func (db *UserDatabase) FindUsers() (*[]User, error) {
	sqlSelect := "SELECT * FROM User"
	rows, err := db.Connection.Query(sqlSelect)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t find users")
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var singleUser User
		err := rows.Scan(&singleUser.ID, &singleUser.Nickname, &singleUser.FirstName,
			&singleUser.LastName, &singleUser.Password, &singleUser.CreatedAt,
			&singleUser.UpdatedAt, &singleUser.DeletedAt)

		if err != nil {
			return nil, err
		}
		users = append(users, singleUser)
	}
	return &users, nil
}

func (db *UserDatabase) DeleteUserByNick(userName string) (string, error) {
	formattedTime := time.Now().Format("2006.01.02 15:04")
	sqlSoftDelete := "UPDATE User SET (DeletedAt) = (?) WHERE NickName = ?"

	_, err := db.Connection.Exec(sqlSoftDelete, formattedTime, userName)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t delete user`s data")
		return "", err
	}

	return userName, nil
}
