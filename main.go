package main

import (
	"checkin/handlers"
	"fmt"
	"os"

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%s", port)))
}
