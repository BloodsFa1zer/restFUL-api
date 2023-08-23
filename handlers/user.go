package handlers

import (
	"app3.1/database"
	"app3.1/response"
	"crypto/subtle"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type UserHandlersInterface interface {
	EditUser(c echo.Context) error
	DeleteUser(c echo.Context) error
	CreateUser(c echo.Context) error
	GetUser(c echo.Context) error
	GetAllUsers(c echo.Context) error
	EnterRootMode(c echo.Context) error
}

var validate = validator.New()

type UserHandler struct {
	DbUser        database.DbInterface
	isRootEnabled bool
}

func NewUserHandler(dbUser database.DbInterface, isRootEnabled bool) *UserHandler {
	return &UserHandler{DbUser: dbUser, isRootEnabled: isRootEnabled}

}

func (uh *UserHandler) CreateUser(c echo.Context) error {
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

	result, err := uh.DbUser.InsertUser(newUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusCreated, response.UserResponse{Status: http.StatusCreated, Message: "success", Data: &echo.Map{"data to interact with profile": result}})
}

func (uh *UserHandler) GetUser(c echo.Context) error {
	ID := c.Param("id")
	input, err := strconv.Atoi(ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	user, err := uh.DbUser.FindByID(input)
	if err == errors.New("user not found") {
		return c.JSON(http.StatusNotFound, response.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": err.Error()}})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": user}})
}

func (uh *UserHandler) EditUser(c echo.Context) error {
	if uh.GivePerm(c) {
		ID := c.Param("id")
		input, err := strconv.Atoi(ID)
		if err != nil {
			return c.JSON(http.StatusNotFound, response.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": err.Error()}})
		}

		if uh.DbUser.IsUserDeleted(input) == true {
			var user database.User

			if err := c.Bind(&user); err != nil {
				return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
			}

			if validationErr := validate.Struct(&user); validationErr != nil {
				return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": validationErr.Error()}})
			}

			updatedUserID, err := uh.DbUser.UpdateUser(input, user)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
			}

			updatedUser, err := uh.DbUser.FindByID(updatedUserID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
			}

			return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": updatedUser}})
		} else {
			return c.JSON(http.StatusConflict, response.UserResponse{Status: http.StatusConflict, Message: "error", Data: &echo.Map{"data": "This user is already deleted"}})
		}
	} else {
		return c.JSON(http.StatusLocked, response.UserResponse{Status: http.StatusLocked, Message: "error", Data: &echo.Map{"data": "Only RootUsers can do that"}})
	}
}

func (uh *UserHandler) GetAllUsers(c echo.Context) error {
	users, err := uh.DbUser.FindUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": users}})
}

func (uh *UserHandler) DeleteUser(c echo.Context) error {
	if uh.GivePerm(c) {
		ID := c.Param("id")
		input, err := strconv.Atoi(ID)
		if err != nil {
			return c.JSON(http.StatusNotFound, response.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": err.Error()}})
		}

		if uh.DbUser.IsUserDeleted(input) == true {
			deletedUserName, err := uh.DbUser.DeleteUserByID(input)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
			}
			updatedUser, _ := uh.DbUser.FindByID(deletedUserName)

			return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": updatedUser}})
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

func (uh *UserHandler) GivePerm(c echo.Context) bool {
	rUser := UserHandler{}

	err := uh.EnterRootMode(c)
	if err != nil {
		rUser.isRootEnabled = true
		return rUser.isRootEnabled
	}

	if c.Response().Status == 200 {
		rUser.isRootEnabled = true
		return rUser.isRootEnabled
	}
	rUser.isRootEnabled = false
	return rUser.isRootEnabled
}

func (uh *UserHandler) EnterRootMode(c echo.Context) error {
	username, password, ok := c.Request().BasicAuth()

	if !ok {
		return c.JSON(http.StatusUnauthorized, response.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &echo.Map{"data": "Credentials are not valid"}})

	}

	if !IsValidCredentials(username, password) {
		return c.JSON(http.StatusUnauthorized, response.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &echo.Map{"data": "Credentials are not valid"}})
	}
	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": "You are in the protected mode"}})

}
