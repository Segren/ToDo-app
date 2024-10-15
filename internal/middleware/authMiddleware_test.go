package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"todo-app/internal/controllers"
)

// Генерация валидного JWT токена для тестов
func generateValidJWT() string {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &controllers.Claims{
		Username: "testuser",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(controllers.JwtSecret)
	return tokenString
}

// тестирование middleware без токена
func TestAuthMiddleware_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	expectedBody := `{"error":"Missing token"}`
	if strings.TrimSpace(w.Body.String()) != expectedBody {
		t.Errorf("Expected response %s, got %s", expectedBody, w.Body.String())
	}
}

// Тестирование middleware с невалидным токеном
func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusUnauthorized, w.Code)
	}

	expectedBody := `{"error":"Imvalid token"}`
	if strings.TrimSpace(w.Body.String()) != expectedBody {
		t.Errorf("Ожидался ответ %s, получен %s", expectedBody, w.Body.String())
	}
}

// Тестирование middleware с валидным токеном
func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	})

	token := generateValidJWT()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, recieved %d", http.StatusOK, w.Code)
	}

	expectedBody := `{"message":"Success"}`
	if strings.TrimSpace(w.Body.String()) != expectedBody {
		t.Errorf("Expected response %s, recieved %s", expectedBody, w.Body.String())
	}
}
