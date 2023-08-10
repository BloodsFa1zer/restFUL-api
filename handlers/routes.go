package handlers

import (
	"github.com/labstack/echo/v4"
)

var userHandler = NewUserHandler(false)

func UserRoute(e *echo.Echo) {

	e.POST("/user", userHandler.CreateUser)
	e.GET("/user/:userName", userHandler.GetUser)
	e.PUT("/user/:userName", userHandler.EditUser)
	e.PUT("/soft_user_delete/:userName", userHandler.SoftDeleteUser)
	e.DELETE("/protected/user/:userName", userHandler.DeleteUser)
	e.GET("/users", userHandler.GetAllUsers)
	e.GET("/protected", userHandler.ProtectedHandler)
}
