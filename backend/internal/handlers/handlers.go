package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type Handler struct {
	version   string
	startedAt time.Time
}

func New(version string, startedAt time.Time) *Handler {
	return &Handler{version: version, startedAt: startedAt}
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status":    "ok",
		"version":   h.version,
		"uptime":    time.Since(h.startedAt).String(),
		"hostname":  hostname,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *Handler) Version(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"version": h.version})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
