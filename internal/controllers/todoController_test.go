package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo-app/internal/database"
	"todo-app/internal/models"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	// Создаём in-memory базу данных для тестов
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.ToDo{})
	database.DB = db
}

func TestGetAllToDos(t *testing.T) {
	setupTestDB()

	// Инициализация Gin и маршрута
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/todos", GetAllToDos)

	// Создадим несколько тестовых задач
	database.DB.Create(&models.ToDo{Title: "Test ToDo 1", Description: "Description 1", Completed: false})
	database.DB.Create(&models.ToDo{Title: "Test ToDo 2", Description: "Description 2", Completed: true})

	// Отправляем GET запрос на /todos
	req, _ := http.NewRequest("GET", "/todos", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем содержимое ответа
	var todos []models.ToDo
	err := json.Unmarshal(w.Body.Bytes(), &todos)
	assert.NoError(t, err)
	assert.Len(t, todos, 2) // Ожидаем, что в ответе будет 2 задачи
}

func TestCreateToDo(t *testing.T) {
	setupTestDB()

	// Инициализация Gin и маршрута
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/todos", CreateToDo)

	// Создадим новую задачу
	todo := models.ToDo{
		Title:       "New ToDo",
		Description: "New Description",
		Completed:   false,
	}
	jsonValue, _ := json.Marshal(todo)

	// Отправляем POST запрос на /todos
	req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем содержимое ответа
	var createdTodo models.ToDo
	err := json.Unmarshal(w.Body.Bytes(), &createdTodo)
	assert.NoError(t, err)
	assert.Equal(t, todo.Title, createdTodo.Title)
	assert.Equal(t, todo.Description, createdTodo.Description)
}

func TestUpdateToDo(t *testing.T) {
	setupTestDB()

	// Инициализация Gin и маршрута
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/todos/:id", UpdateToDo)

	// Создадим тестовую задачу
	todo := models.ToDo{Title: "Old Title", Description: "Old Description", Completed: false}
	database.DB.Create(&todo)

	// Обновим задачу
	updatedTodo := models.ToDo{
		Title:       "Updated Title",
		Description: "Updated Description",
		Completed:   true,
	}
	jsonValue, _ := json.Marshal(updatedTodo)

	// Отправляем PUT запрос на /todos/:id
	req, _ := http.NewRequest("PUT", "/todos/"+strconv.FormatUint(uint64(todo.ID), 10), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем содержимое ответа
	var updatedToDo models.ToDo
	err := json.Unmarshal(w.Body.Bytes(), &updatedToDo)
	assert.NoError(t, err)
	assert.Equal(t, updatedTodo.Title, updatedToDo.Title)
}

func TestDeleteToDo(t *testing.T) {
	setupTestDB()

	// Инициализация Gin и маршрута
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.DELETE("/todos/:id", DeleteToDo)

	// Создадим тестовую задачу
	todo := models.ToDo{Title: "Test ToDo", Description: "Test Description", Completed: false}
	database.DB.Create(&todo)

	// Отправляем DELETE запрос на /todos/:id
	req, _ := http.NewRequest("DELETE", "/todos/"+strconv.FormatUint(uint64(todo.ID), 10), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем, что задача была удалена
	var deletedToDo models.ToDo
	err := database.DB.First(&deletedToDo, todo.ID).Error
	assert.Error(t, err) // Ожидаем ошибку, т.к. задача должна быть удалена
}
