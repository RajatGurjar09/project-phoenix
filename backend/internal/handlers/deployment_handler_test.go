package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/RajatGurjar09/project-phoenix/backend/internal/models"
	"github.com/RajatGurjar09/project-phoenix/backend/internal/service"
)

type mockDeploymentHandlerService struct {
	createFn  func(context.Context, models.Deployment) (models.Deployment, error)
	listFn    func(context.Context, string) ([]models.Deployment, error)
	getFn     func(context.Context, string) (models.Deployment, error)
	stopFn    func(context.Context, string) (models.Deployment, error)
	restartFn func(context.Context, string) (models.Deployment, error)
	removeFn  func(context.Context, string) (models.Deployment, error)
}

func (m *mockDeploymentHandlerService) CreateDeployment(ctx context.Context, deployment models.Deployment) (models.Deployment, error) {
	if m.createFn != nil {
		return m.createFn(ctx, deployment)
	}
	return deployment, nil
}

func (m *mockDeploymentHandlerService) ListDeploymentsByProject(ctx context.Context, projectID string) ([]models.Deployment, error) {
	if m.listFn != nil {
		return m.listFn(ctx, projectID)
	}
	return nil, nil
}

func (m *mockDeploymentHandlerService) GetDeploymentByID(ctx context.Context, id string) (models.Deployment, error) {
	if m.getFn != nil {
		return m.getFn(ctx, id)
	}
	return models.Deployment{}, nil
}

func (m *mockDeploymentHandlerService) StopDeployment(ctx context.Context, id string) (models.Deployment, error) {
	if m.stopFn != nil {
		return m.stopFn(ctx, id)
	}
	return models.Deployment{}, nil
}

func (m *mockDeploymentHandlerService) RestartDeployment(ctx context.Context, id string) (models.Deployment, error) {
	if m.restartFn != nil {
		return m.restartFn(ctx, id)
	}
	return models.Deployment{}, nil
}

func (m *mockDeploymentHandlerService) RemoveDeployment(ctx context.Context, id string) (models.Deployment, error) {
	if m.removeFn != nil {
		return m.removeFn(ctx, id)
	}
	return models.Deployment{}, nil
}

func newDeploymentTestHandler(svc DeploymentService) *Handler {
	return New(
		"test",
		time.Now(),
		nil,
		svc,
	)
}

func TestCreateDeployment(t *testing.T) {
	var received models.Deployment

	svc := &mockDeploymentHandlerService{
		createFn: func(_ context.Context, deployment models.Deployment) (models.Deployment, error) {
			received = deployment
			return models.Deployment{
				ID:        "dep-123",
				ProjectID: deployment.ProjectID,
				Image:     deployment.Image,
				Status:    "running",
			}, nil
		},
	}

	handler := newDeploymentTestHandler(svc)

	req := httptest.NewRequest(
		http.MethodPost,
		"/projects/proj-1/deployments",
		strings.NewReader("{\"image\":\"nginx:latest\",\"status\":\"pending\"}"),
	)
	req.SetPathValue("id", "proj-1")

	rec := httptest.NewRecorder()

	handler.CreateDeployment(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("CreateDeployment() status = %d, want %d", rec.Code, http.StatusCreated)
	}

	if received.ProjectID != "proj-1" {
		t.Errorf("ProjectID = %q, want %q", received.ProjectID, "proj-1")
	}

	if received.Image != "nginx:latest" {
		t.Errorf("Image = %q, want %q", received.Image, "nginx:latest")
	}

}

