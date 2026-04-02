package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ClouGence/cloudcanal-openapi-cli/internal/app"
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/buildinfo"
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/config"
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/console"
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/i18n"
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/repl"
)

func main() {
	if handled, exitCode := handleEarlyCommands(os.Args[1:]); handled {
		os.Exit(exitCode)
	}

	io := console.NewStdIO(os.Stdin, os.Stdout)
	if closer, ok := any(io).(interface{ Close() error }); ok {
		defer func() { _ = closer.Close() }()
	}
	runtime := app.NewRuntime(config.NewService(""))
	ok, err := runtime.InitializeIfNeeded(io)
	if err != nil {
		io.Println(i18n.T("common.fatalErrorPrefix", err.Error()))
		os.Exit(1)
	}
	if !ok {
		return
	}

	shell := repl.NewShell(io, runtime)
	if len(os.Args) > 1 {
		if err := shell.ExecuteArgs(os.Args[1:]); err != nil {
			shell.PrintFatalError(err)
			os.Exit(1)
		}
		return
	}

	if err := shell.Run(); err != nil {
		shell.PrintFatalError(err)
		os.Exit(1)
	}
}

func handleEarlyCommands(args []string) (bool, int) {
	if os.Getenv(repl.CompletionEnvVar) == "1" {
		for _, candidate := range repl.CompletionCandidates(args, false) {
			fmt.Println(candidate)
		}
		return true, 0
	}

	if len(args) > 0 {
		_ = i18n.SetLanguage(config.NewService("").LoadLanguage())
	}
	if helpText, ok := repl.RenderCommandHelp(args); ok {
		fmt.Println(helpText)
		return true, 0
	}
	if handled, exitCode := handleVersionCommand(args); handled {
		return true, exitCode
	}

	if len(args) == 0 {
		return false, 0
	}

	switch strings.ToLower(args[0]) {
	case "completion":
		script, err := repl.RenderCompletionScript(args[1:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return true, 1
		}
		fmt.Print(script)
		return true, 0
	case "__complete":
		for _, candidate := range repl.CompletionCandidates(args[1:], false) {
			fmt.Println(candidate)
		}
		return true, 0
	default:
		return false, 0
	}
}

func handleVersionCommand(args []string) (bool, int) {
	filtered, format, err := extractOutputFormat(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return true, 1
	}
	if len(filtered) == 0 {
		return false, 0
	}

	switch {
	case len(filtered) == 1 && strings.EqualFold(filtered[0], "--version"):
		return true, printVersion(format)
	case strings.EqualFold(filtered[0], "version"):
		if len(filtered) != 1 {
			fmt.Fprintln(os.Stderr, versionUsageText())
			return true, 1
		}
		return true, printVersion(format)
	case containsVersionFlag(filtered):
		fmt.Fprintln(os.Stderr, versionFlagErrorText())
		return true, 1
	default:
		return false, 0
	}
}

func printVersion(format string) int {
	info := buildinfo.Current()
	if format == "json" {
		data, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 1
		}
		fmt.Println(string(data))
		return 0
	}

	fmt.Println("version: " + info.Version)
	fmt.Println("commit: " + info.Commit)
	fmt.Println("buildTime: " + info.BuildTime)
	return 0
}

func extractOutputFormat(args []string) ([]string, string, error) {
	format := "text"
	filtered := make([]string, 0, len(args))
	seen := false

	for i := 0; i < len(args); i++ {
		token := args[i]
		switch {
		case token == "--output":
			if i+1 >= len(args) {
				return nil, "", errors.New(i18n.T("parser.outputOptionRequiresValue"))
			}
			if seen {
				return nil, "", errors.New(i18n.T("parser.duplicateOption", "output"))
			}
			value := strings.ToLower(strings.TrimSpace(args[i+1]))
			if value != "text" && value != "json" {
				return nil, "", errors.New(i18n.T("parser.outputOptionInvalid"))
			}
			format = value
			seen = true
			i++
		case strings.HasPrefix(token, "--output="):
			if seen {
				return nil, "", errors.New(i18n.T("parser.duplicateOption", "output"))
			}
			_, value, _ := strings.Cut(token, "=")
			value = strings.ToLower(strings.TrimSpace(value))
			if value != "text" && value != "json" {
				return nil, "", errors.New(i18n.T("parser.outputOptionInvalid"))
			}
			format = value
			seen = true
		default:
			filtered = append(filtered, token)
		}
	}

	return filtered, format, nil
}

func containsVersionFlag(args []string) bool {
	for _, arg := range args {
		if strings.EqualFold(arg, "--version") {
			return true
		}
	}
	return false
}

func versionUsageText() string {
	if i18n.CurrentLanguage() == i18n.Chinese {
		return "用法：version"
	}
	return "Usage: version"
}

func versionFlagErrorText() string {
	if i18n.CurrentLanguage() == i18n.Chinese {
		return "--version 只能单独使用，或与 --output 一起使用"
	}
	return "--version can only be used by itself or with --output"
}
