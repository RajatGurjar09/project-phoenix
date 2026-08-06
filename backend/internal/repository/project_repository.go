package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/RajatGurjar09/project-phoenix/backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrProjectNotFound is returned when a project does not exist.
var ErrProjectNotFound = errors.New("project not found")

// ProjectRepository manages project persistence.
type ProjectRepository struct {
	db *pgxpool.Pool
}

// NewProjectRepository creates a ProjectRepository backed by db.
func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// CreateProject stores project and returns it with database-generated fields populated.
func (r *ProjectRepository) CreateProject(ctx context.Context, project models.Project) (models.Project, error) {
	const query = `
		INSERT INTO projects (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description, created_at, updated_at`

	createdProject, err := scanProject(r.db.QueryRow(ctx, query, project.Name, project.Description))
	if err != nil {
		return models.Project{}, fmt.Errorf("create project: %w", err)
	}

	return createdProject, nil
}

// GetProjectByID returns the project identified by id.
func (r *ProjectRepository) GetProjectByID(ctx context.Context, id string) (models.Project, error) {
	const query = `
		SELECT id, name, description, created_at, updated_at
		FROM projects
		WHERE id = $1`

	project, err := scanProject(r.db.QueryRow(ctx, query, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Project{}, fmt.Errorf("get project %q: %w", id, ErrProjectNotFound)
	}
	if err != nil {
		return models.Project{}, fmt.Errorf("get project %q: %w", id, err)
	}

	return project, nil
}

// ListProjects returns all projects ordered from oldest to newest.
func (r *ProjectRepository) ListProjects(ctx context.Context) ([]models.Project, error) {
	const query = `
		SELECT id, name, description, created_at, updated_at
		FROM projects
		ORDER BY created_at ASC, id ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	defer rows.Close()

	projects := make([]models.Project, 0)
	for rows.Next() {
		project, err := scanProject(rows)
		if err != nil {
			return nil, fmt.Errorf("list projects: %w", err)
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}

	return projects, nil
}

// DeleteProject removes the project identified by id.
func (r *ProjectRepository) DeleteProject(ctx context.Context, id string) error {
	const query = `DELETE FROM projects WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete project %q: %w", id, err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("delete project %q: %w", id, ErrProjectNotFound)
	}

	return nil
}

type projectScanner interface {
	Scan(dest ...any) error
}

func scanProject(row projectScanner) (models.Project, error) {
	var project models.Project
	if err := row.Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.CreatedAt,
		&project.UpdatedAt,
	); err != nil {
		return models.Project{}, err
	}

	return project, nil
}
