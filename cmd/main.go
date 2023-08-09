package main

import (
	"app3.1/handlers"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	handlers.UserRoute(e)
	e.Logger.Fatal(e.Start(":6000"))
}
