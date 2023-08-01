package connection

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Nickname  string             `bson:"Nickname,omitempty"`
	FirstName string             `bson:"FirstName,omitempty"`
	LastName  string             `bson:"LastName, omitempty"`
	Password  string             `bson:"Password, omitempty"`
}

func (cl *ClientConnection) findUser(field string, dataToFind any) *User {

	result := cl.collection.FindOne(context.TODO(), bson.M{field: dataToFind})

	// check for errors in the finding
	if result.Err() != nil {
		log.Warn().Err(result.Err()).Msg(" can`t find user")
	}

	// convert the cursor result to bson
	var user User
	// check for errors in the conversion
	if err := result.Decode(&user); err != mongo.ErrNoDocuments {
		log.Warn().Err(err).Msg(" no results to convert")
		return nil
	} else if err != nil {
		log.Warn().Err(err).Msg(" can`t convert results")
		return nil
	}
	return &user
}

func (cl *ClientConnection) createUser(user User) {
	userInfo := bson.D{{"Nickname", user.Nickname}, {"FirstName", user.FirstName}, {"LastName", user.LastName}}
	_, err := cl.collection.InsertOne(context.TODO(), userInfo)
	if err != nil {
		log.Warn().Err(err).Msg(" can`t insert user`s data into database")
		return
	}
	log.Info().Msg("successfully insert user`s data")
}

func (cl *ClientConnection) updateUser(id *primitive.ObjectID, user User) {

	update := bson.D{{"Nickname", user.Nickname}, {"FirstName", user.FirstName}, {"LastName", user.LastName}}
	_, err := cl.collection.UpdateByID(context.Background(), id, update)
	if err != nil {
		panic(err)
	}

}

func (cl *ClientConnection) deleteUser(id *primitive.ObjectID, user User) {
	result, err := cl.collection.DeleteOne(context.Background(), id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deleted %d documents. \n", result.DeletedCount)

}
