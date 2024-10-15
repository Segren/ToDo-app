package routes

import (
	"github.com/gin-gonic/gin"
	"todo-app/internal/controllers"
	"todo-app/internal/middleware"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to ToDo API",
		})
	})

	// Маршруты для аутентификации
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// Защищенные маршруты для задач
	protected := r.Group("/todos")
	protected.Use(middleware.AuthMiddleware())

	// CRUD операции для задач
	protected.GET("/", controllers.GetAllToDos)
	protected.POST("/", controllers.CreateToDo)
	protected.PUT("/:id", controllers.UpdateToDo)
	protected.DELETE("/:id", controllers.DeleteToDo)
}
