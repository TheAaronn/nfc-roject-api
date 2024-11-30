package handlers

import (
	"checkin/db"
	"checkin/utils"
	"encoding/csv"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

func GetReport(c echo.Context) error {
	// Create temporary file (will be deleted once closed)
	report, err := os.CreateTemp("", "report-*.csv")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.JSON{"error": "Failed to create temp report file"})
	}
	defer os.Remove(report.Name())

	db := db.GetDb()
	defer db.Close()

	// Query all logs (only name and lastCheckin)
	rows, err := db.Query("SELECT user.nombre, checkinLog.date FROM checkinLog JOIN card ON checkinLog.idTarjeta = card.id JOIN user ON card.userID = user.id")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.JSON{"error": err.Error()})
	}
	defer rows.Close()

	logs := make([][]string, 0)
	var log utils.Log

	// Read each row returned from the query
	for rows.Next() {
		if err := rows.Scan(&log.Name, &log.Date); err != nil {
			return c.JSON(http.StatusInternalServerError, utils.JSON{"error": "Failed to read report row from db"})
		}
		logs = append(logs, []string{log.Name, log.Date.Format(time.RFC3339)})
	}

	writer := csv.NewWriter(report)

	// Instead of writing one log by one log, write them all at once
	writer.WriteAll(logs)
	if err := writer.Error(); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.JSON{"error": "Error writing to csv report"})
	}
	// Close file to be able to save the changes before serving it to client
	if err := report.Close(); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.JSON{"error": err.Error()})
	}

	if err := c.Attachment(report.Name(), "report.csv"); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to serve CSV file"})
	}

	return nil
}
