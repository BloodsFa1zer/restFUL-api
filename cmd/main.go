package main

import (
	"app3.1/ENV"
	"app3.1/connection"
	"app3.1/server"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	Nickname  string             `bson:"Nickname,omitempty" json:"Nickname"`
	FirstName string             `bson:"FirstName,omitempty" json:"FirstName"`
	LastName  string             `bson:"LastName, omitempty" json:"LastName"`
	Password  string             `bson:"Password, omitempty" json:"Password"`
}

func main() {
	ENV.LoadENV("ENV/.env")

	connection.NewConnection()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	server.UserRoute(e)

	e.Logger.Fatal(e.Start(":1323"))

}

func save(c echo.Context) error {
	// Get name and email
	name := c.FormValue("name")
	email := c.FormValue("email")
	return c.String(http.StatusOK, "name:"+name+", email:"+email)
}
