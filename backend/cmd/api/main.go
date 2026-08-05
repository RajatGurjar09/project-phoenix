package main

import (
	"log"
	"time"

	"github.com/RajatGurjar09/project-phoenix/backend/internal/config"
	"github.com/RajatGurjar09/project-phoenix/backend/internal/server"
	"github.com/RajatGurjar09/project-phoenix/backend/internal/version"
)

func main() {
	cfg := config.Load()
	api := server.New(cfg.Address, version.Version, time.Now())

	log.Printf("API listening on %s", cfg.Address)
	if err := api.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
