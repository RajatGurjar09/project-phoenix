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
