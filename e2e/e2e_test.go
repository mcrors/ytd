//go:build e2e

package e2e

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"
)

func baseURL() string {
	if u := os.Getenv("YTD_BASE_URL"); u != "" {
		return u
	}
	return "http://localhost:8080"
}

func TestHealthz(t *testing.T) {
	resp, err := http.Get(baseURL() + "/healthz")
	if err != nil {
		t.Fatalf("GET /healthz: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", resp.StatusCode)
	}

	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decoding body: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("status: got %q, want \"ok\"", body["status"])
	}
}

func TestReadyz_DegradedWithoutYtDlp(t *testing.T) {
	resp, err := http.Get(baseURL() + "/readyz")
	if err != nil {
		t.Fatalf("GET /readyz: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("status: got %d, want 503", resp.StatusCode)
	}

	var body struct {
		Status string            `json:"status"`
		Checks map[string]string `json:"checks"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decoding body: %v", err)
	}
	if body.Status != "degraded" {
		t.Errorf("status: got %q, want \"degraded\"", body.Status)
	}
	if body.Checks["baseDir"] != "ok" {
		t.Errorf("baseDir check: got %q, want \"ok\"", body.Checks["baseDir"])
	}
	if body.Checks["db"] != "ok" {
		t.Errorf("db check: got %q, want \"ok\"", body.Checks["db"])
	}
}

func TestGetDirectories(t *testing.T) {
	resp, err := http.Get(baseURL() + "/api/directories")
	if err != nil {
		t.Fatalf("GET /api/directories: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", resp.StatusCode)
	}

	var body struct {
		Directories []string `json:"directories"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decoding body: %v", err)
	}
}

func TestCreateDirectory(t *testing.T) {
	resp, err := http.Post(
		baseURL()+"/api/directory",
		"application/json",
		strings.NewReader(`{"dir":"history/mary-beard"}`),
	)
	if err != nil {
		t.Fatalf("POST /api/directory: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("status: got %d, want 201", resp.StatusCode)
	}
}

func TestCreateDirectory_TraversalRejected(t *testing.T) {
	resp, err := http.Post(
		baseURL()+"/api/directory",
		"application/json",
		strings.NewReader(`{"dir":"../../etc/passwd"}`),
	)
	if err != nil {
		t.Fatalf("POST /api/directory: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", resp.StatusCode)
	}
}
