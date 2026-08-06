package docker

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestCreateContainerPullsImageBeforeCreatingContainer(t *testing.T) {
	var events []string
	pull := func(_ context.Context, imageName string) error {
		events = append(events, "pull:"+imageName)
		return nil
	}
	create := func(_ context.Context, imageName string) (string, error) {
		events = append(events, "create:"+imageName)
		return "container-id", nil
	}

	containerID, err := createContainer(context.Background(), " alpine:latest ", pull, create)
	if err != nil {
		t.Fatalf("createContainer() error = %v", err)
	}
	if containerID != "container-id" {
		t.Fatalf("createContainer() ID = %q, want %q", containerID, "container-id")
	}

	wantEvents := []string{"pull:alpine:latest", "create:alpine:latest"}
	if !reflect.DeepEqual(events, wantEvents) {
		t.Fatalf("events = %v, want %v", events, wantEvents)
	}
}

func TestCreateContainerDoesNotCreateWhenPullFails(t *testing.T) {
	pullErr := errors.New("pull failed")
	createCalled := false

	_, err := createContainer(
		context.Background(),
		"alpine:latest",
		func(context.Context, string) error { return pullErr },
		func(context.Context, string) (string, error) {
			createCalled = true
			return "", nil
		},
	)
	if !errors.Is(err, pullErr) {
		t.Fatalf("createContainer() error = %v, want wrapped %v", err, pullErr)
	}
	if createCalled {
		t.Fatal("container creation was called after image pull failed")
	}
}

func TestStopContainerValidatesID(t *testing.T) {
	err := stopContainer(context.Background(), "   ", func(_ context.Context, _ string) error {
		return nil
	})
	if err == nil {
		t.Fatal("stopContainer() expected error for empty ID, got nil")
	}
}

func TestStopContainerCallsStopFunction(t *testing.T) {
	var stoppedID string
	err := stopContainer(context.Background(), " container-123 ", func(_ context.Context, id string) error {
		stoppedID = id
		return nil
	})
	if err != nil {
		t.Fatalf("stopContainer() unexpected error = %v", err)
	}
	if stoppedID != "container-123" {
		t.Fatalf("stoppedID = %q, want %q", stoppedID, "container-123")
	}
}
