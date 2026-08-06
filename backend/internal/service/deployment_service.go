package service

import (
	"context"
	"errors"
	"strings"

	"github.com/RajatGurjar09/project-phoenix/backend/internal/models"
)

var (
	// ErrDeploymentImageRequired is returned when a deployment has no image.
	ErrDeploymentImageRequired = errors.New("deployment image is required")
	// ErrDeploymentStatusRequired is returned when a deployment has no status.
	ErrDeploymentStatusRequired = errors.New("deployment status is required")
)

// DeploymentStore defines the persistence operations needed by DeploymentService.
type DeploymentStore interface {
	CreateDeployment(ctx context.Context, deployment models.Deployment) (models.Deployment, error)
	ListDeploymentsByProject(ctx context.Context, projectID string) ([]models.Deployment, error)
	GetDeploymentByID(ctx context.Context, id string) (models.Deployment, error)
}

// DeploymentService contains deployment domain operations.
type DeploymentService struct {
	deployments DeploymentStore
}

// NewDeploymentService creates a DeploymentService using deployments as its data store.
func NewDeploymentService(deployments DeploymentStore) *DeploymentService {
	return &DeploymentService{deployments: deployments}
}

// CreateDeployment validates and persists deployment.
func (s *DeploymentService) CreateDeployment(ctx context.Context, deployment models.Deployment) (models.Deployment, error) {
	deployment.Image = strings.TrimSpace(deployment.Image)
	if deployment.Image == "" {
		return models.Deployment{}, ErrDeploymentImageRequired
	}

	deployment.Status = strings.TrimSpace(deployment.Status)
	if deployment.Status == "" {
		return models.Deployment{}, ErrDeploymentStatusRequired
	}

	return s.deployments.CreateDeployment(ctx, deployment)
}

// ListDeploymentsByProject returns deployments for projectID.
func (s *DeploymentService) ListDeploymentsByProject(ctx context.Context, projectID string) ([]models.Deployment, error) {
	return s.deployments.ListDeploymentsByProject(ctx, projectID)
}

// GetDeploymentByID returns the deployment identified by id.
func (s *DeploymentService) GetDeploymentByID(ctx context.Context, id string) (models.Deployment, error) {
	return s.deployments.GetDeploymentByID(ctx, id)
}
