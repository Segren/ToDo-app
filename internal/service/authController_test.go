package service

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"todo-app/internal/database"
	"todo-app/internal/models"
)

func TestMain(m *testing.M) {
	os.Setenv("GO_ENV", "test")
	database.Connect()

	code := m.Run()
	os.Exit(code)
}

func setupTestDatabase() {
	if database.DB == nil {
		log.Fatal("Database is not initialized")
	}
	// Очищаем базу данных перед каждым тестом
	database.DB.Exec("DELETE FROM users")
}

func TestRegister(t *testing.T) {
	setupTestDatabase()
	// Инициализация Gin
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Регистрируем маршрут
	r.POST("/register", Register)

	// Пример данных для запроса
	user := models.User{
		Username: "testuser",
		Password: "password123",
	}
	jsonValue, _ := json.Marshal(user)

	// Создаём новый HTTP запрос
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	// Используем httptest для имитации ответа
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем, что код ответа 200 OK
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Проверяем, что пользователь добавлен в базу данных
	var createdUser models.User
	if err := database.DB.Where("username = ?", user.Username).First(&createdUser).Error; err == gorm.ErrRecordNotFound {
		t.Fatalf("User is not created")
	}
}

func TestLogin(t *testing.T) {
	setupTestDatabase()
	// Создаём пользователя напрямую в базе данных для теста логина
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{
		Username: "testuser",
		Password: string(hashedPassword),
	}
	database.DB.Create(&user)

	// Инициализация Gin
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Регистрируем маршрут
	r.POST("/login", Login)

	// Пример данных для запроса
	loginPayload := models.User{
		Username: "testuser",
		Password: "password123",
	}
	jsonValue, _ := json.Marshal(loginPayload)

	// Создаём новый HTTP запрос
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	// Используем httptest для имитации ответа
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем, что код ответа 200 OK
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Проверяем, что в ответе есть JWT токен
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil || response["token"] == "" {
		t.Fatalf("Token wos not created or invalid")
	}
}
