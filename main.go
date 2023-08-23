package main

import (
	"app3.1/routes"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	routes.UserRoute(e)
	e.Logger.Fatal(e.Start(":6000"))
}
