package main

import (
	"ToDo/database"
	routes "ToDo/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Ошибка при загрузке файла .env: %v", err)
	}

	database.Init()
	port := "8080"
	router := gin.New()
	routes.ToDoRoutes(router)

	router.Run(":" + port)

}
