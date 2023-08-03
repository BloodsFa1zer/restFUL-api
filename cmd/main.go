package main

import (
	"app3.1/server"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	Nickname  string             `bson:"Nickname,omitempty" json:"Nickname"`
	FirstName string             `bson:"FirstName,omitempty" json:"FirstName"`
	LastName  string             `bson:"LastName, omitempty" json:"LastName"`
	Password  string             `bson:"Password, omitempty" json:"Password"`
}

func main() {

	e := echo.New()

	server.UserRoute(e)

	e.Logger.Fatal(e.Start(":6000"))

}
