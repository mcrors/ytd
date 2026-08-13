package web_test

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/mcrors/ytd/internal/db"
	"github.com/mcrors/ytd/internal/format"
	"github.com/mcrors/ytd/internal/web"
)

// --- Test doubles ---

type mockQueue struct {
	getTitle  func(context.Context, string) (string, error)
	enqueue   func(context.Context, string, string, string, format.Format, string) (int64, error)
	cancelled []int64
}

func (m *mockQueue) GetTitle(ctx context.Context, u string) (string, error) {
	if m.getTitle != nil {
		return m.getTitle(ctx, u)
	}
	return "Test Video", nil
}

func (m *mockQueue) Enqueue(ctx context.Context, u, targetDir, newName string, f format.Format, title string) (int64, error) {
	if m.enqueue != nil {
		return m.enqueue(ctx, u, targetDir, newName, f, title)
	}
	return 0, fmt.Errorf("enqueue not configured")
}

func (m *mockQueue) Cancel(id int64) {
	m.cancelled = append(m.cancelled, id)
}

// --- Helpers ---

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	database.SetMaxOpenConns(1)
	if err := db.Migrate(database); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() { database.Close() })
	return database
}

func newTestServer(t *testing.T, q web.Queue, database *sql.DB) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	if err := web.RegisterRoutes(mux, q, t.TempDir(), database, false); err != nil {
		t.Fatalf("RegisterRoutes: %v", err)
	}
	return httptest.NewServer(mux)
}

func insertDownload(t *testing.T, database *sql.DB, status, title string, progress int) int64 {
	t.Helper()
	res, err := database.Exec(
		`INSERT INTO downloads (url, target_dir, format, status, title, progress) VALUES (?, ?, ?, ?, ?, ?)`,
		"https://example.com/video", "/tmp/media", "best", status, title, progress,
	)
	if err != nil {
		t.Fatalf("insertDownload: %v", err)
	}
	id, _ := res.LastInsertId()
	return id
}

func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return string(b)
}

// --- Tests ---

func TestSubmitHandler(t *testing.T) {
	database := newTestDB(t)

	q := &mockQueue{
		getTitle: func(_ context.Context, _ string) (string, error) { return "My Cool Video", nil },
		enqueue: func(_ context.Context, u, targetDir, _ string, f format.Format, title string) (int64, error) {
			res, err := database.Exec(
				`INSERT INTO downloads (url, target_dir, format, status, title) VALUES (?, ?, ?, 'queued', ?)`,
				u, targetDir, string(f), title,
			)
			if err != nil {
				return 0, err
			}
			return res.LastInsertId()
		},
	}

	srv := newTestServer(t, q, database)
	defer srv.Close()

	form := url.Values{"url": {"https://example.com/video"}, "format": {"best"}}
	resp, err := srv.Client().Post(srv.URL+"/downloads", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("POST /downloads: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	body := readBody(t, resp)
	if !strings.Contains(body, "My Cool Video") {
		t.Errorf("body should contain video title, got: %s", body)
	}
	if !strings.Contains(body, "queued") {
		t.Errorf("body should contain status 'queued', got: %s", body)
	}
}

func TestSubmitHandler_MissingURL(t *testing.T) {
	database := newTestDB(t)
	srv := newTestServer(t, &mockQueue{}, database)
	defer srv.Close()

	form := url.Values{"format": {"best"}}
	resp, err := srv.Client().Post(srv.URL+"/downloads", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("POST /downloads: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", resp.StatusCode)
	}
}

func TestStatusHandler(t *testing.T) {
	database := newTestDB(t)
	id := insertDownload(t, database, "downloading", "Test Video", 45)

	srv := newTestServer(t, &mockQueue{}, database)
	defer srv.Close()

	resp, err := srv.Client().Get(fmt.Sprintf("%s/downloads/%d/status", srv.URL, id))
	if err != nil {
		t.Fatalf("GET /downloads/%d/status: %v", id, err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	body := readBody(t, resp)
	if !strings.Contains(body, "45") {
		t.Errorf("body should contain progress 45, got: %s", body)
	}
}

func TestStatusHandler_NotFound(t *testing.T) {
	database := newTestDB(t)
	srv := newTestServer(t, &mockQueue{}, database)
	defer srv.Close()

	resp, err := srv.Client().Get(srv.URL + "/downloads/9999/status")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestCancelHandler(t *testing.T) {
	database := newTestDB(t)
	id := insertDownload(t, database, "downloading", "Test Video", 20)
	mq := &mockQueue{}

	srv := newTestServer(t, mq, database)
	defer srv.Close()

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/downloads/%d/cancel", srv.URL, id), nil)
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("DELETE: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("status = %d, want 204", resp.StatusCode)
	}
	resp.Body.Close()

	if len(mq.cancelled) == 0 || mq.cancelled[0] != id {
		t.Errorf("Cancel not called with id %d, got: %v", id, mq.cancelled)
	}
}

func TestHistoryHandler(t *testing.T) {
	database := newTestDB(t)
	insertDownload(t, database, "completed", "Finished Video", 100)
	insertDownload(t, database, "failed", "Broken Video", 0)
	insertDownload(t, database, "downloading", "Active Video", 50) // should not appear

	srv := newTestServer(t, &mockQueue{}, database)
	defer srv.Close()

	resp, err := srv.Client().Get(srv.URL + "/downloads/history")
	if err != nil {
		t.Fatalf("GET /downloads/history: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	body := readBody(t, resp)
	if !strings.Contains(body, "Finished Video") {
		t.Errorf("body should contain completed download, got: %s", body)
	}
	if !strings.Contains(body, "Broken Video") {
		t.Errorf("body should contain failed download, got: %s", body)
	}
	if strings.Contains(body, "Active Video") {
		t.Errorf("body should not contain active download, got: %s", body)
	}
}

func TestHistoryHandler_Empty(t *testing.T) {
	database := newTestDB(t)
	srv := newTestServer(t, &mockQueue{}, database)
	defer srv.Close()

	resp, err := srv.Client().Get(srv.URL + "/downloads/history")
	if err != nil {
		t.Fatalf("GET /downloads/history: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	body := readBody(t, resp)
	if !strings.Contains(body, "No history yet") {
		t.Errorf("empty history should show placeholder, got: %s", body)
	}
}
