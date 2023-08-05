package database

import (
	"app3.1/password"
	"database/sql"
	"errors"
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

//func (db *Database) CheckUserNickName(user User) error {
//	sqlSelect := `SELECT COUNT(*) FROM User WHERE NickName = user.NickName`
//	count := 0
//	err := db.connection.QueryRow(sqlSelect).Scan(&count)
//	if err != nil {
//		log.Warn().Err(err).Msg(" can`t find user")
//		return err
//	}
//
//	if count > 0 {
//		fmt.Println("that NickName is already used", user.Nickname)
//		return errors.New("that NickName is already used")
//	}
//
//	return nil
//}

func (db *Database) FindUser(userName string) (*User, error) {
	sqlSelect := `SELECT * FROM User WHERE NickName = ?`
	var selectedUser User
	row := db.connection.QueryRow(sqlSelect, userName)
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

func (db *Database) InsertUser(user User) (string, error) {
	time := time.Now().Format("2006.01.02 15:04")

	sqlInsert := "INSERT INTO User (NickName, FirstName, LastName, Password, CreatedAt) VALUES (?, ?, ?, ?, ?)"

	hashedPassword, err := password.Hash(user.Password)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t hashed user`s password")
	}

	_, err = db.connection.Exec(sqlInsert, user.Nickname, user.FirstName, user.LastName, hashedPassword, time)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t insert user")
		return "", err
	}

	return user.Nickname, nil

}

// BasicAuth required to execute UpdateUser!!!
func (db *Database) UpdateUser(userName string, user User) (string, error) {
	hashedPassword, err := password.Hash(user.Password)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t hashed user`s password")
	}
	time := time.Now().Format("2006.01.02 15:04")
	sqlUpdate := "UPDATE User SET NickName = ?, FirstName = ?, LastName = ?, Password = ?, UpdatedAt = ? WHERE Nickname = ?"

	_, err = db.connection.Exec(sqlUpdate, user.Nickname, user.FirstName, user.LastName, hashedPassword, time, userName)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t update user`s data")
		return "", err
	}

	return user.Nickname, nil

}

func (db *Database) FindUsers() (*[]User, error) {
	sqlSelect := "SELECT * FROM User"
	rows, err := db.connection.Query(sqlSelect)
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

func (db *Database) SoftDeleteUser(userName string) (string, error) {
	time := time.Now().Format("2006.01.02 15:04")
	sqlSoftDelete := "UPDATE User SET (DeletedAt) = (?) WHERE NickName = ?"

	_, err := db.connection.Exec(sqlSoftDelete, time, userName)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t delete user`s data")
		return "", err
	}

	return userName, nil
}

// BasicAuth required to execute DeleteUser!!!
func (db *Database) DeleteUser(userName string) error {
	sqlDelete := "DELETE FROM User WHERE NickName = ?"

	_, err := db.connection.Exec(sqlDelete, userName)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t delete user`s data")
		return err
	}
	return nil
}

//func (cl *ClientConnection) DeleteUser(id *primitive.ObjectID) (*mongo.DeleteResult, error) {
//	result, err := cl.collection.DeleteOne(context.Background(), bson.M{"_id": id})
//	if err != nil {
//		log.Warn().Err(err).Msg(" can`t delete user`s data")
//		return &mongo.DeleteResult{}, err
//	}
//
//	return result, nil
//}

//
//func (cl *ClientConnection) FindUser(field string, data any) (*User, error) {
//
//	result := cl.collection.FindOne(context.TODO(), bson.M{field: data})
//
//	// check for errors in the finding
//	if result.Err() != nil {
//		log.Warn().Err(result.Err()).Msg(" can`t find user")
//	}
//	log.Info().Msg(" find users")
//
//	// convert the cursor result to bson
//	var user User
//	// check for errors in the conversion
//	if err := result.Decode(&user); err == mongo.ErrNoDocuments {
//		log.Warn().Err(err).Msg(" no results to convert")
//		return nil, err
//	} else if err != nil {
//		log.Warn().Err(err).Msg(" can`t convert results")
//		return nil, err
//	}
//	return &user, nil
//}

//
//func (cl *ClientConnection) FindUsers() (*[]User, error) {
//
//	results, err := cl.collection.Find(context.TODO(), bson.M{})
//	if err != nil {
//		log.Warn().Err(results.Err()).Msg(" can`t find users")
//	}
//	log.Info().Msg(" find users")
//	// convert the cursor result to bson
//	var users []User
//
//	for results.Next(context.TODO()) {
//		var singleUser User
//		if err = results.Decode(&singleUser); err != nil {
//			return nil, err
//		}
//		users = append(users, singleUser)
//	}
//	return &users, nil
//}
//

//func (cl *ClientConnection) InsertUser(user User) (*mongo.InsertOneResult, error) {
//	//	password.Hash(user.Password)
//	time := time.Now().Format("2006.01.02 15:04")
//	userInfo := bson.D{
//		{"Nickname", user.Nickname},
//		{"FirstName", user.FirstName},
//		{"LastName", user.LastName},
//		{"Password", user.Password},
//		{"created_at", time},
//	}
//
//	result, err := cl.collection.InsertOne(context.TODO(), userInfo)
//	if err != nil {
//		log.Warn().Err(err).Msg(" can`t insert user`s data into database")
//		return nil, err
//	}
//	log.Info().Msg("successfully insert user`s data")
//	return result, nil
//}

//func (cl *ClientConnection) UpdateUser(id *primitive.ObjectID, user User) (*mongo.UpdateResult, error) {
//	//	password.Hash(user.Password)
//	time := time.Now().Format("2006.01.02 15:04")
//	update := bson.D{{"$set", bson.D{
//		{"Nickname", user.Nickname},
//		{"FirstName", user.FirstName},
//		{"LastName", user.LastName},
//		{"Password", user.Password},
//		{"updated_at", time},
//	}}}
//	result, err := cl.collection.UpdateByID(context.Background(), id, update)
//	if err != nil {
//		log.Warn().Err(err).Msg(" can`t update user`s data")
//		return nil, err
//	}
//	return result, nil
//
//}

//func (cl *ClientConnection) DeleteUser(id *primitive.ObjectID) (*mongo.DeleteResult, error) {
//	result, err := cl.collection.DeleteOne(context.Background(), bson.M{"_id": id})
//	if err != nil {
//		log.Warn().Err(err).Msg(" can`t delete user`s data")
//		return &mongo.DeleteResult{}, err
//	}
//
//	return result, nil
//}
