package repl

import (
	"cloudcanal-openapi-cli/internal/app"
	"cloudcanal-openapi-cli/internal/console"
	"cloudcanal-openapi-cli/internal/util"
	"io"
	"strings"
)

const prompt = "cloudcanal> "

type Shell struct {
	io      console.IO
	runtime app.RuntimeContext
}

func NewShell(io console.IO, runtime app.RuntimeContext) *Shell {
	return &Shell{io: io, runtime: runtime}
}

func (s *Shell) ExecuteArgs(args []string) error {
	if len(args) == 0 {
		return nil
	}
	commandLine := strings.Join(args, " ")
	return s.handleTokens(args, commandLine)
}

func (s *Shell) Run() error {
	s.io.Println("Type 'help' to see available commands.")
	for {
		line, err := s.io.ReadLine(prompt)
		if err != nil {
			if err == io.EOF {
				s.io.Println("")
				return nil
			}
			return err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.EqualFold(line, "exit") || strings.EqualFold(line, "quit") {
			return nil
		}

		if err := s.handle(line); err != nil {
			s.io.Println("Error: " + util.SummarizeError(err))
		}
	}
}

func (s *Shell) handle(commandLine string) error {
	tokens, err := splitCommandLine(commandLine)
	if err != nil {
		return err
	}
	return s.handleTokens(tokens, commandLine)
}

func (s *Shell) handleTokens(tokens []string, commandLine string) error {
	if len(tokens) == 0 {
		return nil
	}

	switch strings.ToLower(tokens[0]) {
	case "help":
		s.printHelp()
		return nil
	case "jobs":
		return s.handleJobs(tokens)
	case "datasources":
		return s.handleDataSources(tokens)
	case "clusters":
		return s.handleClusters(tokens)
	case "workers":
		return s.handleWorkers(tokens)
	case "consolejobs":
		return s.handleConsoleJobs(tokens)
	case "job-config", "jobconfig":
		return s.handleJobConfig(tokens)
	case "config":
		return s.handleConfig(tokens)
	default:
		s.io.Println("Unknown command: " + commandLine)
		s.io.Println("Use 'help' to see available commands.")
		return nil
	}
}

func (s *Shell) handleConfig(tokens []string) error {
	if len(tokens) != 2 {
		s.io.Println("Usage: config show | config init")
		return nil
	}
	switch strings.ToLower(tokens[1]) {
	case "show":
		cfg := s.runtime.Config()
		s.io.Println("apiBaseUrl: " + cfg.APIBaseURL)
		s.io.Println("accessKey: " + util.MaskSecret(cfg.AccessKey))
		return nil
	case "init":
		updated, err := s.runtime.Reinitialize(s.io)
		if err != nil {
			return err
		}
		if updated {
			s.io.Println("Configuration updated.")
		}
		return nil
	default:
		s.io.Println("Usage: config show | config init")
		return nil
	}
}

func (s *Shell) printHelp() {
	s.io.Println("Available commands:")
	s.io.Println("  jobs list [--name NAME] [--type TYPE] [--desc DESC] [--source-id ID] [--target-id ID]")
	s.io.Println("  jobs show <jobId>")
	s.io.Println("  jobs schema <jobId>")
	s.io.Println("  jobs start <jobId>")
	s.io.Println("  jobs stop <jobId>")
	s.io.Println("  jobs delete <jobId>")
	s.io.Println("  jobs replay <jobId> [--auto-start] [--reset-to-created]")
	s.io.Println("  datasources list [--id ID] [--type TYPE] [--deploy-type TYPE] [--host-type TYPE] [--lifecycle STATE]")
	s.io.Println("  datasources show <dataSourceId>")
	s.io.Println("  clusters list [--name NAME] [--desc DESC] [--cloud CLOUD] [--region REGION]")
	s.io.Println("  workers list [--cluster-id ID] [--source-id ID] [--target-id ID]")
	s.io.Println("  workers start <workerId>")
	s.io.Println("  workers stop <workerId>")
	s.io.Println("  consolejobs show <consoleJobId>")
	s.io.Println("  job-config specs [--type TYPE] [--initial-sync=true|false] [--short-term-sync=true|false]")
	s.io.Println("  config show")
	s.io.Println("  config init")
	s.io.Println("  help")
	s.io.Println("  exit | quit")
}
