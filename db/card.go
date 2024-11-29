package db

import "errors"

func CheckNameAvailability(name string) error {
	db := GetDb()
	defer db.Close()

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM user WHERE nombre = ?", name).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("name already taken")
	}
	return nil
}

func CreateCard(UUID []byte, name string) error {
	db := GetDb()
	defer db.Close()

	// Initiate a transaction, guaranteeing that the insertion of both fields is atomic, so every single one is apllied or none
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	result, err := tx.Exec("INSERT INTO user (nombre) VALUES (?)", name)
	if err != nil {
		tx.Rollback()
		return err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO card (UUID, userID) VALUES (?, ?)", UUID, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func ModifyCard(UUID []byte, newName string) error {
	db := GetDb()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Get userID from associated card UUID
	var userID int
	err = tx.QueryRow("SELECT userID FROM card WHERE UUID = ?", UUID).Scan(&userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Replace name with new one
	_, err = tx.Exec("UPDATE user SET nombre = ? WHERE id = ?", newName, userID)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func DeleteCard(UUID []byte) error {
	db := GetDb()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Get userID to delete the card record first
	var userID int
	err = tx.QueryRow("SELECT userID FROM card WHERE UUID = ?", UUID).Scan(&userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM card WHERE UUID = ?", UUID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM user WHERE id = ?", userID)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
