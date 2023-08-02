package server

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func UserRoute(e *echo.Echo) {

}

func save(c echo.Context) error {
	// Get name and email
	name := c.FormValue("name")
	email := c.FormValue("email")
	return c.String(http.StatusOK, "name:"+name+", email:"+email)
}
