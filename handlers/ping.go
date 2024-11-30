package handlers

import (
	"checkin/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, utils.JSON{"confirmation": "hello"})
}
