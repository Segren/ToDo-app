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

	shutdownError := make(chan error)

	//graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit
		app.infoLog.Printf("Caught signal %s, shutting down server...", s)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.infoLog.Printf("completing background tasks")

		app.wg.Wait()
		shutdownError <- nil
	}()

	app.infoLog.Printf("Starting server on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	err := <-shutdownError
	if err != nil {
		return err
	}

	app.infoLog.Printf("Server stoped")

	return nil
}
