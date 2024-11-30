package db

import (
	"database/sql"
	"errors"
	"time"
)

func Checkin(UUID []byte) error {
	date := time.Now()

	db := GetDb()
	defer db.Close()

	var cardID int
	err := db.QueryRow("SELECT id FROM card WHERE UUID = ?", UUID).Scan(&cardID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("card not found in db")
		}
		return err
	}

	_, err = db.Exec("INSERT INTO checkinLog (idTarjeta, date) VALUES (?, ?)", cardID, date)
	if err != nil {
		return err
	}

	return nil
}
