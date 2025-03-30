package utils

import (
	"ToDo/database"
	"errors"
)

func CheckTaskOwnership(userID, taskID int) (bool, error) {
	queryCheck := `SELECT userid FROM tasks WHERE id = $1`
	var checkUserID int
	err := database.DB.QueryRow(queryCheck, taskID).Scan(&checkUserID)
	if err != nil {
		return false, errors.New("task not found")
	}

	return checkUserID == userID, nil
}
