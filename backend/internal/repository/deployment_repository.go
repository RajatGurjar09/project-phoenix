package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/RajatGurjar09/project-phoenix/backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrDeploymentNotFound is returned when a deployment does not exist.
var ErrDeploymentNotFound = errors.New("deployment not found")

// DeploymentRepository manages deployment persistence.
type DeploymentRepository struct {
	db *pgxpool.Pool
}

// NewDeploymentRepository creates a DeploymentRepository backed by db.
func NewDeploymentRepository(db *pgxpool.Pool) *DeploymentRepository {
	return &DeploymentRepository{db: db}
}

// CreateDeployment stores deployment and returns it with database-generated fields populated.
func (r *DeploymentRepository) CreateDeployment(ctx context.Context, deployment models.Deployment) (models.Deployment, error) {
	const query = `
		INSERT INTO deployments (project_id, image, status, container_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, project_id, image, status, container_id, created_at, updated_at`

	createdDeployment, err := scanDeployment(r.db.QueryRow(
		ctx,
		query,
		deployment.ProjectID,
		deployment.Image,
		deployment.Status,
		deployment.ContainerID,
	))
	if err != nil {
		return models.Deployment{}, fmt.Errorf("create deployment: %w", err)
	}

	return createdDeployment, nil
}

// UpdateDeployment updates mutable fields (status, container_id) for the deployment identified by id.
func (r *DeploymentRepository) UpdateDeployment(ctx context.Context, id string, status string, containerID *string) (models.Deployment, error) {
	const query = `
		UPDATE deployments
		SET status = $1, container_id = COALESCE($2, container_id), updated_at = NOW()
		WHERE id = $3
		RETURNING id, project_id, image, status, container_id, created_at, updated_at`

	updatedDeployment, err := scanDeployment(r.db.QueryRow(ctx, query, status, containerID, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Deployment{}, fmt.Errorf("update deployment %q: %w", id, ErrDeploymentNotFound)
	}
	if err != nil {
		return models.Deployment{}, fmt.Errorf("update deployment %q: %w", id, err)
	}

	return updatedDeployment, nil
}

// UpdateDeploymentStatus updates the status and updated_at timestamp of the deployment identified by id.
func (r *DeploymentRepository) UpdateDeploymentStatus(ctx context.Context, id string, status string) (models.Deployment, error) {
	return r.UpdateDeployment(ctx, id, status, nil)
}

// ListDeploymentsByProject returns deployments for projectID from oldest to newest.
func (r *DeploymentRepository) ListDeploymentsByProject(ctx context.Context, projectID string) ([]models.Deployment, error) {
	const query = `
		SELECT id, project_id, image, status, container_id, created_at, updated_at
		FROM deployments
		WHERE project_id = $1
		ORDER BY created_at ASC, id ASC`

	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("list deployments for project %q: %w", projectID, err)
	}
	defer rows.Close()

	deployments := make([]models.Deployment, 0)
	for rows.Next() {
		deployment, err := scanDeployment(rows)
		if err != nil {
			return nil, fmt.Errorf("list deployments for project %q: %w", projectID, err)
		}
		deployments = append(deployments, deployment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list deployments for project %q: %w", projectID, err)
	}

	return deployments, nil
}

// GetDeploymentByID returns the deployment identified by id.
func (r *DeploymentRepository) GetDeploymentByID(ctx context.Context, id string) (models.Deployment, error) {
	const query = `
		SELECT id, project_id, image, status, container_id, created_at, updated_at
		FROM deployments
		WHERE id = $1`

	deployment, err := scanDeployment(r.db.QueryRow(ctx, query, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Deployment{}, fmt.Errorf("get deployment %q: %w", id, ErrDeploymentNotFound)
	}
	if err != nil {
		return models.Deployment{}, fmt.Errorf("get deployment %q: %w", id, err)
	}

	return deployment, nil
}

type deploymentScanner interface {
	Scan(dest ...any) error
}

func scanDeployment(row deploymentScanner) (models.Deployment, error) {
	var deployment models.Deployment
	if err := row.Scan(
		&deployment.ID,
		&deployment.ProjectID,
		&deployment.Image,
		&deployment.Status,
		&deployment.ContainerID,
		&deployment.CreatedAt,
		&deployment.UpdatedAt,
	); err != nil {
		return models.Deployment{}, err
	}

	return deployment, nil
}
