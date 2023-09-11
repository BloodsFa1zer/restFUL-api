package routes

import (
	"app3.1/config"
	"app3.1/handlers"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

var userHandler = handlers.NewUserHandler()

func UserRoute(e *echo.Echo) {
	protected := e.Group("")
	protected.Use(echojwt.WithConfig(config.NewConfig()))
	protected.PUT("/user/:id", userHandler.EditUser)
	protected.DELETE("/user/:id", userHandler.DeleteUser)

	e.POST("/user", userHandler.CreateUser)
	e.GET("/user/:id", userHandler.GetUser)
	e.GET("/users", userHandler.GetAllUsers)
	e.POST("/login", userHandler.Login)

}
