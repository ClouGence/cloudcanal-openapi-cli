package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"cloudcanal-openapi-cli/internal/config"
)

func TestServiceSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	service := config.NewService(filepath.Join(dir, "config.json"))

	cfg := config.AppConfig{
		APIBaseURL: "https://cc.example.com",
		AccessKey:  "access-key",
		SecretKey:  "secret-key",
	}
	if err := service.Save(cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := service.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.APIBaseURL != cfg.APIBaseURL || loaded.AccessKey != cfg.AccessKey || loaded.SecretKey != cfg.SecretKey || loaded.Language != "en" {
		t.Fatalf("loaded config = %+v, want %+v", loaded, cfg)
	}
}

func TestServiceRejectsInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte("{invalid"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	service := config.NewService(path)
	if _, err := service.Load(); err == nil {
		t.Fatal("Load() error = nil, want non-nil")
	}
}

func TestServiceRejectsMissingField(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(`{"apiBaseUrl":"https://cc.example.com","accessKey":"ak"}`), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	service := config.NewService(path)
	if _, err := service.Load(); err == nil {
		t.Fatal("Load() error = nil, want non-nil")
	}
}

func TestConfigNetworkSettingsDefaultsAndValidation(t *testing.T) {
	cfg := config.AppConfig{
		APIBaseURL: "https://cc.example.com",
		AccessKey:  "access-key",
		SecretKey:  "secret-key",
	}
	if got := cfg.HTTPTimeoutSecondsValue(); got != 10 {
		t.Fatalf("HTTPTimeoutSecondsValue() = %d, want 10", got)
	}
	if got := cfg.HTTPReadMaxRetriesValue(); got != 0 {
		t.Fatalf("HTTPReadMaxRetriesValue() = %d, want 0", got)
	}
	if got := cfg.HTTPReadRetryBackoffMillisValue(); got != 250 {
		t.Fatalf("HTTPReadRetryBackoffMillisValue() = %d, want 250", got)
	}

	cfg.HTTPReadMaxRetries = -1
	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want non-nil for negative retry count")
	}
}

func TestDefaultPathUsesCloudCanalCLIDirectory(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	got := config.DefaultPath()
	want := filepath.Join(home, ".cloudcanal-cli", "config.json")
	if got != want {
		t.Fatalf("DefaultPath() = %q, want %q", got, want)
	}
}

func TestServiceLoadMigratesLegacyDefaultConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	legacyPath := config.LegacyDefaultPath()
	if err := os.MkdirAll(filepath.Dir(legacyPath), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	content := []byte(`{"apiBaseUrl":"https://cc.example.com","accessKey":"legacy-ak","secretKey":"legacy-sk","language":"zh"}`)
	if err := os.WriteFile(legacyPath, content, 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	service := config.NewService("")
	if !service.Exists() {
		t.Fatal("Exists() = false, want true when only legacy config exists")
	}

	loaded, err := service.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.APIBaseURL != "https://cc.example.com" || loaded.AccessKey != "legacy-ak" || loaded.SecretKey != "legacy-sk" || loaded.Language != "zh" {
		t.Fatalf("loaded config = %+v, want legacy values", loaded)
	}

	newPath := config.DefaultPath()
	if _, err := os.Stat(newPath); err != nil {
		t.Fatalf("new config path stat error = %v", err)
	}
	if _, err := os.Stat(legacyPath); !os.IsNotExist(err) {
		t.Fatalf("legacy config still exists, stat err = %v", err)
	}
}

func TestServiceLoadWithCustomPathDoesNotReadLegacyConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	legacyPath := config.LegacyDefaultPath()
	if err := os.MkdirAll(filepath.Dir(legacyPath), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(legacyPath, []byte(`{"apiBaseUrl":"https://cc.example.com","accessKey":"legacy-ak","secretKey":"legacy-sk"}`), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	service := config.NewService(filepath.Join(t.TempDir(), "config.json"))
	if service.Exists() {
		t.Fatal("Exists() = true, want false for custom path without file")
	}
	if _, err := service.Load(); !os.IsNotExist(err) {
		t.Fatalf("Load() error = %v, want os.ErrNotExist", err)
	}
}
