package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/RajatGurjar09/project-phoenix/backend/internal/models"
	"github.com/RajatGurjar09/project-phoenix/backend/internal/repository"
	"github.com/RajatGurjar09/project-phoenix/backend/internal/service"
)

type createProjectRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// CreateProject creates a project from the request body.
func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var request createProjectRequest
	if err := decodeJSON(r, &request); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	project, err := h.projectService.CreateProject(r.Context(), models.Project{
		Name:        request.Name,
		Description: request.Description,
	})
	if err != nil {
		handleProjectError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, project)
}

// ListProjects returns every project.
func (h *Handler) ListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.projectService.ListProjects(r.Context())
	if err != nil {
		handleProjectError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, projects)
}

// GetProject returns the project identified by the id path parameter.
func (h *Handler) GetProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "project id is required"})
		return
	}

	project, err := h.projectService.GetProjectByID(r.Context(), id)
	if err != nil {
		handleProjectError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, project)
}

// DeleteProject removes the project identified by the id path parameter.
func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	if err := h.projectService.DeleteProject(r.Context(), r.PathValue("id")); err != nil {
		handleProjectError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func decodeJSON(r *http.Request, destination any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return fmt.Errorf("decode request body: %w", err)
	}

	return nil
}

func handleProjectError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrProjectNameRequired):
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
	case errors.Is(err, repository.ErrProjectNotFound):
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "project not found"})
	default:
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
	}
}
