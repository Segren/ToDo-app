package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo-app/internal/database"
	"todo-app/internal/routes"
)

func (app *application) startServerWithGracefulShutdown() error {
	r := gin.Default()
	database.Connect()
	routes.SetupRoutes(r)

	addr := fmt.Sprintf(":%d", app.config.port)

	server := &http.Server{
		Addr:     addr,
		ErrorLog: app.errorLog,
		Handler:  r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.errorLog.Fatalf("Server not started: %v", err)
		}
	}()

	app.infoLog.Printf("Starting server on %s", addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	app.infoLog.Println("Stop signal recieved, server stop...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		app.errorLog.Printf("Error while stopping server: %v", err)
		return err
	}

	app.infoLog.Println("Server stopped successfully")
	return nil
}
