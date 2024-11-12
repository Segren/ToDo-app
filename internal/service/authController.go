package service

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
	"todo-app/internal/database"
	"todo-app/internal/models"
)

var JwtSecret = []byte("секретный_ключ")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with a username and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   user  body     models.User  true  "User data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /register [post]
func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can not hash password"})
	}

	user.Password = string(hashedPassword)

	// Сохранение нового пользователя в БД
	if err := database.DB.Create(&user).Error; err != nil {
		log.Println("Error while creating user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can not create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User successfully created"})
}

// Login godoc
// @Summary Log in a user
// @Description Log in a user and receive a JWT token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   user  body     models.User  true  "User credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /login [post]
func Login(c *gin.Context) {
	var user models.User
	var input models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка существования пользователя
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверное имя пользователя или пароль"})
		return
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверное имя пользователя или пароль"})
	}

	// Создание JWT токена
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: input.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать JWT токен"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
