package Middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func UserAuth(e *echo.Echo) {

	e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

}
