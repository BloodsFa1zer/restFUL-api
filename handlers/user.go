package handlers

import (
	"app3.1/database"
	"app3.1/response"
	"crypto/subtle"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserHandlersInterface interface {
	EditUser(c echo.Context) error
	SoftDeleteUser(c echo.Context) error
	DeleteUser(c echo.Context) error
	CreateUser(c echo.Context) error
	GetUser(c echo.Context) error
	GetAllUsers(c echo.Context) error
	ProtectedHandler(c echo.Context) error
	database.UserDatabase
}

type UserHandler struct {
	validate *validator.Validate
	*database.UserDatabase
	isRootEnabled bool
}

func NewUserHandler(isRootEnabled bool) *UserHandler {
	return &UserHandler{validate: validator.New(), UserDatabase: database.NewDatabase(), isRootEnabled: isRootEnabled}
}

func (rt *UserHandler) CreateUser(c echo.Context) error {
	var user database.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	if validationErr := rt.validate.Struct(&user); validationErr != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": validationErr.Error()}})
	}

	newUser := database.User{
		Nickname:  user.Nickname,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
	}

	result, err := rt.InsertUser(newUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusCreated, response.UserResponse{Status: http.StatusCreated, Message: "success", Data: &echo.Map{"data to interact with profile": result}})
}

func (rt *UserHandler) GetUser(c echo.Context) error {
	userName := c.Param("userName")

	user, err := rt.FindByUserName(userName)
	if err == errors.New("user not found") {
		return c.JSON(http.StatusNotFound, response.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": err.Error()}})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": user}})
}

func (rt *UserHandler) EditUser(c echo.Context) error {
	if rt.GivePerm(c) {
		userName := c.Param("userName")
		if rt.IsUserDeleted(userName) == true {
			var user database.User

			if err := c.Bind(&user); err != nil {
				return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
			}

			if validationErr := rt.validate.Struct(&user); validationErr != nil {
				return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": validationErr.Error()}})
			}

			updatedUserName, err := rt.UpdateUser(userName, user)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
			}

			updatedUser, _ := rt.FindByUserName(updatedUserName)

			return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": updatedUser}})
		} else {
			return c.JSON(http.StatusConflict, response.UserResponse{Status: http.StatusConflict, Message: "error", Data: &echo.Map{"data": "This user is already deleted"}})
		}
	} else {
		return c.JSON(http.StatusLocked, response.UserResponse{Status: http.StatusLocked, Message: "error", Data: &echo.Map{"data": "Only RootUsers can do that"}})
	}
}

func (rt *UserHandler) GetAllUsers(c echo.Context) error {
	users, err := rt.FindUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": users}})
}

func (rt *UserHandler) SoftDeleteUser(c echo.Context) error {
	if rt.GivePerm(c) {
		userName := c.Param("userName")
		if rt.IsUserDeleted(userName) == true {
			deletedUserName, err := rt.SoftDelete(userName)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
			}
			updatedUser, _ := rt.FindByUserName(deletedUserName)

			return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": updatedUser}})
		} else {
			return c.JSON(http.StatusConflict, response.UserResponse{Status: http.StatusConflict, Message: "error", Data: &echo.Map{"data": "This user is already deleted"}})
		}
	} else {
		return c.JSON(http.StatusLocked, response.UserResponse{Status: http.StatusLocked, Message: "error", Data: &echo.Map{"data": "Only RootUsers can do that"}})
	}
}

func (rt *UserHandler) DeleteUser(c echo.Context) error {
	if rt.GivePerm(c) {
		userName := c.Param("userName")
		if rt.IsUserDeleted(userName) == true {
			err := rt.DeleteUserByNick(userName)
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

func (rt *UserHandler) GivePerm(c echo.Context) bool {
	rUser := UserHandler{}

	rt.ProtectedHandler(c)
	if c.Response().Status == 200 {
		rUser.isRootEnabled = true
		return rUser.isRootEnabled
	}
	rUser.isRootEnabled = false
	return rUser.isRootEnabled
}

func (rt *UserHandler) ProtectedHandler(c echo.Context) error {
	username, password, ok := c.Request().BasicAuth()

	if !ok {
		return c.JSON(http.StatusUnauthorized, response.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &echo.Map{"data": "Credentials are not valid"}})

	}

	if !IsValidCredentials(username, password) {
		return c.JSON(http.StatusUnauthorized, response.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &echo.Map{"data": "Credentials are not valid"}})
	}
	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": "You are in the protected mode"}})

}
