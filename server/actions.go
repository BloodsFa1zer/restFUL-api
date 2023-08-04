package server

import (
	"app3.1/connection"
	"app3.1/response"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var cl = connection.GetCollection()
var validate = validator.New()

func CreateUser(c echo.Context) error {
	var user connection.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": validationErr.Error()}})
	}

	newUser := connection.User{
		ID:        primitive.NewObjectID(),
		Nickname:  user.Nickname,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
	}
	_, FindErr := cl.FindUser("Nickname", user.Nickname)
	if FindErr != mongo.ErrNoDocuments {
		return c.JSON(http.StatusLocked, response.UserResponse{Status: http.StatusLocked, Message: "error", Data: &echo.Map{"data": "That Nickname is already taken"}})
	}

	result, err := cl.InsertUser(newUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusCreated, response.UserResponse{Status: http.StatusCreated, Message: "success", Data: &echo.Map{"data": result}})
}

func GetUser(c echo.Context) error {
	userId := c.Param("userId")

	objectID, _ := primitive.ObjectIDFromHex(userId)
	fmt.Println(objectID)

	user, err := cl.FindUser("_id", objectID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": user}})
}

func EditUser(c echo.Context) error {
	userId := c.Param("userId")
	var user connection.User

	objectID, _ := primitive.ObjectIDFromHex(userId)

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": validationErr.Error()}})
	}

	result, err := cl.UpdateUser(&objectID, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	if result.MatchedCount == 1 {
		updatedUser, err := cl.FindUser("_id", objectID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
		}
		return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": updatedUser}})
	}

	return err
}

func DeleteUser(c echo.Context) error {
	userId := c.Param("userId")

	objectID, _ := primitive.ObjectIDFromHex(userId)

	result, err := cl.DeleteUser(&objectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.JSON(http.StatusNotFound, response.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": "User with that ID is not found!"}})
	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": "User successfully deleted"}})
}

func GetAllUsers(c echo.Context) error {
	cl := connection.GetCollection()

	users, err := cl.FindUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": users}})
}

func CheckAuth(c echo.Context) error {

	return nil
}

func Verify(hashed, password string) bool {
	user, err := cl.FindUser("Password", password)
	if err != nil {
		return false
	}
	errPassword := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(user.Password))
	if errPassword == nil {
		return true
	}
	return false
}
