package routes

import (
	"app3.1/config"
	"app3.1/database"
	"app3.1/handlers"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

var validate = validator.New()
var userHandler = handlers.NewUserHandler(database.NewUserDatabase(), validate)

func UserRoute(e *echo.Echo) {
	r := e.Group("/restricted")
	r.Use(echojwt.WithConfig(JWT()))
	r.GET("", userHandler.Restricted)
	r.PUT("/user/:id", userHandler.EditUser)
	r.DELETE("/user/:id", userHandler.DeleteUser)
	// r.GET("/user/:id", userHandler.GetUser)

	e.POST("/user", userHandler.CreateUser)
	e.GET("/user/:id", userHandler.GetUser)
	e.GET("/users", userHandler.GetAllUsers)
	e.POST("/login", userHandler.Login)
	e.GET("/", userHandler.Accessible)

}

// Where should i put this func? which directory or .go file
func JWT() echojwt.Config {
	cfg := config.LoadENV("config/.env")
	cfg.ParseENV()

	Config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(handlers.JwtCustomClaims)
		},
		SigningKey: []byte(cfg.SigningKey),
	}
	return Config
}
