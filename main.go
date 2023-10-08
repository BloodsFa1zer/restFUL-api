package main

import (
	"app3.1/Middleware"
	"app3.1/routes"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.HidePort = true

	Middleware.UserAuth(e)
	routes.UserRoute(e)

	e.Logger.Fatal(e.Start(":6000"))
}