func TestCreateDeploymentInvalidJSON(t *testing.T) {
	handler := newDeploymentTestHandler(&mockDeploymentHandlerService{})

	req := httptest.NewRequest(
		http.MethodPost,
		"/projects/proj-1/deployments",
		strings.NewReader("{\"image\":"),
	)
	req.SetPathValue("id", "proj-1")

	rec := httptest.NewRecorder()

	handler.CreateDeployment(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("CreateDeployment() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}

}

func TestCreateDeploymentServiceError(t *testing.T) {
	expectedErr := errors.New("create failed")

	svc := &mockDeploymentHandlerService{
		createFn: func(_ context.Context, _ models.Deployment) (models.Deployment, error) {
			return models.Deployment{}, expectedErr
		},
	}

	handler := newDeploymentTestHandler(svc)

	req := httptest.NewRequest(
		http.MethodPost,
		"/projects/proj-1/deployments",
		strings.NewReader("{\"image\":\"nginx:latest\"}"),
	)
	req.SetPathValue("id", "proj-1")

	rec := httptest.NewRecorder()

	handler.CreateDeployment(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("CreateDeployment() status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}

}

func TestListDeploymentsByProject(t *testing.T) {
	svc := &mockDeploymentHandlerService{
		listFn: func(_ context.Context, projectID string) ([]models.Deployment, error) {
			if projectID != "proj-1" {
				t.Errorf("projectID = %q, want %q", projectID, "proj-1")
			}

			return []models.Deployment{
				{
					ID:        "dep-123",
					ProjectID: "proj-1",
					Image:     "nginx:latest",
					Status:    "running",
				},
			}, nil
		},
	}

	handler := newDeploymentTestHandler(svc)

	req := httptest.NewRequest(
		http.MethodGet,
		"/projects/proj-1/deployments",
		nil,
	)
	req.SetPathValue("id", "proj-1")

	rec := httptest.NewRecorder()

	handler.ListDeploymentsByProject(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("ListDeploymentsByProject() status = %d, want %d", rec.Code, http.StatusOK)
	}

}

func TestGetDeployment(t *testing.T) {
	svc := &mockDeploymentHandlerService{
		getFn: func(_ context.Context, id string) (models.Deployment, error) {
			if id != "dep-123" {
				t.Errorf("id = %q, want %q", id, "dep-123")
			}

			return models.Deployment{
				ID:        "dep-123",
				ProjectID: "proj-1",
				Image:     "nginx:latest",
				Status:    "running",
			}, nil
		},
	}

	handler := newDeploymentTestHandler(svc)

	req := httptest.NewRequest(
		http.MethodGet,
		"/deployments/dep-123",
		nil,
	)
	req.SetPathValue("id", "dep-123")

	rec := httptest.NewRecorder()

	handler.GetDeployment(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GetDeployment() status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestStopDeployment(t *testing.T) {
	var receivedID string

	svc := &mockDeploymentHandlerService{
		stopFn: func(_ context.Context, id string) (models.Deployment, error) {
			receivedID = id

			return models.Deployment{
				ID:        id,
				ProjectID: "proj-1",
				Image:     "nginx:latest",
				Status:    "stopped",
			}, nil
		},
	}

	handler := newDeploymentTestHandler(svc)

	req := httptest.NewRequest(
		http.MethodPost,
		"/deployments/dep-123/stop",
		nil,
	)
	req.SetPathValue("id", "dep-123")

	rec := httptest.NewRecorder()

	handler.StopDeployment(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("StopDeployment() status = %d, want %d", rec.Code, http.StatusOK)
	}

	if receivedID != "dep-123" {
		t.Errorf("id = %q, want %q", receivedID, "dep-123")
	}
}

func TestStopDeploymentServiceError(t *testing.T) {
	svc := &mockDeploymentHandlerService{
		stopFn: func(_ context.Context, _ string) (models.Deployment, error) {
			return models.Deployment{}, service.ErrInvalidDeploymentState
		},
	}

	handler := newDeploymentTestHandler(svc)

	req := httptest.NewRequest(
		http.MethodPost,
		"/deployments/dep-123/stop",
		nil,
	)
	req.SetPathValue("id", "dep-123")

	rec := httptest.NewRecorder()

	handler.StopDeployment(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("StopDeployment() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestRestartDeployment(t *testing.T) {
	var receivedID string

	svc := &mockDeploymentHandlerService{
		restartFn: func(_ context.Context, id string) (models.Deployment, error) {
			receivedID = id

			return models.Deployment{
				ID:        id,
				ProjectID: "proj-1",
				Image:     "nginx:latest",
				Status:    "running",
			}, nil
		},
	}

	handler := newDeploymentTestHandler(svc)

	req := httptest.NewRequest(
		http.MethodPost,
		"/deployments/dep-123/restart",
		nil,
	)
	req.SetPathValue("id", "dep-123")

	rec := httptest.NewRecorder()

	handler.RestartDeployment(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("RestartDeployment() status = %d, want %d", rec.Code, http.StatusOK)
	}

	if receivedID != "dep-123" {
		t.Errorf("id = %q, want %q", receivedID, "dep-123")
	}
}

func TestRemoveDeployment(t *testing.T) {
	var receivedID string

	svc := &mockDeploymentHandlerService{
		removeFn: func(_ context.Context, id string) (models.Deployment, error) {
			receivedID = id

			return models.Deployment{
				ID:        id,
				ProjectID: "proj-1",
				Image:     "nginx:latest",
				Status:    "removed",
			}, nil
		},
	}

	handler := newDeploymentTestHandler(svc)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/deployments/dep-123",
		nil,
	)
	req.SetPathValue("id", "dep-123")

	rec := httptest.NewRecorder()

	handler.RemoveDeployment(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("RemoveDeployment() status = %d, want %d", rec.Code, http.StatusOK)
	}

	if receivedID != "dep-123" {
		t.Errorf("id = %q, want %q", receivedID, "dep-123")
	}
}
