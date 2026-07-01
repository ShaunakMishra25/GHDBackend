package main

import (
	"context"
	"encoding/base64"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gumla-hds/gumla-backend/config"
	"github.com/gumla-hds/gumla-backend/internal/server"
	"github.com/gumla-hds/gumla-backend/pkg/database"
	"github.com/gumla-hds/gumla-backend/pkg/firebase"
	"github.com/gumla-hds/gumla-backend/pkg/razorpay"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Database connected")

	if err := runMigrations(ctx, db, "migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed")

	var fb *firebase.Client
	if cfg.Firebase.CredentialsJSON != "" {
		creds := cfg.Firebase.CredentialsJSON
		if decoded, err := base64.StdEncoding.DecodeString(creds); err == nil {
			creds = string(decoded)
		}
		fb, err = firebase.NewClient(ctx, creds, cfg.Firebase.StorageBucket)
		if err != nil {
			log.Fatalf("Failed to init firebase: %v", err)
		}
		log.Println("Firebase initialized")
	}

	rp := razorpay.NewClient(cfg.Razorpay.KeyID, cfg.Razorpay.KeySecret, cfg.Razorpay.WebhookSecret)
	log.Println("Razorpay client initialized")

	srv := server.New(cfg, db, fb, rp)

	if err := srv.Start(cfg.Port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func runMigrations(ctx context.Context, pool *pgxpool.Pool, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		data, err := os.ReadFile(dir + "/" + entry.Name())
		if err != nil {
			return err
		}

		if _, err := pool.Exec(ctx, string(data)); err != nil {
			return err
		}

		log.Printf("Migration applied: %s", entry.Name())
	}

	return nil
}
