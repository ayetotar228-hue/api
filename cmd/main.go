package main

import (
	"api/internal/config"
	database "api/internal/db"
	"api/internal/server"
	"fmt"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db, "./migrations"); err != nil {
		log.Printf("Failed to run migrations: %v", err)
	}

	srv := server.NewServer(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), db)
	if err := srv.Run(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
