package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_Defaults(t *testing.T) {
	clearEnv(t)

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	d := defaults()
	if cfg.MediaDir != d.MediaDir {
		t.Errorf("MediaDir: got %q, want %q", cfg.MediaDir, d.MediaDir)
	}
	if cfg.DBPath != d.DBPath {
		t.Errorf("DBPath: got %q, want %q", cfg.DBPath, d.DBPath)
	}
	if cfg.Port != d.Port {
		t.Errorf("Port: got %q, want %q", cfg.Port, d.Port)
	}
	if cfg.MaxConcurrentDL != d.MaxConcurrentDL {
		t.Errorf("MaxConcurrentDL: got %d, want %d", cfg.MaxConcurrentDL, d.MaxConcurrentDL)
	}
	if cfg.PollInterval != d.PollInterval {
		t.Errorf("PollInterval: got %v, want %v", cfg.PollInterval, d.PollInterval)
	}
}

func TestLoad_FileOverridesDefaults(t *testing.T) {
	clearEnv(t)

	f := writeTempConfig(t, `
media_dir: /mnt/media
db_path: /mnt/ytdlp.db
port: "9090"
max_concurrent_downloads: 4
poll_interval: 30m
`)

	cfg, err := Load(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.MediaDir != "/mnt/media" {
		t.Errorf("MediaDir: got %q", cfg.MediaDir)
	}
	if cfg.DBPath != "/mnt/ytdlp.db" {
		t.Errorf("DBPath: got %q", cfg.DBPath)
	}
	if cfg.Port != "9090" {
		t.Errorf("Port: got %q", cfg.Port)
	}
	if cfg.MaxConcurrentDL != 4 {
		t.Errorf("MaxConcurrentDL: got %d", cfg.MaxConcurrentDL)
	}
	if cfg.PollInterval != 30*time.Minute {
		t.Errorf("PollInterval: got %v", cfg.PollInterval)
	}
}

func TestLoad_EnvOverridesFile(t *testing.T) {
	clearEnv(t)

	f := writeTempConfig(t, `
media_dir: /from/file
port: "9090"
max_concurrent_downloads: 4
poll_interval: 30m
`)

	t.Setenv("YTD_MEDIA_DIR", "/from/env")
	t.Setenv("YTD_PORT", "7070")
	t.Setenv("YTD_MAX_CONCURRENT_DOWNLOADS", "8")
	t.Setenv("YTD_POLL_INTERVAL", "2h")

	cfg, err := Load(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.MediaDir != "/from/env" {
		t.Errorf("MediaDir: got %q, want env value", cfg.MediaDir)
	}
	if cfg.Port != "7070" {
		t.Errorf("Port: got %q, want env value", cfg.Port)
	}
	if cfg.MaxConcurrentDL != 8 {
		t.Errorf("MaxConcurrentDL: got %d, want env value", cfg.MaxConcurrentDL)
	}
	if cfg.PollInterval != 2*time.Hour {
		t.Errorf("PollInterval: got %v, want env value", cfg.PollInterval)
	}
}

func TestLoad_InvalidEnv(t *testing.T) {
	clearEnv(t)

	t.Setenv("YTD_MAX_CONCURRENT_DOWNLOADS", "notanumber")
	_, err := Load("")
	if err == nil {
		t.Error("expected error for invalid YTD_MAX_CONCURRENT_DOWNLOADS")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	clearEnv(t)

	f := writeTempConfig(t, `media_dir: [unclosed bracket`)
	_, err := Load(f)
	if err == nil {
		t.Error("expected error for malformed YAML")
	}
}

// clearEnv unsets all YTD_ env vars for the duration of the test.
func clearEnv(t *testing.T) {
	t.Helper()
	vars := []string{
		"YTD_MEDIA_DIR", "YTD_DB_PATH", "YTD_PORT",
		"YTD_MAX_CONCURRENT_DOWNLOADS", "YTD_POLL_INTERVAL", "YTD_CONFIG_PATH",
	}
	for _, v := range vars {
		t.Setenv(v, "")
		os.Unsetenv(v)
	}
}

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "config.yaml")
	if err != nil {
		t.Fatalf("creating temp config: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	f.Close()
	return f.Name()
}
