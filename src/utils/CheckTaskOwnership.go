package utils

import (
	"ToDo/database"
	"errors"
	"log"
)

func CheckTaskOwnership(userEmail string, taskID int) (bool, error) {
	queryCheck := `SELECT id FROM tasks WHERE id = $1`
	var checkUserEmail string
	err := database.DB.QueryRow(queryCheck, taskID).Scan(&checkUserEmail)
	log.Printf("email in check: %s", checkUserEmail)
	if err != nil {
		//c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return false, errors.New("task not found")
	}

	return checkUserEmail == userEmail, nil
}
