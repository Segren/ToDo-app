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
	"flag"
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
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	r := gin.Default()
	database.Connect()
	routes.SetupRoutes(r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return startServerWithGracefulShutdown(r)
}

func startServerWithGracefulShutdown(handler http.Handler) error {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	addr := flag.String("addr", ":8080", "HTTP network address")
	flag.Parse()

	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  handler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errorLog.Fatalf("Server not started: %v", err)
		}
	}()

	infoLog.Printf("Starting server on %s", *addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	infoLog.Println("Stop signal recieved, server stop...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		errorLog.Printf("Error while stopping server: %v", err)
		return err
	}

	infoLog.Println("Server stopped successfully")
	return nil
}
