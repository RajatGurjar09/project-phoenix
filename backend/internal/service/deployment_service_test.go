package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/RajatGurjar09/project-phoenix/backend/internal/models"
)

type mockDeploymentStore struct {
	created []models.Deployment
	updated map[string]string
}

func newMockDeploymentStore() *mockDeploymentStore {
	return &mockDeploymentStore{
		created: make([]models.Deployment, 0),
		updated: make(map[string]string),
	}
}

func (m *mockDeploymentStore) CreateDeployment(_ context.Context, d models.Deployment) (models.Deployment, error) {
	d.ID = "dep-123"
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
	m.created = append(m.created, d)
	return d, nil
}

func (m *mockDeploymentStore) UpdateDeployment(_ context.Context, id string, status string, containerID *string) (models.Deployment, error) {
	m.updated[id] = status
	for i, d := range m.created {
		if d.ID == id {
			d.Status = status
			if containerID != nil {
				d.ContainerID = containerID
			}
			d.UpdatedAt = time.Now()
			m.created[i] = d
			return d, nil
		}
	}
	return models.Deployment{}, errors.New("not found")
}

func (m *mockDeploymentStore) UpdateDeploymentStatus(ctx context.Context, id string, status string) (models.Deployment, error) {
	return m.UpdateDeployment(ctx, id, status, nil)
}

func (m *mockDeploymentStore) ListDeploymentsByProject(_ context.Context, _ string) ([]models.Deployment, error) {
	return m.created, nil
}

func (m *mockDeploymentStore) GetDeploymentByID(_ context.Context, id string) (models.Deployment, error) {
	for _, d := range m.created {
		if d.ID == id {
			return d, nil
		}
	}
	return models.Deployment{}, errors.New("not found")
}

type mockContainerRuntime struct {
	createFn func(ctx context.Context, imageName string) (string, error)
	startFn  func(ctx context.Context, containerID string) error
	stopFn   func(ctx context.Context, containerID string) error
}

func (m *mockContainerRuntime) CreateContainer(ctx context.Context, imageName string) (string, error) {
	if m.createFn != nil {
		return m.createFn(ctx, imageName)
	}
	return "container-123", nil
}

func (m *mockContainerRuntime) StartContainer(ctx context.Context, containerID string) error {
	if m.startFn != nil {
		return m.startFn(ctx, containerID)
	}
	return nil
}

func (m *mockContainerRuntime) StopContainer(ctx context.Context, containerID string) error {
	if m.stopFn != nil {
		return m.stopFn(ctx, containerID)
	}
	return nil
}

func TestCreateDeploymentSuccess(t *testing.T) {
	store := newMockDeploymentStore()
	runtime := &mockContainerRuntime{}
	svc := NewDeploymentService(store, runtime)

	d, err := svc.CreateDeployment(context.Background(), models.Deployment{
		ProjectID: "proj-1",
		Image:     "nginx:latest",
	})
	if err != nil {
		t.Fatalf("CreateDeployment() unexpected error = %v", err)
	}

	if d.Status != "running" {
		t.Errorf("CreateDeployment() status = %q, want %q", d.Status, "running")
	}
	if d.ContainerID == nil || *d.ContainerID != "container-123" {
		t.Errorf("CreateDeployment() containerID = %v, want %q", d.ContainerID, "container-123")
	}

	if store.updated["dep-123"] != "running" {
		t.Errorf("store updated status = %q, want %q", store.updated["dep-123"], "running")
	}
}

func TestCreateDeploymentImageRequired(t *testing.T) {
	store := newMockDeploymentStore()
	svc := NewDeploymentService(store, nil)

	_, err := svc.CreateDeployment(context.Background(), models.Deployment{
		ProjectID: "proj-1",
		Image:     "   ",
	})
	if !errors.Is(err, ErrDeploymentImageRequired) {
		t.Errorf("CreateDeployment() error = %v, want %v", err, ErrDeploymentImageRequired)
	}
}

func TestCreateDeploymentContainerCreateError(t *testing.T) {
	store := newMockDeploymentStore()
	createErr := errors.New("docker pull failed")
	runtime := &mockContainerRuntime{
		createFn: func(_ context.Context, _ string) (string, error) {
			return "", createErr
		},
	}
	svc := NewDeploymentService(store, runtime)

	_, err := svc.CreateDeployment(context.Background(), models.Deployment{
		ProjectID: "proj-1",
		Image:     "invalid:image",
	})
	if !errors.Is(err, createErr) {
		t.Fatalf("CreateDeployment() error = %v, want %v", err, createErr)
	}

	if store.updated["dep-123"] != "failed" {
		t.Errorf("store updated status = %q, want %q", store.updated["dep-123"], "failed")
	}
}

func TestCreateDeploymentContainerStartError(t *testing.T) {
	store := newMockDeploymentStore()
	startErr := errors.New("container start timeout")
	runtime := &mockContainerRuntime{
		startFn: func(_ context.Context, _ string) error {
			return startErr
		},
	}
	svc := NewDeploymentService(store, runtime)

	_, err := svc.CreateDeployment(context.Background(), models.Deployment{
		ProjectID: "proj-1",
		Image:     "nginx:latest",
	})
	if !errors.Is(err, startErr) {
		t.Fatalf("CreateDeployment() error = %v, want %v", err, startErr)
	}

	if store.updated["dep-123"] != "failed" {
		t.Errorf("store updated status = %q, want %q", store.updated["dep-123"], "failed")
	}
}

func TestStopDeploymentSuccess(t *testing.T) {
	store := newMockDeploymentStore()
	runtime := &mockContainerRuntime{}
	svc := NewDeploymentService(store, runtime)

	created, err := svc.CreateDeployment(context.Background(), models.Deployment{
		ProjectID: "proj-1",
		Image:     "nginx:latest",
	})
	if err != nil {
		t.Fatalf("CreateDeployment() unexpected error = %v", err)
	}

	stopped, err := svc.StopDeployment(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("StopDeployment() unexpected error = %v", err)
	}

	if stopped.Status != "stopped" {
		t.Errorf("StopDeployment() status = %q, want %q", stopped.Status, "stopped")
	}
}

func TestStopDeploymentNoContainerID(t *testing.T) {
	store := newMockDeploymentStore()
	svc := NewDeploymentService(store, nil)

	dep, _ := store.CreateDeployment(context.Background(), models.Deployment{
		ProjectID: "proj-1",
		Image:     "nginx:latest",
		Status:    "running",
	})

	_, err := svc.StopDeployment(context.Background(), dep.ID)
	if !errors.Is(err, ErrDeploymentContainerNotFound) {
		t.Errorf("StopDeployment() error = %v, want %v", err, ErrDeploymentContainerNotFound)
	}
}
