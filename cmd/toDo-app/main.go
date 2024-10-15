// @title ToDo API
// @version 1.0
// @description This is a simple ToDo API application.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "todo-app/docs"
	"todo-app/internal/database"
	"todo-app/internal/routes"
)

func main() {
	r := gin.Default()
	database.Connect()
	routes.SetupRoutes(r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	startServerWithGracefulShutdown(r)
}

func startServerWithGracefulShutdown(handler http.Handler) {
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Не удалось запустить сервер: %v", err)
		}
	}()

	log.Println("Сервер запущен на http://localhost:8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Получен сигнал завершения работы, завершение сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при завершении работы сервера: %v", err)
	}
	log.Println("Сервер успешно завершил работу")
}
