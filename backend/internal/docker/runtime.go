package docker

import "context"

// Runtime implements ContainerRuntime using Docker Engine API primitives.
type Runtime struct{}

// NewRuntime creates a new Docker Runtime instance.
func NewRuntime() *Runtime {
	return &Runtime{}
}

// CreateContainer pulls imageName if necessary and creates a stopped container.
func (r *Runtime) CreateContainer(ctx context.Context, imageName string) (string, error) {
	return CreateContainer(ctx, imageName)
}

// StartContainer starts an existing container identified by containerID.
func (r *Runtime) StartContainer(ctx context.Context, containerID string) error {
	return StartContainer(ctx, containerID)
}
