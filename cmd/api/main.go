package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

	"lead-automation/internal/ai"
	"lead-automation/internal/automation"
	"lead-automation/internal/lead"
)

func main() {
	ctx := context.Background()
	if err := godotenv.Load(); err != nil {
    log.Println("Warning: .env file not found")
}
		databaseURL := os.Getenv("DATABASE_URL")
		if databaseURL == "" {
			log.Fatal("DATABASE_URL is required")
		}

	db, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)

	leadRepo := lead.NewRepository(db)

	qualifier := ai.NewQualifier(
		os.Getenv("AI_URL"),
		"",
	)
	
	n8n := automation.NewN8N(
		os.Getenv("N8N_WEBHOOK_URL"),
	)

	leadService := lead.NewService(
		leadRepo,
		qualifier,
		n8n,
	)

	leadHandler := lead.NewHandler(leadService)

	router := chi.NewRouter()

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	router.Post("/api/v1/leads", leadHandler.Create)

	server := &http.Server{
		Addr:              ":8081",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("API listening on :8081")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}