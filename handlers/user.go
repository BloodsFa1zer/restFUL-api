package handlers

import (
	"app3.1/database"
	"app3.1/response"
	"crypto/subtle"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserActions interface {
	EditUser(c echo.Context) error
	SoftDeleteUser(c echo.Context) error
	DeleteUser(c echo.Context) error
}

var RT = RootUser{}

type RootUser struct {
	isRootEnabled bool
}

var db = database.Connection()
var validate = validator.New()

type UserHandler struct {
	dbUser UserActions
}

func NewUserHandler(User UserActions) *UserHandler {
	return &UserHandler{dbUser: User}
}

func CreateUser(c echo.Context) error {
	var user database.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": validationErr.Error()}})
	}

	newUser := database.User{
		Nickname:  user.Nickname,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
	}

	result, err := db.InsertUser(newUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusCreated, response.UserResponse{Status: http.StatusCreated, Message: "success", Data: &echo.Map{"data to interact with profile": result}})
}

func GetUser(c echo.Context) error {
	userName := c.Param("userName")

	user, err := db.FindUser(userName)
	if err == errors.New("user not found") {
		return c.JSON(http.StatusNotFound, response.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": err.Error()}})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": user}})
}

func (rt *RootUser) EditUser(c echo.Context) error {

	if rt.GivePerm(c) {
		userName := c.Param("userName")
		if db.IsUserDeleted(userName) == true {
			var user database.User

			if err := c.Bind(&user); err != nil {
				return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
			}

			if validationErr := validate.Struct(&user); validationErr != nil {
				return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": validationErr.Error()}})
			}

			updatedUserName, err := db.UpdateUser(userName, user)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
			}

			updatedUser, _ := db.FindUser(updatedUserName)

			return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": updatedUser}})
		} else {
			return c.JSON(http.StatusConflict, response.UserResponse{Status: http.StatusConflict, Message: "error", Data: &echo.Map{"data": "This user is already deleted"}})
		}
	} else {
		return c.JSON(http.StatusLocked, response.UserResponse{Status: http.StatusLocked, Message: "error", Data: &echo.Map{"data": "Only RootUsers can do that"}})
	}
}

func GetAllUsers(c echo.Context) error {
	users, err := db.FindUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": users}})
}

func (rt *RootUser) SoftDeleteUser(c echo.Context) error {
	if rt.GivePerm(c) {
		userName := c.Param("userName")
		if db.IsUserDeleted(userName) == true {
			deletedUserName, err := db.SoftDeleteUser(userName)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
			}
			updatedUser, _ := db.FindUser(deletedUserName)

			return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": updatedUser}})
		} else {
			return c.JSON(http.StatusConflict, response.UserResponse{Status: http.StatusConflict, Message: "error", Data: &echo.Map{"data": "This user is already deleted"}})
		}
	} else {
		return c.JSON(http.StatusLocked, response.UserResponse{Status: http.StatusLocked, Message: "error", Data: &echo.Map{"data": "Only RootUsers can do that"}})
	}
}

func (rt *RootUser) DeleteUser(c echo.Context) error {

	if rt.GivePerm(c) {
		userName := c.Param("userName")
		if db.IsUserDeleted(userName) == true {
			err := db.DeleteUser(userName)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
			}

			return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": "user successfully deleted"}})
		} else {
			return c.JSON(http.StatusConflict, response.UserResponse{Status: http.StatusConflict, Message: "error", Data: &echo.Map{"data": "This user is already deleted"}})
		}

	} else {
		return c.JSON(http.StatusLocked, response.UserResponse{Status: http.StatusLocked, Message: "error", Data: &echo.Map{"data": "Only RootUsers can do that"}})
	}

}

func IsValidCredentials(username, password string) bool {
	if subtle.ConstantTimeCompare([]byte(username), []byte("user")) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), []byte("password")) == 1 {
		return true
	}
	return false
}

func (rt *RootUser) GivePerm(c echo.Context) bool {

	rUser := RootUser{}

	ProtectedHandler(c)
	if c.Response().Status == 200 {
		rUser.isRootEnabled = true
		fmt.Println("Success")
		return rUser.isRootEnabled
	}
	rUser.isRootEnabled = false
	return rUser.isRootEnabled
}

func ProtectedHandler(c echo.Context) error {
	username, password, ok := c.Request().BasicAuth()

	if !ok {
		return c.JSON(http.StatusUnauthorized, response.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &echo.Map{"data": "Credentials are not valid"}})

	}

	if !IsValidCredentials(username, password) {
		return c.JSON(http.StatusUnauthorized, response.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &echo.Map{"data": "Credentials are not valid"}})
	}
	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": "You are in the protected mode"}})

}
