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
	userID, err, respStatus := enterParameter(c)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	user, err, respStatus := uh.userService.GetUser(userID)
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
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	updatedUserID, err, respStatus := uh.userService.EditUser(userID, user)
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

	userID, err, respStatus := enterParameter(c)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	err, respStatus = uh.userService.DeleteUser(userID)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "success", Data: &echo.Map{"data": "User successfully deleted"}})
}

func (uh *UserHandler) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	t, err, respStatus := uh.userService.CreateToken(username, password)
	if err == errors.New("you have no account and will be redirected to registration page") {
		return c.Redirect(respStatus, "/register")
	}
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "success", Data: &echo.Map{"token": t}})
}

func (uh *UserHandler) UserRegistration(c echo.Context) error {
	username := c.FormValue("username")
	firstName := c.FormValue("firstName")
	surName := c.FormValue("lastName")
	password := c.FormValue("password")

	userID, err, respStatus := uh.userService.Registration(username, firstName, surName, password)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	str := "Now you need to login, to get access to actions that is only for registered users." +
		"ID to interact with your profile is:"

	return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{str: userID}})
}

// to Vote user need to register or login and then using his Bearer Token to vote, by POST(/user/:id)
func (uh *UserHandler) Voting(c echo.Context) error {

	userID, err, respStatus := enterParameter(c)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	userName := uh.GiveUserToken(c)

	err, respStatus = uh.userService.Vote(userID, userName)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": "you are successfully voted"}})
}

func (uh *UserHandler) GetUserRate(c echo.Context) error {
	userID, err, respStatus := enterParameter(c)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	user, err, respStatus := uh.userService.GetUserRate(userID)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "success", Data: &echo.Map{"data": user}})
}

func enterParameter(c echo.Context) (int64, error, int) {
	ID := c.Param("id")
	userID, err := strconv.Atoi(ID)
	return int64(userID), err, http.StatusNotFound
}

func (uh *UserHandler) isUserHavePermissionToActions(roleToFind string, c echo.Context) (bool, int) {
	user := c.Get("user")

	return uh.userService.IsUserHavePermission(roleToFind, user)
}

func (uh *UserHandler) GiveUserToken(c echo.Context) string {
	user := c.Get("user")

	return uh.userService.GetUserNameViaToken(user)
}
