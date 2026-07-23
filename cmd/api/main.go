package main

import (
	"Desktop/signoff/internal/api"
	"Desktop/signoff/internal/config"
	"Desktop/signoff/internal/db"
	"Desktop/signoff/internal/jobs"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()
	err := db.Connect(cfg.DataBase_URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	jobs.StartReminderJob()

	app := api.New(&cfg)
	port := cfg.Port

	log.Printf("Server running on http://localhost:%s", port)
	log.Printf("Health check: http://localhost:%s/health", port)
	log.Printf("Webhook: http://localhost:%s/webhook", port)

	if err := http.ListenAndServe(":"+port, app.Router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
