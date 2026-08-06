package docker

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/container"
)

// CreateContainer ensures image is available and creates a stopped Docker container.
func CreateContainer(ctx context.Context, imageName string) (string, error) {
	return createContainer(ctx, imageName, PullImage, createDockerContainer)
}

func createContainer(
	ctx context.Context,
	imageName string,
	pull func(context.Context, string) error,
	create func(context.Context, string) (string, error),
) (string, error) {
	imageName = strings.TrimSpace(imageName)
	if err := pull(ctx, imageName); err != nil {
		return "", fmt.Errorf("ensure image %q is available: %w", imageName, err)
	}

	containerID, err := create(ctx, imageName)
	if err != nil {
		return "", fmt.Errorf("create container from image %q: %w", imageName, err)
	}

	return containerID, nil
}

func createDockerContainer(ctx context.Context, imageName string) (string, error) {
	client, err := NewClient()
	if err != nil {
		return "", fmt.Errorf("create docker client: %w", err)
	}
	defer client.Close()

	response, err := client.ContainerCreate(ctx, &container.Config{Image: imageName}, nil, nil, nil, "")
	if err != nil {
		return "", err
	}

	return response.ID, nil
}

// StartContainer starts an existing Docker container.
func StartContainer(ctx context.Context, containerID string) error {
	return startContainer(ctx, containerID, startDockerContainer)
}

func startContainer(
	ctx context.Context,
	containerID string,
	start func(context.Context, string) error,
) error {
	containerID = strings.TrimSpace(containerID)

	if containerID == "" {
		return fmt.Errorf("container id is required")
	}

	if err := start(ctx, containerID); err != nil {
		return fmt.Errorf("start container %q: %w", containerID, err)
	}

	return nil
}

func startDockerContainer(ctx context.Context, containerID string) error {
	client, err := NewClient()
	if err != nil {
		return fmt.Errorf("create docker client: %w", err)
	}
	defer client.Close()

	return client.ContainerStart(
		ctx,
		containerID,
		container.StartOptions{},
	)
}

// StopContainer stops a running Docker container.
func StopContainer(ctx context.Context, containerID string) error {
	return stopContainer(ctx, containerID, stopDockerContainer)
}

func stopContainer(
	ctx context.Context,
	containerID string,
	stop func(context.Context, string) error,
) error {
	containerID = strings.TrimSpace(containerID)

	if containerID == "" {
		return fmt.Errorf("container id is required")
	}

	if err := stop(ctx, containerID); err != nil {
		return fmt.Errorf("stop container %q: %w", containerID, err)
	}

	return nil
}

func stopDockerContainer(ctx context.Context, containerID string) error {
	client, err := NewClient()
	if err != nil {
		return fmt.Errorf("create docker client: %w", err)
	}
	defer client.Close()

	return client.ContainerStop(
		ctx,
		containerID,
		container.StopOptions{},
	)
}
