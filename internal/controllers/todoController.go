package controllers

import (
	"net/http"
	"todo-app/internal/database"
	"todo-app/internal/models"

	"github.com/gin-gonic/gin"
)

func GetAllToDos(c *gin.Context) {
	var todos []models.ToDo
	database.DB.Find(&todos)
	c.JSON(http.StatusOK, todos)
}

// CreateToDo godoc
// @Summary Create a new ToDo
// @Description Create a new ToDo item
// @Tags todos
// @Accept  json
// @Produce  json
// @Param   todo  body     models.ToDo  true  "ToDo"
// @Success 200 {object} models.ToDo
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /todos [post]
func CreateToDo(c *gin.Context) {
	var todo models.ToDo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Create(&todo)
	c.JSON(http.StatusOK, todo)
}

// UpdateToDo godoc
// @Summary Update a new ToDo
// @Description Update a new ToDo item
// @Tags todos
// @Accept  json
// @Produce  json
// @Param   todo  body     models.ToDo  true  "ToDo"
// @Success 200 {object} models.ToDo
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /todos [put]
func UpdateToDo(c *gin.Context) {
	var todo models.ToDo
	id := c.Param("id")
	if err := database.DB.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
		return
	}
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Save(&todo)
	c.JSON(http.StatusOK, todo)
}

// DeleteToDo godoc
// @Summary Delete a new ToDo
// @Description Delete a new ToDo item
// @Tags todos
// @Accept  json
// @Produce  json
// @Param   todo  body     models.ToDo  true  "ToDo"
// @Success 200 {object} models.ToDo
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /todos [delete]
func DeleteToDo(c *gin.Context) {
	var todo models.ToDo
	id := c.Param("id")
	if err := database.DB.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
		return
	}
	database.DB.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{"message": "Задача удалена"})
}
