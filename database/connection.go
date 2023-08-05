package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

//type ClientConnection struct {
//	collection *mongo.Collection
//}

type Database struct {
	connection *sql.DB
}

func Connection() *Database {
	db, err := sql.Open("sqlite3", "database/test.db")
	if err != nil {
		log.Warn().Err(err).Msg(" can`t connect to SQLite")
		return nil
	}
	database := Database{
		connection: db,
	}
	return &database
}

//func NewConnection() *mongo.Client {
//	cfg := config.LoadENV("config/.env")
//	ctx := context.TODO()
//	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoENV()))
//	if err != nil {
//		log.Warn().Err(err).Msg(" can`t connect to MongoDB")
//		return nil
//	}
//	log.Info().Msg("successfully connected to MongoDB")
//
//	err = client.Ping(ctx, nil)
//	if err != nil {
//		log.Fatal().Err(err).Msg(" can`t ping MongoDB")
//	}
//	fmt.Println("Connected to MongoDB")
//
//	return client
//}
//
//var DB = NewConnection()
//
//func GetCollection() *ClientConnection {
//
//	clientConn := ClientConnection{
//		collection: DB.Database("WebAPI").Collection("Users"),
//	}
//	return &clientConn
//}
