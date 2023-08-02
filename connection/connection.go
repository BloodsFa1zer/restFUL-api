package connection

import (
	"app3.1/ENV"
	"context"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ClientConnection struct {
	collection *mongo.Collection
}

func NewConnection() *mongo.Client {
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(ENV.MongoENV()))
	if err != nil {
		log.Warn().Err(err).Msg(" can`t connect to MongoDB")
		return nil
	}
	log.Info().Msg("success")

	return client
}

func GetCollection() *ClientConnection {
	client := NewConnection()
	clientConn := ClientConnection{
		collection: client.Database("WebAPI").Collection("Users"),
	}
	return &clientConn
}
