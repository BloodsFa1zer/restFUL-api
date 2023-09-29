package handlers

import (
	"github.com/labstack/echo/v4"
)

type UserHandlersInterface interface {
	EditUser(c echo.Context) error
	DeleteUser(c echo.Context) error
	CreateUser(c echo.Context) error
	GetUser(c echo.Context) error
	GetAllUsers(c echo.Context) error
	Login(c echo.Context) error
	isUserHavePermissionToActions(roleToFind string, c echo.Context) (bool, int)
	// GetUsersRate(c echo.Context) error
	// GetUserRate(c echo.Context) error
	// GetUserRateModerator(c echo.Context) error
	PostVote(c echo.Context) error
	//DeleteVote(c echo.Context) error
}
