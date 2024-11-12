package database

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"todo-app/internal/models"
)

var DB *gorm.DB

func Connect() {
	if os.Getenv("GO_ENV") == "test" {
		// Используем SQLite в памяти для тестов
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			log.Fatal("Не удалось подключиться к in-memory базе данных: ", err)
		}
		db.AutoMigrate(&models.User{}, &models.ToDo{})
		DB = db
		return
	}
	// Загружаем переменные из.env файла в окружение
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error while loading .env file: $v", err)
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных: ", err)
	}

	db.AutoMigrate(&models.User{}, &models.ToDo{}) // Автоматическая миграция таблицы для модели ToDo
	DB = db
}
