/*
	All the request including UUID should be sent as base64 string, and then decoded in each side
*/

package handlers

import (
	"checkin/db"
	"checkin/utils"
	"encoding/base64"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Type declaration necessary for declaring the JSON datatype specifically
type (
	UUIDResponse struct {
		UUID string `json:"uuid"`
	}
	// Only the name is needed from the client to create a new card record
	CreateCardRequest struct {
		Name string `json:"name"`
	}
)

func CreateCard(c echo.Context) error {
	req := new(CreateCardRequest)

	// Bind request body to our request data type
	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusBadRequest, utils.JSON{"error": "Invalid Request"})
	}

	// Check name availability in db
	name := req.Name
	if err := db.CheckNameAvailability(name); err != nil {
		return c.JSON(http.StatusConflict, utils.JSON{"error": "name already taken"})
	}

	// Create new UUID
	UUID := uuid.New()
	// Uuid to binary to fit 16 bytes memory block from mifare
	UUIDBytes := UUID[:]
	// Encode base64 (avoid compatibility issues through the wire)
	// Just for sending, the uuid is decoded back to binary in the front and stored in binary too in the db
	encodedUUID := base64.StdEncoding.EncodeToString(UUIDBytes)

	// Insert record of name and card UUID
	if err := db.CreateCard(UUIDBytes, name); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.JSON{"error": err.Error()})
	}
	response := UUIDResponse{UUID: encodedUUID}

	// Return success code
	return c.JSON(http.StatusOK, response)
}

type ModifyCardRequest struct {
	NewName string `json:"name"`
	UUID    string `json:"uuid"`
}

// Patch, only name is modifiable, uuid is needed as identifier
func ModifyCard(c echo.Context) error {
	req := new(ModifyCardRequest)

	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusBadRequest, utils.JSON{"error": "Invalid Request"})
	}

	UUID := req.UUID
	UUIDBinary, err := base64.StdEncoding.DecodeString(UUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.JSON{"error": "UUID Invalid"})
	}
	err = db.ModifyCard(UUIDBinary, req.NewName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.JSON{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, utils.JSON{"confirmation": "Name updated"})
}

type DeleteCardRequest struct {
	UUID string `json:"uuid"`
}

func DeleteCard(c echo.Context) error {
	req := new(DeleteCardRequest)

	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusBadRequest, utils.JSON{"error": "Invalid Request"})
	}

	UUID := req.UUID
	UUIDBinary, err := base64.StdEncoding.DecodeString(UUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.JSON{"error": "UUID Invalid"})
	}

	err = db.DeleteCard(UUIDBinary)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.JSON{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, utils.JSON{"confirmation": "Card deleted"})
}
