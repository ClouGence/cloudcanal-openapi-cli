package repl_test

import (
	"cloudcanal-openapi-cli/internal/config"
	"cloudcanal-openapi-cli/internal/repl"
	"cloudcanal-openapi-cli/test/testsupport"
	"strings"
	"testing"
)

func TestShellRegistersCompleter(t *testing.T) {
	runtime := newCompletionRuntime()
	io := testsupport.NewTestConsole()

	repl.NewShell(io, runtime)

	candidates := io.Complete("jo")
	if !contains(candidates, "jobs") {
		t.Fatalf("completion candidates = %v, want jobs", candidates)
	}
	if !contains(candidates, "job-config") {
		t.Fatalf("completion candidates = %v, want job-config", candidates)
	}
	schemaCandidates := io.Complete("sc")
	if !contains(schemaCandidates, "schemas") {
		t.Fatalf("completion candidates = %v, want schemas", schemaCandidates)
	}
}

func TestCompletionCandidatesSuggestCommandsFlagsAndValues(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		want []string
	}{
		{name: "top level", args: []string{""}, want: []string{"help", "jobs", "config", "schemas"}},
		{name: "top level help flag", args: []string{"--h"}, want: []string{"--help"}},
		{name: "top level global flag", args: []string{"--o"}, want: []string{"--output"}},
		{name: "jobs help flag", args: []string{"jobs", "--h"}, want: []string{"--help"}},
		{name: "config subcommand", args: []string{"config", "la"}, want: []string{"lang"}},
		{name: "config lang value", args: []string{"config", "lang", "set", ""}, want: []string{"en", "zh"}},
		{name: "jobs subcommand", args: []string{"jobs", "re"}, want: []string{"replay"}},
		{name: "jobs create flag", args: []string{"jobs", "create", "--bo"}, want: []string{"--body", "--body-file"}},
		{name: "list flag", args: []string{"jobs", "list", "--so"}, want: []string{"--source-id"}},
		{name: "global flag value", args: []string{"jobs", "list", "--output", ""}, want: []string{"text", "json"}},
		{name: "bool value", args: []string{"job-config", "specs", "--initial-sync", ""}, want: []string{"true", "false"}},
		{name: "inline bool value", args: []string{"job-config", "specs", "--initial-sync=t"}, want: []string{"--initial-sync=true"}},
		{name: "datasource add file flag", args: []string{"datasources", "add", "--sec"}, want: []string{"--security-file"}},
		{name: "worker alert bool value", args: []string{"workers", "update-alert", "5", "--phone", ""}, want: []string{"true", "false"}},
		{name: "transform job type flag", args: []string{"job-config", "transform-job-type", "--sou"}, want: []string{"--source-type"}},
		{name: "schemas flag", args: []string{"schemas", "list-trans-objs-by-meta", "--src"}, want: []string{"--src-db", "--src-schema", "--src-trans-obj"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			candidates := repl.CompletionCandidates(tc.args, false)
			for _, want := range tc.want {
				if !contains(candidates, want) {
					t.Fatalf("completion candidates for %v = %v, want %q", tc.args, candidates, want)
				}
			}
		})
	}
}

func TestCompletionCandidatesHideInternalAndAliasCommands(t *testing.T) {
	candidates := repl.CompletionCandidates([]string{""}, true)
	for _, hidden := range []string{"clear", "cls", "completion", "jobconfig", "schema", "lang", "language", "quit"} {
		if contains(candidates, hidden) {
			t.Fatalf("completion candidates = %v, should not contain %q", candidates, hidden)
		}
	}

	helpCandidates := repl.CompletionCandidates([]string{"help", ""}, false)
	for _, hidden := range []string{"completion", "lang"} {
		if contains(helpCandidates, hidden) {
			t.Fatalf("help completion candidates = %v, should not contain %q", helpCandidates, hidden)
		}
	}
}

func TestCompletionCandidatesDoNotSuggestFlagsWhileTypingFreeformValues(t *testing.T) {
	candidates := repl.CompletionCandidates([]string{"jobs", "list", "--name", ""}, false)
	if len(candidates) != 0 {
		t.Fatalf("completion candidates = %v, want empty while typing --name value", candidates)
	}
}

func TestShellCompletionCommands(t *testing.T) {
	runtime := newCompletionRuntime()
	io := testsupport.NewTestConsole()
	shell := repl.NewShell(io, runtime)

	if err := shell.ExecuteArgs([]string{"completion", "zsh"}); err != nil {
		t.Fatalf("ExecuteArgs(completion zsh) error = %v", err)
	}
	if !strings.Contains(io.Output(), "#compdef cloudcanal") {
		t.Fatalf("zsh completion output missing compdef: %q", io.Output())
	}
	if !strings.Contains(io.Output(), repl.CompletionEnvVar+"=1") {
		t.Fatalf("zsh completion output missing hidden completion env var: %q", io.Output())
	}
	if strings.Contains(io.Output(), "__complete") {
		t.Fatalf("zsh completion output should not expose __complete: %q", io.Output())
	}

	io = testsupport.NewTestConsole()
	shell = repl.NewShell(io, runtime)
	if err := shell.ExecuteArgs([]string{"__complete", "jobs", "l"}); err != nil {
		t.Fatalf("ExecuteArgs(__complete) error = %v", err)
	}
	if !strings.Contains(io.Output(), "list") {
		t.Fatalf("hidden completion output missing list: %q", io.Output())
	}
}

func newCompletionRuntime() *fakeRuntime {
	return &fakeRuntime{
		cfg:         config.AppConfig{APIBaseURL: "https://cc.example.com", AccessKey: "abcdefghijkl", SecretKey: "qrstuvwxyz1234"},
		dataJobs:    &fakeDataJobs{},
		dataSources: &fakeDataSources{},
		clusters:    &fakeClusters{},
		workers:     &fakeWorkers{},
		consoleJobs: &fakeConsoleJobs{},
		jobConfigs:  &fakeJobConfigs{},
	}
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
