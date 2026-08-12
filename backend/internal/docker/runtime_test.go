package docker

import (
	"context"
	"testing"
)

func TestNewRuntime(t *testing.T) {
	runtime := NewRuntime()

	if runtime == nil {
		t.Fatal("NewRuntime() returned nil")
	}
}

func TestRuntimeImplementsContainerRuntime(t *testing.T) {
	var _ interface {
		CreateContainer(context.Context, string) (string, error)
		StartContainer(context.Context, string) error
		StopContainer(context.Context, string) error
		RestartContainer(context.Context, string) error
		RemoveContainer(context.Context, string) error
	} = (*Runtime)(nil)
}
