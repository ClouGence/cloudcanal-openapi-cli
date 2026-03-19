package config_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cloudcanal-openapi-cli/internal/config"
	"cloudcanal-openapi-cli/test/testsupport"
)

func TestWizardSavesConfigAfterSuccessfulValidation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	service := config.NewService(path)
	io := testsupport.NewTestConsole("https://cc.example.com", "test-ak", "test-sk")

	wizard := config.NewWizard(io, service, func(cfg config.AppConfig) error {
		return nil
	}, config.AppConfig{})

	cfg, err := wizard.Run()
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if cfg == nil {
		t.Fatal("Run() returned nil config")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("saved config missing: %v", err)
	}
}

func TestWizardDoesNotPersistOnValidationFailureThenExit(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	service := config.NewService(path)
	io := testsupport.NewTestConsole("https://cc.example.com", "test-ak", "test-sk", "exit")

	wizard := config.NewWizard(io, service, func(cfg config.AppConfig) error {
		return errors.New("authentication failed")
	}, config.AppConfig{})

	cfg, err := wizard.Run()
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if cfg != nil {
		t.Fatalf("Run() config = %+v, want nil", *cfg)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("config file exists unexpectedly, err = %v", err)
	}
	if out := io.Output(); out == "" || !strings.Contains(out, "Configuration validation failed") {
		t.Fatalf("wizard output missing validation failure: %q", out)
	}
}

func TestWizardReusesCurrentValuesAndDoesNotPrintSecret(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	service := config.NewService(path)
	io := testsupport.NewTestConsole("", "", "")

	wizard := config.NewWizard(io, service, func(cfg config.AppConfig) error {
		return nil
	}, config.AppConfig{
		APIBaseURL: "https://cc.example.com",
		AccessKey:  "current-ak",
		SecretKey:  "current-sk",
	})

	cfg, err := wizard.Run()
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if cfg == nil {
		t.Fatal("Run() returned nil config")
	}
	if cfg.SecretKey != "current-sk" {
		t.Fatalf("SecretKey = %q, want current-sk", cfg.SecretKey)
	}

	out := io.Output()
	for _, want := range []string{
		"Press Enter to keep the current value.",
		"apiHost [https://cc.example.com]: ",
		"ak [current-ak]: ",
		"sk [hidden]: ",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("wizard output missing %q in %q", want, out)
		}
	}
	if strings.Contains(out, "current-sk") {
		t.Fatalf("wizard output leaked secret key: %q", out)
	}
}
