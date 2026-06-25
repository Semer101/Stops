package main

import (
	"context"
	"log"

	"stops/backend/internal/config"
	"stops/backend/internal/migrations"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	db, err := config.NewDatabasePool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database setup failed: %v", err)
	}
	defer db.Close()

	if err := migrations.NewRunner(db).Up(ctx); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("migrations applied successfully")
}
