package main

import (
	"checkin/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Name
	e.POST("/card/create", handlers.CreateCard)
	// NewName
	e.PATCH("/card/edit", handlers.ModifyCard)
	// UUID
	e.DELETE("/card/delete", handlers.DeleteCard)
	// UUID
	e.POST("/checkin", handlers.Checkin)

	// Nothing
	e.GET("/report", handlers.GetReport)

	e.Logger.Fatal(e.Start("127.0.0.1:6969"))
}
