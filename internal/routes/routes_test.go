package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo-app/internal/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
	"todo-app/internal/database"
	"todo-app/internal/service"
)

func setupRoutesTestDB() {
	// Создаём in-memory базу данных для тестов
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{}, &models.ToDo{})
	database.DB = db
}

func setupRouter() *gin.Engine {
	// Инициализация Gin в тестовом режиме
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	SetupRoutes(r)
	return r
}

// Вспомогательная функция для хеширования пароля
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Вспомогательная функция для генерации JWT токена
func generateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &service.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(service.JwtSecret)
	return tokenString, err
}

func TestHomeRoute(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем код ответа и содержимое
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Welcome to ToDo API")
}

func TestRegisterRoute(t *testing.T) {
	setupRoutesTestDB()
	router := setupRouter()

	// Создадим запрос на регистрацию
	user := models.User{
		Username: "testuser",
		Password: "password123",
	}
	jsonValue, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User successfully created")
}

func TestLoginRoute(t *testing.T) {
	setupRoutesTestDB()
	router := setupRouter()

	// Создадим тестового пользователя
	user := models.User{
		Username: "testuser",
		Password: "password123",
	}
	hashedPassword, _ := hashPassword(user.Password)
	user.Password = hashedPassword
	database.DB.Create(&user)

	// Создадим запрос на логин
	loginData := models.User{
		Username: "testuser",
		Password: "password123",
	}
	jsonValue, _ := json.Marshal(loginData)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем наличие JWT токена в ответе
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotEmpty(t, response["token"])
}

func TestProtectedRouteWithoutToken(t *testing.T) {
	router := setupRouter()

	// Попытаемся получить задачи без токена
	req, _ := http.NewRequest("GET", "/todos/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Ожидаем 401 статус, так как запрос без токена
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Missing token")
}

func TestProtectedRouteWithToken(t *testing.T) {
	setupRoutesTestDB()
	router := setupRouter()

	// Создадим тестового пользователя
	user := models.User{
		Username: "testuser",
		Password: "password123",
	}
	hashedPassword, _ := hashPassword(user.Password)
	user.Password = hashedPassword
	database.DB.Create(&user)

	// Генерируем JWT токен для пользователя
	tokenString, _ := generateJWT(user.Username)

	// Создадим запрос для получения задач с токеном
	req, _ := http.NewRequest("GET", "/todos/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Ожидаем 200 статус
	assert.Equal(t, http.StatusOK, w.Code)
}
