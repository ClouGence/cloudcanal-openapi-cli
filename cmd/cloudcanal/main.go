package main

import (
	"fmt"
	"os"
	"strings"

	"cloudcanal-openapi-cli/internal/app"
	"cloudcanal-openapi-cli/internal/config"
	"cloudcanal-openapi-cli/internal/console"
	"cloudcanal-openapi-cli/internal/i18n"
	"cloudcanal-openapi-cli/internal/repl"
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
