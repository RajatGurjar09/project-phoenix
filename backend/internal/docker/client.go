// Package docker provides integration points for the Docker Engine API.
package docker

import "github.com/docker/docker/client"

// NewClient creates a Docker Engine client from the standard Docker environment variables.
// API version negotiation is deferred until the first API request.
func NewClient() (*client.Client, error) {
	return client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
}
