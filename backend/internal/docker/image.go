package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types/image"
)

var errImageRequired = errors.New("image is required")

type imagePullClient interface {
	ImagePull(ctx context.Context, ref string, options image.PullOptions) (io.ReadCloser, error)
}

// PullImage pulls image through the Docker Engine API.
// Docker returns a successful pull stream when the image already exists locally.
func PullImage(ctx context.Context, imageName string) error {
	imageName = strings.TrimSpace(imageName)
	if imageName == "" {
		return errImageRequired
	}

	client, err := NewClient()
	if err != nil {
		return fmt.Errorf("create docker client: %w", err)
	}
	defer client.Close()

	return pullImage(ctx, client, imageName)
}

func pullImage(ctx context.Context, client imagePullClient, imageName string) error {
	stream, err := client.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("pull image %q: %w", imageName, err)
	}
	defer stream.Close()

	if _, err := io.Copy(io.Discard, stream); err != nil {
		return fmt.Errorf("read image pull stream for %q: %w", imageName, err)
	}

	return nil
}
