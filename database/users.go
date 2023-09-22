package database

import (
	"app3.1/config"
	"app3.1/hash"
	"database/sql"
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
	Rating    *int    `db:"rating"   json:"Rating,omitempty"`
	Votes     *[]int  `db:"votes"   json:"UserVotes,omitempty"`
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
	sqlSelect := `SELECT * FROM Users WHERE ID = ? AND deleted_at == 'NULL'`
	var selectedUser User

	row := db.Connection.QueryRow(sqlSelect, ID)
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

func (db *UserDatabase) FindByNicknameToGetUserPassword(nickname string) (*User, error) {
	sqlSelect := `SELECT * FROM Users WHERE nick_name = ? AND deleted_at == 'NULL'`
	var selectedUser User

	row := db.Connection.QueryRow(sqlSelect, nickname)
	err := row.Scan(&selectedUser.ID, &selectedUser.Nickname, &selectedUser.FirstName,
		&selectedUser.LastName, &selectedUser.Password, &selectedUser.CreatedAt,
		&selectedUser.UpdatedAt, &selectedUser.DeletedAt, &selectedUser.Role, &selectedUser.Rating)
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
			&singleUser.UpdatedAt, &singleUser.DeletedAt, &singleUser.Role, &singleUser.Rating)

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

func (db *UserDatabase) VoteForUser(ID int64) error {
	sqlAddVote := "UPDATE Voting SET rating = rating + ? WHERE ID = ?"
	result, err := db.Connection.Exec(sqlAddVote, 1, ID)

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

type UserRating struct {
	ID       int64  `db:"ID" json:"ID"`
	Nickname string `db:"nick_name" json:"Nickname" validate:"required"`
	Rating   int    `db:"rating"   json:"Rating"`
}

func (db *UserDatabase) GetUserRating(ID int64) (*UserRating, error) {
	sqlGetRating := "SELECT ID, nick_name FROM Users WHERE ID = ? AND deleted_at = 'NULL'"

	var user UserRating

	row := db.Connection.QueryRow(sqlGetRating, ID)
	err := row.Scan(&user.ID, &user.Nickname)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle "not found" scenario
			return nil, err
		}
		log.Warn().Err(err).Msg(" can`t find user")
		return nil, err
	}
	sqlGetRatingFromVoting := "SELECT rating FROM Voting WHERE ID = ?"
	row = db.Connection.QueryRow(sqlGetRatingFromVoting, ID)
	err = row.Scan(&user.Rating)

	return &user, nil
}

func (db *UserDatabase) WriteUserVotes(voteTime map[string]time.Time, userVotes map[string][]int64, userName string) error {
	sqlInsertVotes := "UPDATE Voting SET (votes, vote_time) = (?, ?)"

	result, err := db.Connection.Exec(sqlInsertVotes, userVotes[userName], voteTime[userName])

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
