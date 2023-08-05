package routes

import (
	"app3.1/handlers"
	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Echo) {
	e.POST("/user", handlers.CreateUser)
	e.GET("/user/:userName", handlers.GetUser)
	e.PUT("/user/:userName", handlers.EditUser)
	e.DELETE("/soft_user_delete/:userName", handlers.SoftDeleteUser)
	e.DELETE("/user/:userName", handlers.DeleteUser)
	e.GET("/users", handlers.GetAllUsers)
}
