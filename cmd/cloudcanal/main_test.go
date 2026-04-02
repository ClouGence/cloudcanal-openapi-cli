package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
