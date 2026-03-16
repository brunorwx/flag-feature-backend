package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/brunorwx/flagAPI/internal/application"
	"github.com/brunorwx/flagAPI/internal/infrastructure"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	repo := infrastructure.NewInMemoryFeatureFlagRepository()
	service := application.NewFeatureFlagService(repo)
	handler := application.NewHandler(service)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"healthy"}`)
	})

	r.Post("/flags", handler.CreateFlag)
	r.Get("/flags", handler.ListFlags)
	r.Get("/flags/{key}", handler.GetFlag)
	r.Put("/flags/{key}/global", handler.SetGlobalState)
	r.Put("/flags/{key}/users/{userId}", handler.SetUserOverride)
	r.Get("/evaluate/{key}", handler.EvaluateFlag)

	port := ":8080"
	log.Printf("Starting Feature Flag API on %s\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
