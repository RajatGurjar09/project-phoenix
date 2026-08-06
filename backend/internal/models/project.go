package models

import "time"

// Project represents a deployable project managed by the platform.
type Project struct {
	ID          string
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
