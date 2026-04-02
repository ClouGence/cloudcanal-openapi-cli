package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ClouGence/cloudcanal-openapi-cli/internal/config"
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/i18n"
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/updatecheck"
)

func TestHandleEarlyCommandsSupportsVersionWithLegacyConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	configPath := filepath.Join(home, ".cloudcanal-cli", "config.json")
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	legacyConfig := `{"language":"zh","apiBaseUrl":"https://cc.example.com","accessKey":"ak","secretKey":"sk"}`
	if err := os.WriteFile(configPath, []byte(legacyConfig), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	stdout, stderr := captureProcessOutput(t, func() (bool, int) {
		return handleEarlyCommands([]string{"version"})
	})
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
	for _, want := range []string{"version: ", "commit: ", "buildTime: "} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("stdout missing %q in %q", want, stdout)
		}
	}
}

func TestHandleEarlyCommandsSupportsVersionFlagJSONWithInvalidConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	configPath := filepath.Join(home, ".cloudcanal-cli", "config.json")
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(configPath, []byte("{invalid"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	stdout, stderr := captureProcessOutput(t, func() (bool, int) {
		return handleEarlyCommands([]string{"--version", "--output", "json"})
	})
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("json.Unmarshal(stdout) error = %v, stdout = %q", err, stdout)
	}
	for _, key := range []string{"version", "commit", "buildTime"} {
		if _, ok := payload[key]; !ok {
			t.Fatalf("payload missing %q: %#v", key, payload)
		}
	}
}

func TestCommandContextLine(t *testing.T) {
	runtime := fakeCommandContextRuntime{
		currentProfile: "test",
		cfg:            config.AppConfig{APIBaseURL: "https://test.example.com"},
	}

	line, ok := commandContextLine([]string{"jobs", "list"}, runtime)
	if !ok {
		t.Fatal("commandContextLine() ok = false, want true")
	}
	if line != "Current profile: test (https://test.example.com)" {
		t.Fatalf("commandContextLine() = %q", line)
	}

	if _, ok := commandContextLine([]string{"jobs", "list", "--output", "json"}, runtime); ok {
		t.Fatal("commandContextLine(json) ok = true, want false")
	}
	if _, ok := commandContextLine([]string{"config", "show"}, runtime); ok {
		t.Fatal("commandContextLine(config show) ok = true, want false")
	}
}

func TestStartupUpdateLines(t *testing.T) {
	originalLanguage := i18n.CurrentLanguage()
	t.Cleanup(func() {
		_ = i18n.SetLanguage(originalLanguage)
	})
	_ = i18n.SetLanguage(i18n.English)

	checker := fakeUpdateNoticeChecker{
		notice: updatecheck.Notice{
			CurrentVersion: "v0.1.2",
			LatestVersion:  "v0.1.3",
			UpgradeCommand: "curl -fsSL https://example.com/install.sh | bash",
		},
	}

	lines := startupUpdateLines([]string{"jobs", "list"}, false, checker)
	want := []string{
		"New version available: v0.1.3 (current: v0.1.2)",
		"Upgrade command: curl -fsSL https://example.com/install.sh | bash",
	}
	if !reflect.DeepEqual(lines, want) {
		t.Fatalf("startupUpdateLines() = %#v, want %#v", lines, want)
	}

	if lines := startupUpdateLines([]string{"jobs", "list", "--output", "json"}, false, checker); lines != nil {
		t.Fatalf("startupUpdateLines(json) = %#v, want nil", lines)
	}
	if lines := startupUpdateLines([]string{"version"}, false, checker); lines != nil {
		t.Fatalf("startupUpdateLines(version) = %#v, want nil", lines)
	}
	if lines := startupUpdateLines(nil, true, checker); !reflect.DeepEqual(lines, want) {
		t.Fatalf("startupUpdateLines(interactive) = %#v, want %#v", lines, want)
	}
}

func captureProcessOutput(t *testing.T, fn func() (bool, int)) (string, string) {
	t.Helper()

	originalStdout := os.Stdout
	originalStderr := os.Stderr

	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe(stdout) error = %v", err)
	}
	stderrReader, stderrWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe(stderr) error = %v", err)
	}

	os.Stdout = stdoutWriter
	os.Stderr = stderrWriter

	handled, exitCode := fn()

	_ = stdoutWriter.Close()
	_ = stderrWriter.Close()
	os.Stdout = originalStdout
	os.Stderr = originalStderr

	stdoutBytes, err := io.ReadAll(stdoutReader)
	if err != nil {
		t.Fatalf("ReadAll(stdout) error = %v", err)
	}
	stderrBytes, err := io.ReadAll(stderrReader)
	if err != nil {
		t.Fatalf("ReadAll(stderr) error = %v", err)
	}
	_ = stdoutReader.Close()
	_ = stderrReader.Close()

	if !handled || exitCode != 0 {
		t.Fatalf("handled=%v exitCode=%d, want true/0", handled, exitCode)
	}

	return strings.TrimSpace(string(stdoutBytes)), strings.TrimSpace(string(stderrBytes))
}

type fakeCommandContextRuntime struct {
	currentProfile string
	cfg            config.AppConfig
}

func (f fakeCommandContextRuntime) CurrentProfile() string {
	return f.currentProfile
}

func (f fakeCommandContextRuntime) Config() config.AppConfig {
	return f.cfg
}

type fakeUpdateNoticeChecker struct {
	notice updatecheck.Notice
	err    error
}

func (f fakeUpdateNoticeChecker) Check(string) (updatecheck.Notice, error) {
	return f.notice, f.err
}
