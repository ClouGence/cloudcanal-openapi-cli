package config_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/ClouGence/cloudcanal-openapi-cli/internal/config"
)

func TestServiceSaveAndLoadState(t *testing.T) {
	dir := t.TempDir()
	service := config.NewService(filepath.Join(dir, "config.json"))

	state := config.State{
		Language:       "zh",
		CurrentProfile: "prod",
		Profiles: map[string]config.AppConfig{
			"dev": {
				APIBaseURL: "https://dev.example.com",
				AccessKey:  "dev-ak",
				SecretKey:  "dev-sk",
			},
			"prod": {
				APIBaseURL: "https://cc.example.com",
				AccessKey:  "access-key",
				SecretKey:  "secret-key",
			},
		},
	}
	if err := service.Save(state); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := service.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.Language != "zh" {
		t.Fatalf("Language = %q, want zh", loaded.Language)
	}
	if loaded.CurrentProfile != "prod" {
		t.Fatalf("CurrentProfile = %q, want prod", loaded.CurrentProfile)
	}
	if got := loaded.Profiles["prod"].APIBaseURL; got != "https://cc.example.com" {
		t.Fatalf("prod apiBaseUrl = %q, want https://cc.example.com", got)
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

func TestServiceRejectsMissingCurrentProfile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(`{"language":"en","profiles":{"dev":{"apiBaseUrl":"https://cc.example.com","accessKey":"ak","secretKey":"sk"}}}`), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	service := config.NewService(path)
	if _, err := service.Load(); err == nil {
		t.Fatal("Load() error = nil, want non-nil")
	}
}

func TestServiceLoadLanguageFromNewAndLegacyConfig(t *testing.T) {
	testCases := []struct {
		name    string
		content string
		want    string
	}{
		{name: "new schema", content: `{"language":"zh","currentProfile":"dev","profiles":{"dev":{"apiBaseUrl":"https://cc.example.com","accessKey":"ak","secretKey":"sk"}}}`, want: "zh"},
		{name: "legacy schema", content: `{"language":"zh","apiBaseUrl":"https://cc.example.com","accessKey":"ak","secretKey":"sk"}`, want: "zh"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "config.json")
			if err := os.WriteFile(path, []byte(tc.content), 0o600); err != nil {
				t.Fatalf("WriteFile() error = %v", err)
			}

			service := config.NewService(path)
			if got := service.LoadLanguage(); got != tc.want {
				t.Fatalf("LoadLanguage() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestServiceDetectsLegacyFormat(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{"language":"en","apiBaseUrl":"https://cc.example.com","accessKey":"ak","secretKey":"sk"}`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	service := config.NewService(path)
	_, err := service.Load()
	if !errors.Is(err, config.ErrLegacyFormat) {
		t.Fatalf("Load() error = %v, want legacy format error", err)
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
