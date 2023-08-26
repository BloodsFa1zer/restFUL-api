package database

import (
	"app3.1/config"
	"app3.1/hash"
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"time"
)

type User struct {
	ID        int64   `db:"ID" json:"ID"`
	Nickname  string  `db:"nick_name" json:"Nickname" validate:"required"`
	FirstName string  `db:"first_name" json:"FirstName" validate:"required"`
	LastName  string  `db:"last_name" json:"LastName" validate:"required"`
	Password  string  `db:"password" json:"Password" validate:"required"`
	CreatedAt string  `db:"created_at" json:"created_at"`
	UpdatedAt *string `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt *string `db:"deleted_at" json:"deleted_at,omitempty"`
}

type DbInterface interface {
	FindByID(ID int64) (*User, error)
	IsUserDeleted(ID int64) bool
	InsertUser(user User) (int64, error)
	UpdateUser(ID int64, user User) (int64, error)
	FindUsers() (*[]User, error)
	DeleteUserByID(ID int64) error
}

type UserDatabase struct {
	Connection *sql.DB
}

func NewUserDatabase() *UserDatabase {
	cfg := config.LoadENV("config/.env")
	cfg.ParseENV()

	db, err := sql.Open(cfg.DbName, cfg.DbPath)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t connect to SQLite")
		return nil
	}
	log.Info().Msg("successfully connected to SQLite")

	return &UserDatabase{Connection: db}
}

func (db *UserDatabase) FindByID(ID int64) (*User, error) {
	sqlSelect := `SELECT * FROM Users WHERE ID = ?`
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

func (db *UserDatabase) IsUserDeleted(ID int64) bool {

	sqlSelect := `select ID from Users where ID = ? and deleted_at != 'NULL'`
	var selectedUserID int64
	row := db.Connection.QueryRow(sqlSelect, ID)
	err := row.Scan(&selectedUserID)

	if err == sql.ErrNoRows {
		// not found
		return true
	}

	return false
}

func (db *UserDatabase) InsertUser(user User) (int64, error) {
	formattedTime := time.Now().Format("2006.01.02 15:04")

	sqlInsert := "INSERT INTO Users (nick_name, first_name, last_name, password, created_at) VALUES (?, ?, ?, ?, ?)"

	hashedPassword, err := hash.Hash(user.Password)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t hash user`s password")
	}

	var selectedUser User

	row := db.Connection.QueryRow(sqlInsert, user.Nickname, user.FirstName, user.LastName, hashedPassword, formattedTime)
	if row.Err() != nil {
		log.Warn().Err(row.Err()).Msg(" can`t insert user")
		return 0, row.Err()
	}

	err = row.Scan(&selectedUser.ID, &selectedUser.Nickname, &selectedUser.FirstName,
		&selectedUser.LastName, &selectedUser.Password, &selectedUser.CreatedAt,
		&selectedUser.UpdatedAt, &selectedUser.DeletedAt)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t scan user`s password")
		return 0, err
	}

	return selectedUser.ID, nil

}

func (db *UserDatabase) UpdateUser(ID int64, user User) (int64, error) {

	hashedPassword, err := hash.Hash(user.Password)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t hashed user`s password")
	}
	formattedTime := time.Now().Format("2006.01.02 15:04")
	sqlUpdate := "UPDATE Users SET nick_name = ?, first_name = ?, last_name = ?, password = ?, updated_at = ? WHERE ID = ?"

	_, err = db.Connection.Exec(sqlUpdate, user.Nickname, user.FirstName, user.LastName, hashedPassword, formattedTime, ID)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t update user`s data")
		return 0, err
	}

	return ID, nil

}

func (db *UserDatabase) FindUsers() (*[]User, error) {
	sqlSelect := "SELECT * FROM Users"
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

func (db *UserDatabase) DeleteUserByID(ID int64) error {
	formattedTime := time.Now().Format("2006.01.02 15:04")
	sqlSoftDelete := "UPDATE Users SET (deleted_at) = (?) WHERE ID = ?"

	_, err := db.Connection.Exec(sqlSoftDelete, formattedTime, ID)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t delete user`s data")
		return err
	}

	return nil
}
