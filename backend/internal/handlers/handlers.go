package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/RajatGurjar09/project-phoenix/backend/internal/models"
)

// ProjectService defines the project operations used by HTTP handlers.
type ProjectService interface {
	CreateProject(rctx context.Context, project models.Project) (models.Project, error)
	GetProjectByID(rctx context.Context, id string) (models.Project, error)
	UpdateProject(rctx context.Context, id string, project models.Project) (models.Project, error)
	ListProjects(rctx context.Context) ([]models.Project, error)
	DeleteProject(rctx context.Context, id string) error
}

// DeploymentService defines the deployment operations used by HTTP handlers.
type DeploymentService interface {
	CreateDeployment(rctx context.Context, deployment models.Deployment) (models.Deployment, error)
	ListDeploymentsByProject(rctx context.Context, projectID string) ([]models.Deployment, error)
	GetDeploymentByID(rctx context.Context, id string) (models.Deployment, error)
	StopDeployment(rctx context.Context, id string) (models.Deployment, error)
	RestartDeployment(rctx context.Context, id string) (models.Deployment, error)
	RemoveDeployment(rctx context.Context, id string) (models.Deployment, error)
}

type Handler struct {
	version           string
	startedAt         time.Time
	projectService    ProjectService
	deploymentService DeploymentService
}

func New(
	version string,
	startedAt time.Time,
	projectService ProjectService,
	deploymentService DeploymentService,
) *Handler {
	return &Handler{
		version:           version,
		startedAt:         startedAt,
		projectService:    projectService,
		deploymentService: deploymentService,
	}
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status":    "ok",
		"version":   h.version,
		"uptime":    time.Since(h.startedAt).String(),
		"hostname":  hostname,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *Handler) Version(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"version": h.version})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
