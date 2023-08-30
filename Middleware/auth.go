package Middleware

import (
	"app3.1/database"
	"github.com/labstack/echo/v4"
)

var userHandler = NewUserHandler(database.NewUserDatabase())

func UserAuth(e *echo.Echo) {
	// e.Use(middleware.BasicAuth(IsValidCredentials))

	e.POST("/login", userHandler.Login)
}
