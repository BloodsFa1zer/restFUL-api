package routes

import (
	"app3.1/database"
	"app3.1/handlers"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var validate = validator.New()
var userHandler = handlers.NewUserHandler(database.NewUserDatabase(), validate)

func UserRoute(e *echo.Echo) {
	e.POST("/user", userHandler.CreateUser)
	e.GET("/user/:id", userHandler.GetUser)
	e.PUT("/user/:id", userHandler.EditUser)
	//e.PUT("/user/:id", userHandler.EditUser, Middleware.IsLoggedIn)
	e.DELETE("/user/:id", userHandler.DeleteUser)
	//e.DELETE("/user/:id", userHandler.DeleteUser, Middleware.IsLoggedIn)
	e.GET("/users", userHandler.GetAllUsers)

}
