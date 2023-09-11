package handlers

import (
	"app3.1/services"
	"github.com/labstack/echo/v4"
)

type UserHandlersInterface interface {
	EditUser(c echo.Context) error
	DeleteUser(c echo.Context) error
	CreateUser(c echo.Context) error
	GetUser(c echo.Context) error
	GetAllUsers(c echo.Context) error
	Login(c echo.Context) error
	isUserHavePermissionToActions(roleToFind string, c echo.Context) bool
}

type UserHandler struct {
	service services.UserServiceInterface
}

// Which way is better in that case? NewUserHandler(service services.UserServiceInterface) or written one?
func NewUserHandler() *UserHandler {
	return &UserHandler{service: services.NewUserService()}
}
