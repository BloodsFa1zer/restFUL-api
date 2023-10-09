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
	Role      string  `db:"role" json:"Role"`
	CreatedAt string  `db:"created_at" json:"CreatedAt"`
	UpdatedAt *string `db:"updated_at" json:"UpdatedAt,omitempty"`
	DeletedAt *string `db:"deleted_at" json:"DeletedAt,omitempty"`
	Rating    int64   `db:"rating" json:"Rating"`
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
	sqlSelect := `SELECT Users.*, sum(Voting.vote_value) FROM Users
	LEFT JOIN Voting ON Voting.user_id=Users.ID WHERE ID = ?;`
	var num sql.NullInt64
	var selectedUser User

	row := db.Connection.QueryRow(sqlSelect, ID)
	err := row.Scan(&selectedUser.ID, &selectedUser.Nickname, &selectedUser.FirstName,
		&selectedUser.LastName, &selectedUser.Password, &selectedUser.CreatedAt,
		&selectedUser.UpdatedAt, &selectedUser.DeletedAt, &selectedUser.Role, &num)
	selectedUser.Rating = num.Int64
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		log.Warn().Err(err).Msg(" can`t find user")
		return nil, err
	}

	return &selectedUser, nil
}

func (db *UserDatabase) FindByNicknameToGetUserPassword(nickname string) (*User, error) {
	sqlSelect := `SELECT * FROM Users WHERE nick_name = ? AND deleted_at == 'NULL'`
	var selectedUser User

	row := db.Connection.QueryRow(sqlSelect, nickname)
	err := row.Scan(&selectedUser.ID, &selectedUser.Nickname, &selectedUser.FirstName,
		&selectedUser.LastName, &selectedUser.Password, &selectedUser.CreatedAt,
		&selectedUser.UpdatedAt, &selectedUser.DeletedAt, &selectedUser.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle "not found" scenario
			return nil, err
		}
		log.Warn().Err(err).Msg(" can`t find user")
		return nil, err
	}

	return &selectedUser, nil
}

func (db *UserDatabase) InsertUser(user User) (int64, error) {
	formattedTime := time.Now().Format("2006.01.02 15:04")

	sqlInsert := "INSERT INTO Users (nick_name, first_name, last_name, password, created_at) VALUES (?, ?, ?, ?, ?)"

	hashedPassword := hash.Hash(user.Password)

	result, err := db.Connection.Exec(sqlInsert, user.Nickname, user.FirstName, user.LastName, hashedPassword, formattedTime)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t insert user")
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		log.Warn().Err(err).Msg(" can`t find userID")
		return 0, err
	}

	return userID, nil
}

func (db *UserDatabase) UpdateUser(ID int64, user User) (int64, error) {

	hashedPassword := hash.Hash(user.Password)

	formattedTime := time.Now().Format("2006.01.02 15:04")
	sqlUpdate := "UPDATE Users SET nick_name = ?, first_name = ?, last_name = ?, password = ?, updated_at = ? WHERE ID = ? AND deleted_at == 'NULL'"

	result, err := db.Connection.Exec(sqlUpdate, user.Nickname, user.FirstName, user.LastName, hashedPassword, formattedTime, ID)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t update user`s data")
		return 0, err
	}

	affectedRow, err := result.RowsAffected()
	if err != nil {
		log.Warn().Err(err).Msg(" error getting affected rows")
		return 0, err
	}

	if affectedRow == 0 {
		log.Warn().Msg(" no rows affected")
		return 0, sql.ErrNoRows
	}

	return ID, nil
}

func (db *UserDatabase) FindUsers() (*[]User, error) {

	sqlSelect := "SELECT * FROM Users WHERE deleted_at == 'NULL'"
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
			&singleUser.UpdatedAt, &singleUser.DeletedAt, &singleUser.Role)
		if err != nil {
			return nil, err
		}

		users = append(users, singleUser)
	}
	return &users, nil
}

func (db *UserDatabase) DeleteUserByID(ID int64) error {
	formattedTime := time.Now().Format("2006.01.02 15:04")
	sqlSoftDelete := "UPDATE Users SET (deleted_at) = (?) WHERE ID = ? AND deleted_at == 'NULL'"

	result, err := db.Connection.Exec(sqlSoftDelete, formattedTime, ID)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t delete user`s data")
		return err
	}

	affectedRow, err := result.RowsAffected()
	if err != nil {
		log.Warn().Err(err).Msg(" error getting affected rows")
		return err
	}

	if affectedRow == 0 {
		log.Warn().Msg(" no rows affected")
		return sql.ErrNoRows
	}

	return nil
}

