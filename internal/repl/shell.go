package repl

import (
	"cloudcanal-openapi-cli/internal/app"
	"cloudcanal-openapi-cli/internal/console"
	"cloudcanal-openapi-cli/internal/i18n"
	"cloudcanal-openapi-cli/internal/util"
	"io"
	"strconv"
	"strings"
)

const prompt = "cloudcanal> "

type Shell struct {
	io           console.IO
	runtime      app.RuntimeContext
	outputFormat outputFormat
}

func NewShell(io console.IO, runtime app.RuntimeContext) *Shell {
	_ = i18n.SetLanguage(runtime.Config().NormalizedLanguage())
	shell := &Shell{io: io, runtime: runtime, outputFormat: outputText}
	if completable, ok := io.(console.TabCompletable); ok {
		completable.SetCompleter(shell.completeLine)
	}
	return shell
}

func (s *Shell) ExecuteArgs(args []string) error {
	if len(args) == 0 {
		return nil
	}
	return s.handleTokens(args)
}

func (s *Shell) Run() error {
	s.io.Println(i18n.T("common.typeHelp"))
	for {
		line, err := s.io.ReadLine(prompt)
		if err != nil {
			if err == io.EOF {
				s.io.Println("")
				return nil
			}
			if console.IsPromptAborted(err) {
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
			s.PrintError(err)
		}
	}
}

func (s *Shell) handle(commandLine string) error {
	tokens, err := splitCommandLine(commandLine)
	if err != nil {
		return err
	}
	return s.handleTokens(tokens)
}

func (s *Shell) handleTokens(tokens []string) error {
	filteredTokens, format, err := extractOutputFormat(tokens)
	if err != nil {
		return wrapCommandError(err, outputText)
	}
	tokens = filteredTokens
	if len(tokens) == 0 {
		return nil
	}
	previousFormat := s.outputFormat
	s.outputFormat = format
	defer func() {
		s.outputFormat = previousFormat
	}()

	if helpText, ok := RenderCommandHelp(tokens); ok {
		s.io.Println(helpText)
		return nil
	}

	switch strings.ToLower(tokens[0]) {
	case "clear", "cls":
		s.io.ClearScreen()
		return nil
	case "completion":
		return wrapCommandError(s.handleCompletion(tokens), format)
	case "__complete":
		s.printHiddenCompletions(tokens[1:])
		return nil
	case "jobs":
		return wrapCommandError(s.handleJobs(tokens), format)
	case "datasources":
		return wrapCommandError(s.handleDataSources(tokens), format)
	case "clusters":
		return wrapCommandError(s.handleClusters(tokens), format)
	case "workers":
		return wrapCommandError(s.handleWorkers(tokens), format)
	case "consolejobs":
		return wrapCommandError(s.handleConsoleJobs(tokens), format)
	case "job-config", "jobconfig":
		return wrapCommandError(s.handleJobConfig(tokens), format)
	case "schemas", "schema":
		return wrapCommandError(s.handleSchemas(tokens), format)
	case "config":
		return wrapCommandError(s.handleConfig(tokens), format)
	case "lang", "language":
		return wrapCommandError(s.handleLang(tokens), format)
	default:
		s.printUnknownCommand(tokens[0])
		return nil
	}
}

func (s *Shell) handleConfig(tokens []string) error {
	if len(tokens) < 2 {
		s.io.Println(s.usageConfig())
		return nil
	}
	switch strings.ToLower(tokens[1]) {
	case "show":
		if len(tokens) != 2 {
			s.io.Println(s.usageConfigShow())
			return nil
		}
		cfg := s.runtime.Config()
		if s.isJSONOutput() {
			return s.printJSON(map[string]any{
				"apiBaseUrl":                 cfg.APIBaseURL,
				"accessKeyMasked":            util.MaskSecret(cfg.AccessKey),
				"language":                   cfg.NormalizedLanguage(),
				"httpTimeoutSeconds":         cfg.HTTPTimeoutSecondsValue(),
				"httpReadMaxRetries":         cfg.HTTPReadMaxRetriesValue(),
				"httpReadRetryBackoffMillis": cfg.HTTPReadRetryBackoffMillisValue(),
			})
		}
		s.io.Println(i18n.T("config.apiBaseUrlLabel") + ": " + cfg.APIBaseURL)
		s.io.Println(i18n.T("config.accessKeyLabel") + ": " + util.MaskSecret(cfg.AccessKey))
		s.io.Println(i18n.T("config.languageLabel") + ": " + cfg.NormalizedLanguage())
		s.io.Println(i18n.T("config.httpTimeoutLabel") + ": " + strconv.Itoa(cfg.HTTPTimeoutSecondsValue()))
		s.io.Println(i18n.T("config.httpReadMaxRetriesLabel") + ": " + strconv.Itoa(cfg.HTTPReadMaxRetriesValue()))
		s.io.Println(i18n.T("config.httpReadRetryBackoffMillisLabel") + ": " + strconv.Itoa(cfg.HTTPReadRetryBackoffMillisValue()))
		return nil
	case "init":
		if len(tokens) != 2 {
			s.io.Println(s.usageConfigInit())
			return nil
		}
		updated, err := s.runtime.Reinitialize(s.io)
		if err != nil {
			return err
		}
		if updated {
			s.io.Println(i18n.T("runtime.configUpdated"))
		}
		return nil
	case "lang":
		return s.handleLanguageTokens(tokens[2:], "config lang", s.usageConfigLang())
	default:
		s.printUnknownSubcommand("config", tokens[1], configSubcommands, s.usageConfig())
		return nil
	}
}

func (s *Shell) handleLang(tokens []string) error {
	return s.handleLanguageTokens(tokens[1:], "config lang", s.usageConfigLang())
}

func (s *Shell) handleLanguageTokens(tokens []string, group string, usage string) error {
	if len(tokens) == 0 || strings.EqualFold(tokens[0], "show") {
		if len(tokens) > 1 {
			s.io.Println(usage)
			return nil
		}
		if s.isJSONOutput() {
			return s.printJSON(map[string]any{
				"language":  s.runtime.Config().NormalizedLanguage(),
				"supported": []string{"en", "zh"},
			})
		}
		s.io.Println(i18n.T("lang.current", s.runtime.Config().NormalizedLanguage()))
		s.io.Println(i18n.T("common.supportedLanguages"))
		return nil
	}
	if strings.EqualFold(tokens[0], "set") {
		if len(tokens) != 2 {
			s.io.Println(usage)
			return nil
		}
		if err := s.runtime.SetLanguage(tokens[1]); err != nil {
			return err
		}
		if s.isJSONOutput() {
			return s.printJSON(map[string]any{
				"language": s.runtime.Config().NormalizedLanguage(),
				"message":  i18n.T("lang.updated", i18n.DisplayName(s.runtime.Config().NormalizedLanguage())),
			})
		}
		s.io.Println(i18n.T("lang.updated", i18n.DisplayName(s.runtime.Config().NormalizedLanguage())))
		return nil
	}
	s.printUnknownSubcommand(group, tokens[0], langSubcommands, usage)
	return nil
}
