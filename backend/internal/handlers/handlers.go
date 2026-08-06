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
	ListProjects(rctx context.Context) ([]models.Project, error)
}

type Handler struct {
	version        string
	startedAt      time.Time
	projectService ProjectService
}

func New(version string, startedAt time.Time, projectService ProjectService) *Handler {
	return &Handler{version: version, startedAt: startedAt, projectService: projectService}
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
