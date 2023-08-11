package handlers

import (
	"app3.1/database"
	"github.com/labstack/echo/v4"
)

var userHandler = NewUserHandler(&database.UserDatabase{Connection: database.NewDatabase()}, false)

func UserRoute(e *echo.Echo) {

	e.POST("/user", userHandler.CreateUser)
	e.GET("/user/:userName", userHandler.GetUser)
	e.PUT("/user/:userName", userHandler.EditUser)
	e.DELETE("/user/:userName", userHandler.DeleteUser)
	e.GET("/users", userHandler.GetAllUsers)
	e.GET("/protected", userHandler.ProtectedHandler)
}
