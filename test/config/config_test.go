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
	if loaded.APIBaseURL != cfg.APIBaseURL || loaded.AccessKey != cfg.AccessKey || loaded.SecretKey != cfg.SecretKey {
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