type UserRating struct {
	UserID    int64  `db:"user_id" json:"ID,omitempty"`
	VoterID   int64  `db:"voter_id" json:"Votes,omitempty"`
	VoteTime  string `db:"updated_at" json:"VoteTime,omitempty"`
	VoteValue int    `db:"vote_value" json:"VoteValue,omitempty"`
}

func (db *UserDatabase) WriteUserVotes(userID, voterID, voteValue int) error {
	sqlInsertVotes := "INSERT INTO Voting (user_id, voter_id, updated_at, vote_value) VALUES (?, ?, ?, ?)"
	result, err := db.Connection.Exec(sqlInsertVotes, userID, voterID, time.Now(), voteValue)

	if err != nil {
		log.Warn().Err(err).Msg(" can`t vote for that user")
		return err
	}

	affectedRow, err := result.RowsAffected()
	if err != nil {
		log.Warn().Err(err).Msg(" error getting affected rows")
		return err
	}

	if affectedRow == 0 {
		log.Warn().Msg(" no rows affected")
		return sql.ErrNoRows
	}

	return nil
}

func (db *UserDatabase) IsSuchVoteExists(userID, voterID int) (bool, error) {
	sqlSelectVotes := "SELECT exists(SELECT 1 FROM Voting WHERE user_id = ? AND voter_id = ?)"
	var exists bool

	row := db.Connection.QueryRow(sqlSelectVotes, userID, voterID)
	err := row.Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Warn().Err(err).Msg(" can`t find user")
		return false, err
	}

	return exists, nil
}

func (db *UserDatabase) WithdrawVote(userID, voterID int) error {
	sqlDeleteVote := "DELETE FROM VOTING WHERE user_id = ? AND voter_id = ?"
	result, err := db.Connection.Exec(sqlDeleteVote, userID, voterID)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t delete user vote")
		return err
	}

	affectedRow, err := result.RowsAffected()
	if err != nil {
		log.Warn().Err(err).Msg(" error getting affected rows")
		return err
	}

	if affectedRow == 0 {
		log.Warn().Msg(" no rows affected")
		return sql.ErrNoRows
	}

	return nil
}

func (db *UserDatabase) ChangeVote(newVoteID, voterID int) error {
	sqlVotesCheckTime := "SELECT user_id FROM Voting WHERE voter_id = ? ORDER BY updated_at DESC LIMIT 1"
	votedID := 0
	row := db.Connection.QueryRow(sqlVotesCheckTime, voterID)
	err := row.Scan(&votedID)
	if err != nil {
		log.Warn().Err(err)
		return err
	}

	if newVoteID == voterID {
		return errors.New(" you are not allowed to vote for yourself")
	}
	exists, err := db.IsSuchVoteExists(newVoteID, voterID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if exists {
		return errors.New("you can`t change your vote to the user you`ve already voted for")
	}

	sqlUpdateVote := "UPDATE VOTING SET (user_id) = ? WHERE user_id = ? AND voter_id = ?"
	_, err = db.Connection.Exec(sqlUpdateVote, newVoteID, votedID, voterID)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t change user vote")
		return err
	}

	return nil
}

func (db *UserDatabase) GetUserLastVoteTime(voterID int) (string, error) {
	sqlVotesCheckTime := "SELECT updated_at FROM Voting WHERE voter_id = ? ORDER BY updated_at DESC LIMIT 1"
	voteTime := ""
	row := db.Connection.QueryRow(sqlVotesCheckTime, voterID)
	err := row.Scan(&voteTime)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Err(err).Msg(" can`t find user vote")
			return "0", err
		}
		log.Warn().Err(err)
		return "", err
	}

	return voteTime, nil
}
