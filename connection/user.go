package connection

import (
	"context"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	Nickname  string             `bson:"Nickname,omitempty" json:"Nickname" validate:"required"`
	FirstName string             `bson:"FirstName,omitempty" json:"FirstName" validate:"required"`
	LastName  string             `bson:"LastName, omitempty" json:"LastName" validate:"required"`
	Password  string             `bson:"Password, omitempty" json:"Password" validate:"required"`
	CreatedAt string             `bson:"created_at" json:"created_at"`
	UpdatedAt string             `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	DeletedAt string             `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

func (cl *ClientConnection) FindUser(field string, data any) (*User, error) {

	result := cl.collection.FindOne(context.TODO(), bson.M{field: data})

	// check for errors in the finding
	if result.Err() != nil {
		log.Warn().Err(result.Err()).Msg(" can`t find user")
	}
	log.Info().Msg(" find users")

	// convert the cursor result to bson
	var user User
	// check for errors in the conversion
	if err := result.Decode(&user); err == mongo.ErrNoDocuments {
		log.Warn().Err(err).Msg(" no results to convert")
		return nil, err
	} else if err != nil {
		log.Warn().Err(err).Msg(" can`t convert results")
		return nil, err
	}
	return &user, nil
}

func (cl *ClientConnection) FindUsers() (*[]User, error) {

	results, err := cl.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Warn().Err(results.Err()).Msg(" can`t find users")
	}
	log.Info().Msg(" find users")
	// convert the cursor result to bson
	var users []User

	for results.Next(context.TODO()) {
		var singleUser User
		if err = results.Decode(&singleUser); err != nil {
			return nil, err
		}
		users = append(users, singleUser)
	}
	return &users, nil
}

func (cl *ClientConnection) InsertUser(user User) (*mongo.InsertOneResult, error) {
	//	password.Hash(user.Password)
	time := time.Now().Format("2006.01.02 15:04")
	userInfo := bson.D{
		{"Nickname", user.Nickname},
		{"FirstName", user.FirstName},
		{"LastName", user.LastName},
		{"Password", user.Password},
		{"created_at", time},
	}

	result, err := cl.collection.InsertOne(context.TODO(), userInfo)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t insert user`s data into database")
		return nil, err
	}
	log.Info().Msg("successfully insert user`s data")
	return result, nil
}

func (cl *ClientConnection) UpdateUser(id *primitive.ObjectID, user User) (*mongo.UpdateResult, error) {
	//	password.Hash(user.Password)
	time := time.Now().Format("2006.01.02 15:04")
	update := bson.D{{"$set", bson.D{
		{"Nickname", user.Nickname},
		{"FirstName", user.FirstName},
		{"LastName", user.LastName},
		{"Password", user.Password},
		{"updated_at", time},
	}}}
	result, err := cl.collection.UpdateByID(context.Background(), id, update)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t update user`s data")
		return nil, err
	}
	return result, nil

}

func (cl *ClientConnection) DeleteUser(id *primitive.ObjectID) (*mongo.DeleteResult, error) {
	result, err := cl.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		log.Warn().Err(err).Msg(" can`t delete user`s data")
		return &mongo.DeleteResult{}, err
	}

	return result, nil
}
