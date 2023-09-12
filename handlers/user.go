package handlers

import (
	"app3.1/database"
	"app3.1/response"
	"app3.1/serviceLayer"
	"database/sql"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type UserHandler struct {
	service serviceLayer.UserServiceInterface
}

func NewUserHandler(service serviceLayer.UserServiceInterface) *UserHandler {
	return &UserHandler{service: service}
}

const (
	adminRole = "Admin"
	// userRole      = "User"
	// moderatorRole = "Moderator"
)

func (uh *UserHandler) CreateUser(c echo.Context) error {
	var user database.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	userID, err := uh.service.CreateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusCreated, response.UserResponse{Status: http.StatusCreated, Message: "success", Data: &echo.Map{"use this ID to interact with user`s profile": userID}})
}

func (uh *UserHandler) GetUser(c echo.Context) error {
	userID, err := enterParameter(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	user, err := uh.service.GetUser(userID)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusConflict, response.UserResponse{Status: http.StatusConflict, Message: "error", Data: &echo.Map{"data": "There is no user with that ID"}})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": user}})
}

func (uh *UserHandler) EditUser(c echo.Context) error {
	if !uh.isUserHavePermissionToActions(adminRole, c) {
		return c.JSON(http.StatusUnauthorized, response.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &echo.Map{"data": "that user has no access to admin actions"}})
	}

	userID, err := enterParameter(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	var user database.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	updatedUserID, err := uh.service.EditUser(userID, user)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusConflict, response.UserResponse{Status: http.StatusConflict, Message: "error", Data: &echo.Map{"data": "There is no user with that ID"}})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	updatedUser, err := uh.service.GetUser(updatedUserID)
	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": updatedUser}})
}

func (uh *UserHandler) GetAllUsers(c echo.Context) error {
	users, err := uh.service.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": users}})
}

func (uh *UserHandler) DeleteUser(c echo.Context) error {
	if !uh.isUserHavePermissionToActions(adminRole, c) {
		return c.JSON(http.StatusUnauthorized, response.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &echo.Map{"data": "that user has no access to admin actions"}})
	}

	userID, err := enterParameter(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	err = uh.service.DeleteUser(userID)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusConflict, response.UserResponse{Status: http.StatusConflict, Message: "error", Data: &echo.Map{"data": "There is no user with that ID"}})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": "User successfully deleted"}})
}

func enterParameter(c echo.Context) (int64, error) {
	ID := c.Param("id")
	userID, err := strconv.Atoi(ID)
	return int64(userID), err
}

func (uh *UserHandler) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	t, err := uh.service.GetUserByName(username, password)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusConflict, response.UserResponse{Status: http.StatusConflict, Message: "error", Data: &echo.Map{"data": "There is no user with that NickName"}})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"token": t}})
}

func (uh *UserHandler) isUserHavePermissionToActions(roleToFind string, c echo.Context) bool {
	user := c.Get("user")

	return uh.service.IsUserHavePermission(roleToFind, user)
}
