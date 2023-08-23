package database

import (
	"app3.1/hash"
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
	FindByID(ID int) (*User, error)
	IsUserDeleted(ID int) bool
	InsertUser(user User) (string, error)
	UpdateUser(ID int, user User) (int, error)
	FindUsers() (*[]User, error)
	DeleteUserByID(ID int) (int, error)
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

func (db *UserDatabase) FindByID(ID int) (*User, error) {
	sqlSelect := `SELECT * FROM User WHERE ID = ?`
	var selectedUser User

	row := db.Connection.QueryRow(sqlSelect, ID)
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

func (db *UserDatabase) IsUserDeleted(ID int) bool {
	sqlSelect := `SELECT NickName FROM User WHERE ID = ? AND DeletedAt IS NOT NULL`
	var checkedUser string
	row := db.Connection.QueryRow(sqlSelect, ID)
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

func (db *UserDatabase) UpdateUser(ID int, user User) (int, error) {

	hashedPassword, err := config.Hash(user.Password)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t hashed user`s password")
	}
	formattedTime := time.Now().Format("2006.01.02 15:04")
	sqlUpdate := "UPDATE User SET NickName = ?, FirstName = ?, LastName = ?, Password = ?, UpdatedAt = ? WHERE ID = ?"

	_, err = db.Connection.Exec(sqlUpdate, user.Nickname, user.FirstName, user.LastName, hashedPassword, formattedTime, ID)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t update user`s data")
		return 0, err
	}

	return ID, nil

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

func (db *UserDatabase) DeleteUserByID(ID int) (int, error) {
	formattedTime := time.Now().Format("2006.01.02 15:04")
	sqlSoftDelete := "UPDATE User SET (DeletedAt) = (?) WHERE ID = ?"

	_, err := db.Connection.Exec(sqlSoftDelete, formattedTime, ID)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t delete user`s data")
		return 0, err
	}

	return ID, nil
}
