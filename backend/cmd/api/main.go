package main

import (
	"context"
	"log"
	"time"

	"github.com/RajatGurjar09/project-phoenix/backend/internal/config"
	"github.com/RajatGurjar09/project-phoenix/backend/internal/database"
	"github.com/RajatGurjar09/project-phoenix/backend/internal/docker"
	"github.com/RajatGurjar09/project-phoenix/backend/internal/repository"
	"github.com/RajatGurjar09/project-phoenix/backend/internal/server"
	"github.com/RajatGurjar09/project-phoenix/backend/internal/service"
	"github.com/RajatGurjar09/project-phoenix/backend/internal/version"
)

func main() {
	cfg := config.Load()
	databaseConfig := database.LoadConfig()
	databaseContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	db, err := database.Connect(databaseContext, databaseConfig)
	cancel()
	if err != nil {
		log.Fatalf("database startup failed: %v", err)
	}
	defer db.Close()

	projectRepository := repository.NewProjectRepository(db)
	projectService := service.NewProjectService(projectRepository)
	deploymentRepository := repository.NewDeploymentRepository(db)
	dockerRuntime := docker.NewRuntime()
	deploymentService := service.NewDeploymentService(deploymentRepository, dockerRuntime)
	api := server.New(cfg.Address, version.Version, time.Now(), projectService, deploymentService)

	log.Printf("API listening on %s", cfg.Address)
	if err := api.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
