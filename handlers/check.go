package handlers

import (
	"checkin/db"
	"checkin/utils"
	"encoding/base64"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CheckinRequest struct {
	UUID string `json:"uuid"`
}

func Checkin(c echo.Context) error {
	req := new(CheckinRequest)

	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusBadRequest, utils.JSON{"error": "UUID Invalid"})
	}

	UUID := req.UUID
	UUIDBinary, err := base64.StdEncoding.DecodeString(UUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.JSON{"error": "UUID Invalid"})
	}
	err = db.Checkin(UUIDBinary)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.JSON{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, utils.JSON{"confirmation": "Checkin recorded"})
}
