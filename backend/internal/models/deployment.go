package models

import "time"

// Deployment represents a container image deployment for a project.
type Deployment struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Image     string    `json:"image"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
