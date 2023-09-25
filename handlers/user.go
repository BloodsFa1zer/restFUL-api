package handlers

import (
	"app3.1/database"
	"app3.1/response"
	"app3.1/service"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type UserHandler struct {
	userService service.UserServiceInterface
}

func NewUserHandler(service service.UserServiceInterface) *UserHandler {
	return &UserHandler{userService: service}
}

const (
	adminRole     = "Admin"
	userRole      = "User"
	moderatorRole = "Moderator"
)

func (uh *UserHandler) CreateUser(c echo.Context) error {
	var user database.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	userID, err, respStatus := uh.userService.CreateUser(user)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "success", Data: &echo.Map{"use this ID to interact with user`s profile": userID}})
}

func (uh *UserHandler) GetUser(c echo.Context) error {
	ID := c.Param("id")
	userID, err := strconv.Atoi(ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	user, err, respStatus := uh.userService.GetUser(int64(userID))
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "success", Data: &echo.Map{"data": user}})
}

func (uh *UserHandler) EditUser(c echo.Context) error {
	if permission, respStatus := uh.isUserHavePermissionToActions(adminRole, c); !permission {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": "that user has no access to admin actions"}})
	}

	ID := c.Param("id")
	userID, err := strconv.Atoi(ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	var user database.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	updatedUserID, err, respStatus := uh.userService.EditUser(int64(userID), user)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	updatedUser, err, respStatus := uh.userService.GetUser(updatedUserID)
	return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "success", Data: &echo.Map{"data": updatedUser}})
}

func (uh *UserHandler) GetAllUsers(c echo.Context) error {
	users, err, respStatus := uh.userService.GetAllUsers()
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "success", Data: &echo.Map{"data": users}})
}

func (uh *UserHandler) DeleteUser(c echo.Context) error {
	if permission, respStatus := uh.isUserHavePermissionToActions(adminRole, c); !permission {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": "that user has no access to admin actions"}})
	}

	ID := c.Param("id")
	userID, err := strconv.Atoi(ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	err, respStatus := uh.userService.DeleteUser(int64(userID))
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "success", Data: &echo.Map{"data": "User successfully deleted"}})
}

func (uh *UserHandler) Login(c echo.Context) error {
	var user database.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	t, err, respStatus := uh.userService.CreateToken(user)
	if err == errors.New("you have no account and will be redirected to registration page") {
		//return c.Redirect(respStatus, "/register")
	}
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "success", Data: &echo.Map{"token": t}})
}

// to Vote user need to register or login and then using his Bearer Token to vote, by POST(/user/:id)
func (uh *UserHandler) PostVote(c echo.Context) error {
	ID := c.Param("id")
	userID, err := strconv.Atoi(ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	user := c.Get("user")

	userName := uh.userService.GetUserNameViaToken(user)
	if userName == "" {
		return c.JSON(http.StatusUnauthorized, response.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &echo.Map{"data": "cannot find that user"}})
	}

	err, respStatus := uh.userService.PostVote(int64(userID), userName)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "success", Data: &echo.Map{"data": "you are successfully voted"}})
}

func (uh *UserHandler) DeleteVote(c echo.Context) error {
	ID := c.Param("id")
	userID, err := strconv.Atoi(ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	user := c.Get("user")

	userName := uh.userService.GetUserNameViaToken(user)
	if userName == "" {
		return c.JSON(http.StatusUnauthorized, response.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &echo.Map{"data": "cannot find that user"}})
	}

	err, respStatus := uh.userService.DeleteVote(int64(userID), userName)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "success", Data: &echo.Map{"data": "you are successfully voted"}})
}

func (uh *UserHandler) GetUserRate(c echo.Context) error {
	ID := c.Param("id")
	userID, err := strconv.Atoi(ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	userRate, err, respStatus := uh.userService.GetUserRate(int64(userID))
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "success", Data: &echo.Map{"data": userRate}})
}

func (uh *UserHandler) GetUsersRate(c echo.Context) error {

	userRate, err, respStatus := uh.userService.GetAllUsersRate()
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "success", Data: &echo.Map{"data": userRate}})
}

func (uh *UserHandler) GetUserRateModerator(c echo.Context) error {
	if permission, respStatus := uh.isUserHavePermissionToActions(userRole, c); permission {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": "that user has no access to admin or moderator actions"}})
	}

	ID := c.Param("id")
	userID, err := strconv.Atoi(ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	userRate, err, respStatus := uh.userService.GetUserRateModerator(int64(userID))
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "success", Data: &echo.Map{"data": userRate}})
}

func (uh *UserHandler) isUserHavePermissionToActions(roleToFind string, c echo.Context) (bool, int) {
	user := c.Get("user")

	return uh.userService.IsUserHavePermission(roleToFind, user)
}
