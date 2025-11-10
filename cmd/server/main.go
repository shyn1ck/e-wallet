package main

import (
	"context"
	"e-wallet/internal/infrastructure/config"
	"e-wallet/internal/infrastructure/container"
	"e-wallet/internal/infrastructure/logger"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load config
	fmt.Println("Loading config...")
	cfg := config.MustLoad("configs/config.yaml")
	fmt.Println("Config loaded")

	// Init logger
	fmt.Println("Initializing logger...")
	if err := logger.Init(cfg.Log); err != nil {
		fmt.Printf("Failed to init logger: %v\n", err)
		os.Exit(1)
	}
	logger.Info.Println("Logger initialized successfully")

	// Initialize DI container
	fmt.Println("Initializing application container...")
	app, err := container.NewContainer(cfg)
	if err != nil {
		logger.Error.Printf("Failed to initialize container: %v", err)
		os.Exit(1)
	}
	defer func() {
		if err := app.Close(); err != nil {
			logger.Error.Printf("Failed to close container: %v", err)
		}
	}()
	fmt.Println("Application container initialized successfully")

	fmt.Println("Initializing server...")
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      app.Router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in goroutine
	go func() {
		fmt.Println("Starting server...")
		fmt.Printf("Starting %s v%s in %s mode\n", cfg.App.Name, cfg.App.Version, cfg.App.Environment)
		fmt.Printf("Server listening on %s\n", server.Addr)
		logger.Info.Printf("Server listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error.Printf("Server failed to start: %v", err)
			fmt.Printf("Server failed to start: %v\n", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info.Println("Shutting down server...")

	// Graceful shutdown with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error.Printf("Server forced to shutdown: %v", err)
		os.Exit(1)
	}

	logger.Info.Println("Server exited gracefully")
}
