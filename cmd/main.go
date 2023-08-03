package main

import (
	"app3.1/server"
	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()

	server.UserRoute(e)

	e.Logger.Fatal(e.Start(":6000"))

}
