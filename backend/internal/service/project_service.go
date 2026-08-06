package service

import (
	"context"
	"errors"
	"strings"

	"github.com/RajatGurjar09/project-phoenix/backend/internal/models"
)

// ErrProjectNameRequired is returned when a project has no name.
var ErrProjectNameRequired = errors.New("project name is required")

// ProjectStore defines the persistence operations needed by ProjectService.
type ProjectStore interface {
	CreateProject(ctx context.Context, project models.Project) (models.Project, error)
	GetProjectByID(ctx context.Context, id string) (models.Project, error)
	UpdateProject(ctx context.Context, id string, project models.Project) (models.Project, error)
	ListProjects(ctx context.Context) ([]models.Project, error)
	DeleteProject(ctx context.Context, id string) error
}

// ProjectService contains project domain operations.
type ProjectService struct {
	projects ProjectStore
}

// NewProjectService creates a ProjectService using projects as its data store.
func NewProjectService(projects ProjectStore) *ProjectService {
	return &ProjectService{projects: projects}
}

// CreateProject validates and persists project.
func (s *ProjectService) CreateProject(ctx context.Context, project models.Project) (models.Project, error) {
	project.Name = strings.TrimSpace(project.Name)
	if project.Name == "" {
		return models.Project{}, ErrProjectNameRequired
	}

	return s.projects.CreateProject(ctx, project)
}

// GetProjectByID returns the project identified by id.
func (s *ProjectService) GetProjectByID(ctx context.Context, id string) (models.Project, error) {
	return s.projects.GetProjectByID(ctx, id)
}

// UpdateProject validates and updates the project identified by id.
func (s *ProjectService) UpdateProject(ctx context.Context, id string, project models.Project) (models.Project, error) {
	project.Name = strings.TrimSpace(project.Name)
	if project.Name == "" {
		return models.Project{}, ErrProjectNameRequired
	}

	return s.projects.UpdateProject(ctx, id, project)
}

// ListProjects returns all projects.
func (s *ProjectService) ListProjects(ctx context.Context) ([]models.Project, error) {
	return s.projects.ListProjects(ctx)
}

// DeleteProject removes the project identified by id.
func (s *ProjectService) DeleteProject(ctx context.Context, id string) error {
	return s.projects.DeleteProject(ctx, id)
}
