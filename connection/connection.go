package connection

import (
	"app3.1/config"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ClientConnection struct {
	collection *mongo.Collection
}

func NewConnection() *mongo.Client {
	cfg := config.LoadENV("config/.env")
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoENV()))
	if err != nil {
		log.Warn().Err(err).Msg(" can`t connect to MongoDB")
		return nil
	}
	log.Info().Msg("successfully connected to MongoDB")

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal().Err(err).Msg(" can`t ping MongoDB")
	}
	fmt.Println("Connected to MongoDB")

	return client
}

var DB = NewConnection()

func GetCollection() *ClientConnection {

	clientConn := ClientConnection{
		collection: DB.Database("WebAPI").Collection("Users"),
	}
	return &clientConn
}
