package routes

import (
	controllers "ToDo/controller"
	"ToDo/midlleware"

	"github.com/gin-gonic/gin"
)

func ToDoRoutes(r *gin.Engine) {
	r.POST("/login", controllers.Login)
	r.POST("/signup", controllers.Signup)
	r.GET("/home", controllers.Home)
	r.GET("/logout", controllers.Logout)
	userGroup := r.Group("/todo").Use(midlleware.IsAuthorized())
	{
		userGroup.POST("/", controllers.CreateTask)
		userGroup.GET("/", controllers.GetAllTasks)
		userGroup.GET("/:id", controllers.GetTaskById)
		userGroup.GET("/user/", controllers.GetTasksByUserId)
		userGroup.PUT("/:id", controllers.UpdateTask)
		userGroup.DELETE("/:id", controllers.DeleteTask)
	}
}
