package handlers

import (
	"app3.1/database"
	"app3.1/response"
	"app3.1/service"
	"github.com/labstack/echo/v4"
	"strconv"
)

type UserHandler struct {
	service service.UserServiceInterface
}

func NewUserHandler(service service.UserServiceInterface) *UserHandler {
	return &UserHandler{service: service}
}

const (
	adminRole     = "Admin"
	userRole      = "User"
	moderatorRole = "Moderator"
)

func (uh *UserHandler) CreateUser(c echo.Context) error {
	var user database.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(400, response.UserResponse{Status: 400, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	userID, err, respStatus := uh.service.CreateUser(user)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "success", Data: &echo.Map{"use this ID to interact with user`s profile": userID}})
}

func (uh *UserHandler) GetUser(c echo.Context) error {
	userID, err, respStatus := enterParameter(c)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	user, err, respStatus := uh.service.GetUser(userID)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "success", Data: &echo.Map{"data": user}})
}

func (uh *UserHandler) EditUser(c echo.Context) error {
	if permission, respStatus := uh.isUserHavePermissionToActions(adminRole, c); !permission {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": "that user has no access to admin actions"}})
	}

	userID, err, respStatus := enterParameter(c)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	var user database.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(400, response.UserResponse{Status: 400, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	updatedUserID, err, respStatus := uh.service.EditUser(userID, user)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	updatedUser, err, respStatus := uh.service.GetUser(updatedUserID)
	return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "success", Data: &echo.Map{"data": updatedUser}})
}

func (uh *UserHandler) GetAllUsers(c echo.Context) error {
	users, err, respStatus := uh.service.GetAllUsers()
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "success", Data: &echo.Map{"data": users}})
}

func (uh *UserHandler) DeleteUser(c echo.Context) error {
	if permission, respStatus := uh.isUserHavePermissionToActions(adminRole, c); !permission {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": "that user has no access to admin actions"}})
	}

	userID, err, respStatus := enterParameter(c)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	err, respStatus = uh.service.DeleteUser(userID)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "success", Data: &echo.Map{"data": "User successfully deleted"}})
}

func enterParameter(c echo.Context) (int64, error, int) {
	ID := c.Param("id")
	userID, err := strconv.Atoi(ID)
	return int64(userID), err, 404
}

func (uh *UserHandler) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	t, err, respStatus := uh.service.GetToken(username, password)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "success", Data: &echo.Map{"token": t}})
}

func (uh *UserHandler) isUserHavePermissionToActions(roleToFind string, c echo.Context) (bool, int) {
	user := c.Get("user")

	return uh.service.IsUserHavePermission(roleToFind, user)
}
