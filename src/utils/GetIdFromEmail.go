package utils

import (
	"ToDo/database"
	"errors"
)

func GetIdFromEmail(email string) (int, error) {
	id := 0
	query := "SELECT id FROM users WHERE email = $1"
	err := database.DB.QueryRow(query, email).Scan(&id)

	if err != nil {
		return 0, errors.New("couldn't take id from email")
	}

	return id, nil
}
