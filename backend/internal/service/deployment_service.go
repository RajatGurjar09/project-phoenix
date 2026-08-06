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

// ContainerRuntime defines container lifecycle operations needed by DeploymentService.
type ContainerRuntime interface {
	CreateContainer(ctx context.Context, imageName string) (string, error)
	StartContainer(ctx context.Context, containerID string) error
}

// DeploymentStore defines the persistence operations needed by DeploymentService.
type DeploymentStore interface {
	CreateDeployment(ctx context.Context, deployment models.Deployment) (models.Deployment, error)
	UpdateDeploymentStatus(ctx context.Context, id string, status string) (models.Deployment, error)
	ListDeploymentsByProject(ctx context.Context, projectID string) ([]models.Deployment, error)
	GetDeploymentByID(ctx context.Context, id string) (models.Deployment, error)
}

// DeploymentService contains deployment domain operations.
type DeploymentService struct {
	deployments DeploymentStore
	runtime     ContainerRuntime
}

// NewDeploymentService creates a DeploymentService using deployments as its data store and runtime for container operations.
func NewDeploymentService(deployments DeploymentStore, runtime ContainerRuntime) *DeploymentService {
	return &DeploymentService{
		deployments: deployments,
		runtime:     runtime,
	}
}

// CreateDeployment validates input, persists initial pending deployment, triggers container creation/start, and updates status.
func (s *DeploymentService) CreateDeployment(ctx context.Context, deployment models.Deployment) (models.Deployment, error) {
	deployment.Image = strings.TrimSpace(deployment.Image)
	if deployment.Image == "" {
		return models.Deployment{}, ErrDeploymentImageRequired
	}

	deployment.Status = "pending"

	createdDeployment, err := s.deployments.CreateDeployment(ctx, deployment)
	if err != nil {
		return models.Deployment{}, err
	}

	if s.runtime != nil {
		containerID, err := s.runtime.CreateContainer(ctx, createdDeployment.Image)
		if err != nil {
			_, _ = s.deployments.UpdateDeploymentStatus(ctx, createdDeployment.ID, "failed")
			return models.Deployment{}, err
		}

		if err := s.runtime.StartContainer(ctx, containerID); err != nil {
			_, _ = s.deployments.UpdateDeploymentStatus(ctx, createdDeployment.ID, "failed")
			return models.Deployment{}, err
		}

		updatedDeployment, err := s.deployments.UpdateDeploymentStatus(ctx, createdDeployment.ID, "running")
		if err != nil {
			return createdDeployment, err
		}

		return updatedDeployment, nil
	}

	return createdDeployment, nil
}

// ListDeploymentsByProject returns deployments for projectID.
func (s *DeploymentService) ListDeploymentsByProject(ctx context.Context, projectID string) ([]models.Deployment, error) {
	return s.deployments.ListDeploymentsByProject(ctx, projectID)
}

// GetDeploymentByID returns the deployment identified by id.
func (s *DeploymentService) GetDeploymentByID(ctx context.Context, id string) (models.Deployment, error) {
	return s.deployments.GetDeploymentByID(ctx, id)
}
