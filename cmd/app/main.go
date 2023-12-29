package main

import (
	"log"

	"lemon_be/config"
	"lemon_be/internal/app"
)

func main() {
	// Configuration

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
