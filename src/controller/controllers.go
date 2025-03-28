package controllers

import (
	"ToDo/database"
	"ToDo/models"
	"ToDo/utils"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// type Task struct {
// 	Id          uint
// 	Title       string
// 	Description string
// 	Completed   bool
// 	UserId      uint
// 	CreatedAt   time.Time
// 	UpdatedAt   time.Time
// }

// Creating tasks fields in database
func CreateTask(c *gin.Context) {
	var task models.Task

	if err := c.ShouldBindBodyWithJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "c.ShouldBindBodyWithJSON(&task) didn't work"})
		return
	}

	userEmail, err := utils.ExtractUserEmailFromToken(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	query := `SELECT id FROM users WHERE email = $1`
	if err := database.DB.QueryRow(query, userEmail).Scan(&task.UserId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Print("UserId:", task.UserId)
	query = `INSERT INTO tasks (title, description, completed, userid, created_at, updated_at) VALUES ($1,$2,$3,$4,NOW(),NOW())`
	if _, err := database.DB.Exec(query, task.Title, task.Description, task.Completed, task.UserId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// Let's say we got an admin func without admin role just for tests
func GetAllTasks(c *gin.Context) {
	//Parameters for pagination
	//If Parameters for pagination wasn't declared, we set default values
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if page < 1 || err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)

	if limit > 10 || err != nil {
		limit = 10
	}

	offset := (page - 1) * limit

	var Tasks []models.Task

	rows, err := database.DB.Query(`SELECT * FROM tasks LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var Task models.Task
		if err := rows.Scan(&Task.Id, &Task.Title, &Task.Description, &Task.Completed, &Task.UserId, &Task.CreatedAt, &Task.UpdatedAt); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		Tasks = append(Tasks, Task)

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	var total int

	err = database.DB.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&total)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": Tasks,
		"meta": gin.H{
			"total":      total,
			"page":       page,
			"limit":      limit,
			"totalPages": int(math.Ceil(float64(total) / float64(limit))),
		},
	})
}

func GetTaskById(c *gin.Context) {
	var Task models.Task
	id := c.Param("id")

	userEmail, err := utils.ExtractUserEmailFromToken(c)
	fmt.Println(userEmail)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, couldn't extract user id"})
		return
	}

	taskID, err := strconv.Atoi(id)
	log.Printf("%d", taskID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//	isOwner, err := utils.CheckTaskOwnership(userId, taskID)
	isOwner, err := utils.CheckTaskOwnership(userEmail, taskID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only get your own tasks"})
		return
	}

	query := "SELECT * FROM tasks WHERE id = $1"

	err = database.DB.QueryRow(query, id).Scan(&Task.Id, &Task.Title, &Task.Description, &Task.Completed, &Task.UserId, &Task.CreatedAt, &Task.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task didn't found"})
		return
	}

	c.JSON(http.StatusOK, Task)
}

func GetTasksByUserId(c *gin.Context) {
	var Tasks []models.Task

	userEmail, err := utils.ExtractUserEmailFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	query := "SELECT * FROM tasks WHERE userid = $1"

	rows, err := database.DB.Query(query, userEmail)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No tasks by this user"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var Task models.Task
		if err := rows.Scan(&Task.Id, &Task.Title, &Task.Description, &Task.Completed, &Task.UserId, &Task.CreatedAt, &Task.UpdatedAt); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		Tasks = append(Tasks, Task)

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, Tasks)
}

func UpdateTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	userEmail, err := utils.ExtractUserEmailFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	taskID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	isOwner, err := utils.CheckTaskOwnership(userEmail, taskID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own tasks"})
		return
	}

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "UPDATE tasks SET title = $1, description = $2, updated_at = NOW() WHERE id = $3 AND email = $4"

	_, err = database.DB.Exec(query, task.Title, task.Description, id, userEmail)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	userId, err := utils.ExtractUserEmailFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	taskID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	isOwner, err := utils.CheckTaskOwnership(userId, taskID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own tasks"})
		return
	}

	query := "DELETE FROM tasks WHERE id=$1"
	if _, err := database.DB.Exec(query, id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task didn't found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
