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

func NewConnection(config ENV.Config) *ClientConnection {
	ctx := context.TODO()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.URI_BD))
	if err != nil {
		log.Warn().Err(err).Msg(" can`t connect to MongoDB")
		return nil
	}
	clientConn := ClientConnection{
		collection: client.Database("telegram").Collection("usersID"),
	}
	return &clientConn
}
