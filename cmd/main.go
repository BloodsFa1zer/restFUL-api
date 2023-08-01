package main

import (
	"app3.1/ENV"
	"app3.1/connection"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	cfg := ENV.LoadENV("ENV/.env")
	cfg.ParseENV()

	con := connection.NewConnection(*cfg)
	fmt.Println(con)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))

	e.POST("/save", save)

}

func save(c echo.Context) error {
	// Get name and email
	name := c.FormValue("name")
	email := c.FormValue("email")
	return c.String(http.StatusOK, "name:"+name+", email:"+email)
}
