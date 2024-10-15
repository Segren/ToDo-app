package database

import (
	"os"
	"testing"
	"todo-app/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestConnectToSQLite(t *testing.T) {
	// Устанавливаем переменную окружения для тестового режима
	os.Setenv("GO_ENV", "test")

	// Вызываем функцию подключения
	Connect()

	// Проверяем, что база данных не равна nil
	assert.NotNil(t, DB)

	// Проверяем, что таблицы для моделей User и ToDo были созданы
	err := DB.AutoMigrate(&models.User{}, &models.ToDo{})
	assert.NoError(t, err)
}
