package Middleware

import (
	"app3.1/config"
	"app3.1/database"
	"app3.1/hash"
	"app3.1/response"
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

//func IsValidCredentials(username, password string, c echo.Context) (bool, error) {
//	if subtle.ConstantTimeCompare([]byte(username), []byte("user")) == 1 &&
//		subtle.ConstantTimeCompare([]byte(password), []byte("password")) == 1 {
//		return true, nil
//	}
//	return false, nil
//}

var IsLoggedIn = echojwt.JWT(middleware.JWTConfig{
	SigningKey: []byte("secret"),
})

type UserHandler struct {
	DbUser database.DbInterface
}

func NewUserHandler(dbUser database.DbInterface) *UserHandler {
	return &UserHandler{DbUser: dbUser}
}

func (uh *UserHandler) Login(c echo.Context) error {

	loginRequest := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := c.Bind(&loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	//fmt.Println(loginRequest.Username)
	//fmt.Println(loginRequest.Password)

	//Username := c.FormValue("username")
	//Password := c.FormValue("password")

	cfg := config.LoadENV("config/.env")
	cfg.ParseENV()

	hashedPass, err := uh.DbUser.FindByNickname(loginRequest.Username)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusConflict, response.UserResponse{Status: http.StatusConflict, Message: "error", Data: &echo.Map{"data": "There is no user with that password"}})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	isPasswordCorrect := hash.Verify(hashedPass, loginRequest.Password)
	if isPasswordCorrect == false {
		return echo.ErrUnauthorized
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = loginRequest.Username
	claims["exp"] = time.Now().Add(48 * time.Hour).Unix()

	t, err := token.SignedString([]byte(cfg.SigningKey))
	if err != nil {
		log.Warn().Err(err)
		return c.JSON(http.StatusUnauthorized, response.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"token": t}})
}
