package routes

import (
	"app3.1/database"
	"app3.1/handlers"
	"github.com/labstack/echo/v4"
)

var userHandler = handlers.NewUserHandler(database.NewUserDatabase())

func UserRoute(e *echo.Echo) {
	e.POST("/user", userHandler.CreateUser)
	e.GET("/user/:id", userHandler.GetUser)
	e.PUT("/user/:id", userHandler.EditUser)
	e.DELETE("/user/:id", userHandler.DeleteUser)
	e.GET("/users", userHandler.GetAllUsers)
}
