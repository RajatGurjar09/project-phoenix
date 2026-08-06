package docker

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/image"
)

func TestPullImageConsumesAndClosesStream(t *testing.T) {
	stream := &trackingReadCloser{reader: strings.NewReader("pull complete")}
	client := imagePullClientFunc(func(_ context.Context, ref string, _ image.PullOptions) (io.ReadCloser, error) {
		if ref != "alpine:latest" {
			t.Fatalf("ImagePull() ref = %q, want %q", ref, "alpine:latest")
		}

		return stream, nil
	})

	if err := pullImage(context.Background(), client, "alpine:latest"); err != nil {
		t.Fatalf("pullImage() error = %v", err)
	}
	if !stream.reachedEOF {
		t.Fatal("pull stream was not consumed to EOF")
	}
	if !stream.closed {
		t.Fatal("pull stream was not closed")
	}
}

func TestPullImageReturnsPullError(t *testing.T) {
	pullErr := errors.New("docker unavailable")
	client := imagePullClientFunc(func(context.Context, string, image.PullOptions) (io.ReadCloser, error) {
		return nil, pullErr
	})

	err := pullImage(context.Background(), client, "alpine:latest")
	if !errors.Is(err, pullErr) {
		t.Fatalf("pullImage() error = %v, want wrapped %v", err, pullErr)
	}
}

type imagePullClientFunc func(context.Context, string, image.PullOptions) (io.ReadCloser, error)

func (f imagePullClientFunc) ImagePull(ctx context.Context, ref string, options image.PullOptions) (io.ReadCloser, error) {
	return f(ctx, ref, options)
}

type trackingReadCloser struct {
	reader     *strings.Reader
	closed     bool
	reachedEOF bool
}

func (r *trackingReadCloser) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	if errors.Is(err, io.EOF) {
		r.reachedEOF = true
	}

	return n, err
}

func (r *trackingReadCloser) Close() error {
	r.closed = true
	return nil
}
