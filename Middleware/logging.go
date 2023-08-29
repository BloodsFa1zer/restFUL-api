package Middleware

import (
	"crypto/subtle"
	"github.com/labstack/echo/v4"
)

func IsValidCredentials(username, password string, c echo.Context) (bool, error) {
	if subtle.ConstantTimeCompare([]byte(username), []byte("user")) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), []byte("password")) == 1 {
		return true, nil
	}
	return false, nil
}
