package server

import (
	"net/http"
	"time"

	"github.com/RajatGurjar09/project-phoenix/backend/internal/handlers"
	"github.com/RajatGurjar09/project-phoenix/backend/internal/middleware"
	"github.com/RajatGurjar09/project-phoenix/backend/internal/service"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 10 * time.Second
	idleTimeout       = 60 * time.Second
)

func New(
	address string,
	version string,
	startedAt time.Time,
	projectService *service.ProjectService,
	deploymentService *service.DeploymentService,
) *http.Server {
	h := handlers.New(version, startedAt, projectService, deploymentService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /version", h.Version)
	mux.HandleFunc("POST /projects", h.CreateProject)
	mux.HandleFunc("GET /projects", h.ListProjects)
	mux.HandleFunc("GET /projects/{id}", h.GetProject)
	mux.HandleFunc("PATCH /projects/{id}", h.UpdateProject)
	mux.HandleFunc("DELETE /projects/{id}", h.DeleteProject)
	mux.HandleFunc("POST /projects/{id}/deployments", h.CreateDeployment)
	mux.HandleFunc("GET /projects/{id}/deployments", h.ListDeploymentsByProject)
	mux.HandleFunc("GET /deployments/{id}", h.GetDeployment)

	return &http.Server{
		Addr:              address,
		Handler:           middleware.RequestLogger(mux),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}
}
