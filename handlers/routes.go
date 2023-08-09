package handlers

import (
	"github.com/labstack/echo/v4"
)

var userHandler = NewUserHandler(&RT)

func UserRoute(e *echo.Echo) {

	e.POST("/user", CreateUser)
	e.GET("/user/:userName", GetUser)
	e.PUT("/user/:userName", userHandler.dbUser.EditUser)
	e.PUT("/soft_user_delete/:userName", userHandler.dbUser.SoftDeleteUser)
	e.DELETE("/protected/user/:userName", userHandler.dbUser.DeleteUser)
	e.GET("/users", GetAllUsers)
	e.GET("/protected", ProtectedHandler)
}
