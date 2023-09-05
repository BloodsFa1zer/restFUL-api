package handlers

import (
	"app3.1/config"
	"app3.1/database"
	"app3.1/hash"
	"app3.1/response"
	"database/sql"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

type UserHandlersInterface interface {
	EditUser(c echo.Context) error
	DeleteUser(c echo.Context) error
	CreateUser(c echo.Context) error
	GetUser(c echo.Context) error
	GetAllUsers(c echo.Context) error
	Login(c echo.Context) error
	Permission(c echo.Context) error
}

type UserHandler struct {
	DbUser   database.DbInterface
	validate *validator.Validate
}

func NewUserHandler(dbUser database.DbInterface, validate *validator.Validate) *UserHandler {
	return &UserHandler{DbUser: dbUser, validate: validate}
}

func (uh *UserHandler) CreateUser(c echo.Context) error {
	var user database.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	if validationErr := uh.validate.Struct(&user); validationErr != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": validationErr.Error()}})
	}

	newUser := database.User{
		Nickname:  user.Nickname,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
	}

	userID, err := uh.DbUser.InsertUser(newUser)
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

	user, err := uh.DbUser.FindByID(userID)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusConflict, response.UserResponse{Status: http.StatusConflict, Message: "error", Data: &echo.Map{"data": "There is no user with that ID"}})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": user}})
}

func (uh *UserHandler) EditUser(c echo.Context) error {
	userID, err := enterParameter(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	var user database.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	if validationErr := uh.validate.Struct(&user); validationErr != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": validationErr.Error()}})
	}

	updatedUserID, err := uh.DbUser.UpdateUser(userID, user)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusConflict, response.UserResponse{Status: http.StatusConflict, Message: "error", Data: &echo.Map{"data": "There is no user with that ID"}})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	updatedUser, err := uh.DbUser.FindByID(updatedUserID)
	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": updatedUser}})
}

func (uh *UserHandler) GetAllUsers(c echo.Context) error {
	users, err := uh.DbUser.FindUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": users}})
}

func (uh *UserHandler) DeleteUser(c echo.Context) error {
	userID, err := enterParameter(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	err = uh.DbUser.DeleteUserByID(userID)
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
	cfg := config.LoadENV("config/.env")
	cfg.ParseENV()

	username := c.FormValue("username")
	password := c.FormValue("password")

	selectedUser, err := uh.DbUser.FindByNickname(username)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusConflict, response.UserResponse{Status: http.StatusConflict, Message: "error", Data: &echo.Map{"data": "There is no user with that ID"}})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	if username != selectedUser.Nickname || hash.Verify(selectedUser.Password, password) != true {
		return echo.ErrUnauthorized
	}

	claims := &config.JwtCustomClaims{
		Name: selectedUser.Nickname,
		Role: selectedUser.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(cfg.SigningKey))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"token": t}})
}

//func (uh *UserHandler) Permission(c echo.Context) error {
//	user := c.Get("user")
//	userToken, ok := user.(jwt.Token)
//	if !ok {
//		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": "token is nil"}})
//	}
//	claims := userToken.Claims.(*config.JwtCustomClaims)
//	log.Info().Type("claims", claims)
//	// log.InfoDump(claims, "claims")
//	expirationTime, err := claims.GetExpirationTime()
//	if err != nil {
//		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
//	}
//	name := claims.Name
//
//	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"message": "Welcome " + name + "! Your permission will continue for: " + expirationTime.Format("2006.01.02 15:04")}})
//}
