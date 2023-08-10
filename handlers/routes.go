package handlers

import (
	"github.com/labstack/echo/v4"
)

var userHandler = NewUserHandler(rootUser, false)

func UserRoute(e *echo.Echo) {

	e.POST("/user", userHandler.dbUser.CreateUser)
	e.GET("/user/:userName", userHandler.dbUser.GetUser)
	e.PUT("/user/:userName", userHandler.dbUser.EditUser)
	e.PUT("/soft_user_delete/:userName", userHandler.dbUser.SoftDeleteUser)
	e.DELETE("/protected/user/:userName", userHandler.dbUser.DeleteUser)
	e.GET("/users", userHandler.dbUser.GetAllUsers)
	e.GET("/protected", userHandler.dbUser.ProtectedHandler)
}
