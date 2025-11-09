package main

import (
	"e-wallet/internal/infrastructure/config"
	"e-wallet/internal/infrastructure/logger"
	"fmt"
)

func main() {
	// Load config
	fmt.Println("[main]: Loading config...")
	cfg := config.MustLoad("configs/config.yaml")
	fmt.Println("[main]: Config loaded")

	// Init logger
	fmt.Println("[main]: Initializing logger...")
	if err := logger.Init(cfg.Log); err != nil {
		fmt.Printf("[main]: Failed to init logger: %v\n", err)
	}
	fmt.Println("[main]: Logger initialized")
}
