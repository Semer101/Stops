package main

import (
	"context"
	"log"

	"stops/backend/internal/config"
	"stops/backend/internal/routes"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	db, err := config.NewDatabasePool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database setup failed: %v", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}

	router := routes.NewRouter(cfg, db)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
