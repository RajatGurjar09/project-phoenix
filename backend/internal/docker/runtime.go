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

// StopContainer stops a running container identified by containerID.
func (r *Runtime) StopContainer(ctx context.Context, containerID string) error {
	return StopContainer(ctx, containerID)
}

// RestartContainer restarts an existing container identified by containerID.
func (r *Runtime) RestartContainer(ctx context.Context, containerID string) error {
	return RestartContainer(ctx, containerID)
}

// RemoveContainer force-removes an existing container identified by containerID.
func (r *Runtime) RemoveContainer(ctx context.Context, containerID string) error {
	return RemoveContainer(ctx, containerID)
}
