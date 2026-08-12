package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/RajatGurjar09/project-phoenix/backend/internal/models"
	"github.com/RajatGurjar09/project-phoenix/backend/internal/repository"
	"github.com/RajatGurjar09/project-phoenix/backend/internal/service"
)

type createDeploymentRequest struct {
	Image  string `json:"image"`
	Status string `json:"status"`
}

// CreateDeployment creates a deployment for the project identified by the id path parameter.
func (h *Handler) CreateDeployment(w http.ResponseWriter, r *http.Request) {
	var request createDeploymentRequest
	if err := decodeJSON(r, &request); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	deployment, err := h.deploymentService.CreateDeployment(r.Context(), models.Deployment{
		ProjectID: r.PathValue("id"),
		Image:     request.Image,
		Status:    request.Status,
	})
	if err != nil {
		handleDeploymentError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, deployment)
}

// ListDeploymentsByProject returns deployments for the project identified by the id path parameter.
func (h *Handler) ListDeploymentsByProject(w http.ResponseWriter, r *http.Request) {
	deployments, err := h.deploymentService.ListDeploymentsByProject(r.Context(), r.PathValue("id"))
	if err != nil {
		handleDeploymentError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, deployments)
}

// GetDeployment returns the deployment identified by the id path parameter.
func (h *Handler) GetDeployment(w http.ResponseWriter, r *http.Request) {
	deployment, err := h.deploymentService.GetDeploymentByID(r.Context(), r.PathValue("id"))
	if err != nil {
		handleDeploymentError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, deployment)
}

// StopDeployment stops the container for the deployment identified by the id path parameter.
func (h *Handler) StopDeployment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "deployment id is required"})
		return
	}

	deployment, err := h.deploymentService.StopDeployment(r.Context(), id)
	if err != nil {
		handleDeploymentError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, deployment)
}

// RestartDeployment restarts the container for the deployment identified by the id path parameter.
func (h *Handler) RestartDeployment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "deployment id is required"})
		return
	}

	deployment, err := h.deploymentService.RestartDeployment(r.Context(), id)
	if err != nil {
		handleDeploymentError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, deployment)
}

// RemoveDeployment removes the container and marks the deployment as removed.
func (h *Handler) RemoveDeployment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "deployment id is required"})
		return
	}

	deployment, err := h.deploymentService.RemoveDeployment(r.Context(), id)
	if err != nil {
		handleDeploymentError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, deployment)
}

func handleDeploymentError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrDeploymentImageRequired),
		errors.Is(err, service.ErrDeploymentStatusRequired),
		errors.Is(err, service.ErrDeploymentContainerNotFound),
		errors.Is(err, service.ErrInvalidDeploymentState):
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
	case errors.Is(err, repository.ErrDeploymentNotFound):
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "deployment not found"})
	default:
	        log.Printf("deployment error: %v", err)
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
	}
}
