package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"todo-app/internal/middleware"
	"todo-app/internal/service"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to ToDo API",
		})
	})

	// Маршруты для аутентификации
	r.POST("/register", service.Register)
	r.POST("/login", service.Login)

	// Защищенные маршруты для задач
	protected := r.Group("/todos")
	protected.Use(middleware.AuthMiddleware())

	// CRUD операции для задач
	protected.GET("/", service.GetAllToDos)
	protected.POST("/", service.CreateToDo)
	protected.PUT("/:id", service.UpdateToDo)
	protected.DELETE("/:id", service.DeleteToDo)

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
