package main

import (
	"api/internal/server"
	"api/pkg/config"
	database "api/pkg/db"
	"api/pkg/storage"
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

	store := storage.NewSQLxStorage(db)
	defer store.Close()

	if err := database.RunMigrations(db, "./migrations"); err != nil {
		log.Printf("Failed to run migrations: %v", err)
	}

	server := server.NewServer(cfg.GetAddr(), store)
	if err := server.Run(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
