package docker

import (
	"context"
	"fmt"
	"io"
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

// RestartContainer restarts an existing Docker container.
func RestartContainer(ctx context.Context, containerID string) error {
	return restartContainer(ctx, containerID, restartDockerContainer)
}

func restartContainer(
	ctx context.Context,
	containerID string,
	restart func(context.Context, string) error,
) error {
	containerID = strings.TrimSpace(containerID)

	if containerID == "" {
		return fmt.Errorf("container id is required")
	}

	if err := restart(ctx, containerID); err != nil {
		return fmt.Errorf("restart container %q: %w", containerID, err)
	}

	return nil
}

func restartDockerContainer(ctx context.Context, containerID string) error {
	client, err := NewClient()
	if err != nil {
		return fmt.Errorf("create docker client: %w", err)
	}
	defer client.Close()

	return client.ContainerRestart(
		ctx,
		containerID,
		container.StopOptions{},
	)
}

// RemoveContainer force-removes an existing Docker container.
func RemoveContainer(ctx context.Context, containerID string) error {
	return removeContainer(ctx, containerID, removeDockerContainer)
}

func removeContainer(
	ctx context.Context,
	containerID string,
	remove func(context.Context, string) error,
) error {
	containerID = strings.TrimSpace(containerID)

	if containerID == "" {
		return fmt.Errorf("container id is required")
	}

	if err := remove(ctx, containerID); err != nil {
		return fmt.Errorf("remove container %q: %w", containerID, err)
	}

	return nil
}

func removeDockerContainer(ctx context.Context, containerID string) error {
	client, err := NewClient()
	if err != nil {
		return fmt.Errorf("create docker client: %w", err)
	}
	defer client.Close()

	return client.ContainerRemove(
		ctx,
		containerID,
		container.RemoveOptions{Force: true},
	)
}

// ContainerLogs returns the recent logs for an existing Docker container.
func ContainerLogs(ctx context.Context, containerID string) (string, error) {
	return containerLogs(ctx, containerID, getDockerContainerLogs)
}

func containerLogs(
	ctx context.Context,
	containerID string,
	getLogs func(context.Context, string) (string, error),
) (string, error) {
	containerID = strings.TrimSpace(containerID)

	if containerID == "" {
		return "", fmt.Errorf("container id is required")
	}

	logs, err := getLogs(ctx, containerID)
	if err != nil {
		return "", fmt.Errorf("get logs for container %q: %w", containerID, err)
	}

	return logs, nil
}

func getDockerContainerLogs(ctx context.Context, containerID string) (string, error) {
	client, err := NewClient()
	if err != nil {
		return "", fmt.Errorf("create docker client: %w", err)
	}
	defer client.Close()

	reader, err := client.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Tail:       "100",
	})
	if err != nil {
		return "", err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
