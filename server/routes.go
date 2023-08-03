package server

import (
	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Echo) {
	e.POST("/user", CreateUser)
	e.GET("/user/:userId", GetUser)
	e.PUT("/user/:userId", EditUser)
	e.DELETE("/user/:userId", DeleteUser)
	e.GET("/users", GetAllUsers)
}
