package config_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/ClouGence/cloudcanal-openapi-cli/internal/config"
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/console"
	"github.com/ClouGence/cloudcanal-openapi-cli/test/testsupport"
	"github.com/peterh/liner"
)

func TestWizardReturnsConfigAfterSuccessfulValidation(t *testing.T) {
	io := testsupport.NewTestConsole("https://cc.example.com", "test-ak", "test-sk")

	wizard := config.NewWizard(io, func(cfg config.AppConfig) error {
		return nil
	}, "dev", config.AppConfig{})

	cfg, err := wizard.Run()
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if cfg == nil {
		t.Fatal("Run() returned nil config")
	}
	if cfg.APIBaseURL != "https://cc.example.com" {
		t.Fatalf("APIBaseURL = %q, want https://cc.example.com", cfg.APIBaseURL)
	}
}

func TestWizardDoesNotPersistOnValidationFailureThenExit(t *testing.T) {
	io := testsupport.NewTestConsole("https://cc.example.com", "test-ak", "test-sk", "exit")

	wizard := config.NewWizard(io, func(cfg config.AppConfig) error {
		return errors.New("authentication failed")
	}, "prod", config.AppConfig{})

	cfg, err := wizard.Run()
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if cfg != nil {
		t.Fatalf("Run() config = %+v, want nil", *cfg)
	}
	if out := io.Output(); out == "" || !strings.Contains(out, "Configuration validation failed") || !strings.Contains(out, "apiHost [https://cc.example.com]: ") {
		t.Fatalf("wizard output missing validation failure reuse prompts: %q", out)
	}
}

func TestWizardReusesCurrentValuesAndDoesNotPrintSecret(t *testing.T) {
	io := testsupport.NewTestConsole("", "", "")

	wizard := config.NewWizard(io, func(cfg config.AppConfig) error {
		return nil
	}, "prod", config.AppConfig{
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
		"CloudCanal CLI profile initialization (prod)",
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

func TestWizardTreatsPromptAbortAsCancellation(t *testing.T) {
	io := &promptAbortConsole{}

	wizard := config.NewWizard(io, func(cfg config.AppConfig) error {
		return nil
	}, "dev", config.AppConfig{})

	cfg, err := wizard.Run()
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if cfg != nil {
		t.Fatalf("Run() config = %+v, want nil", *cfg)
	}
	if !strings.Contains(io.output.String(), "apiHost: ") {
		t.Fatalf("output = %q, want apiHost prompt", io.output.String())
	}
}

type promptAbortConsole struct {
	output strings.Builder
}

func (p *promptAbortConsole) ReadLine(prompt string) (string, error) {
	p.output.WriteString(prompt)
	return "", liner.ErrPromptAborted
}

func (p *promptAbortConsole) ReadSecret(prompt string) (string, error) {
	p.output.WriteString(prompt)
	return "", liner.ErrPromptAborted
}

func (p *promptAbortConsole) Println(text string) {
	p.output.WriteString(text)
	p.output.WriteString("\n")
}

func (p *promptAbortConsole) ClearScreen() {}

var _ console.IO = (*promptAbortConsole)(nil)
