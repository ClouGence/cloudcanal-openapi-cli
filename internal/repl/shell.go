package repl

import (
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/app"
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/buildinfo"
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/config"
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/console"
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/i18n"
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/util"
	"io"
	"strconv"
	"strings"
)

type Shell struct {
	io           console.IO
	runtime      app.RuntimeContext
	outputFormat outputFormat
}

func NewShell(io console.IO, runtime app.RuntimeContext) *Shell {
	_ = i18n.SetLanguage(runtime.Language())
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
		line, err := s.io.ReadLine(s.prompt())
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

	if spec := findRootCommand(tokens[0]); spec != nil && spec.run != nil {
		return wrapCommandError(spec.run(s, tokens), format)
	}

	s.printUnknownCommand(tokens[0])
	return nil
}

func (s *Shell) handleConfig(tokens []string) error {
	return s.dispatchRegisteredCommand(tokens)
}

func (s *Shell) handleLang(tokens []string) error {
	return s.dispatchRegisteredCommand(tokens)
}

func (s *Shell) handleVersion(tokens []string) error {
	return s.dispatchRegisteredCommand(tokens)
}

func (s *Shell) runConfigShow(tokens []string) error {
	if len(tokens) != 2 {
		s.io.Println(s.usageConfigShow())
		return nil
	}

	cfg := s.runtime.Config()
	if s.isJSONOutput() {
		return s.printJSON(map[string]any{
			"currentProfile":             s.runtime.CurrentProfile(),
			"apiBaseUrl":                 cfg.APIBaseURL,
			"accessKeyMasked":            util.MaskSecret(cfg.AccessKey),
			"language":                   s.runtime.Language(),
			"httpTimeoutSeconds":         cfg.HTTPTimeoutSecondsValue(),
			"httpReadMaxRetries":         cfg.HTTPReadMaxRetriesValue(),
			"httpReadRetryBackoffMillis": cfg.HTTPReadRetryBackoffMillisValue(),
		})
	}
	s.io.Println(i18n.T("config.currentProfileLabel") + ": " + s.runtime.CurrentProfile())
	s.io.Println(i18n.T("config.apiBaseUrlLabel") + ": " + cfg.APIBaseURL)
	s.io.Println(i18n.T("config.accessKeyLabel") + ": " + util.MaskSecret(cfg.AccessKey))
	s.io.Println(i18n.T("config.languageLabel") + ": " + s.runtime.Language())
	s.io.Println(i18n.T("config.httpTimeoutLabel") + ": " + strconv.Itoa(cfg.HTTPTimeoutSecondsValue()))
	s.io.Println(i18n.T("config.httpReadMaxRetriesLabel") + ": " + strconv.Itoa(cfg.HTTPReadMaxRetriesValue()))
	s.io.Println(i18n.T("config.httpReadRetryBackoffMillisLabel") + ": " + strconv.Itoa(cfg.HTTPReadRetryBackoffMillisValue()))
	return nil
}

func (s *Shell) runConfigInit(tokens []string) error {
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
}

func (s *Shell) runLanguageShow(tokens []string) error {
	expectedLen := languageValueIndex(tokens)
	if len(tokens) != expectedLen {
		s.io.Println(s.usageConfigLang())
		return nil
	}
	if s.isJSONOutput() {
		return s.printJSON(map[string]any{
			"language":  s.runtime.Language(),
			"supported": []string{"en", "zh"},
		})
	}
	s.io.Println(i18n.T("lang.current", s.runtime.Language()))
	s.io.Println(i18n.T("common.supportedLanguages"))
	return nil
}

func (s *Shell) runLanguageSet(tokens []string) error {
	valueIndex := languageValueIndex(tokens)
	if len(tokens) != valueIndex+1 {
		s.io.Println(s.usageConfigLang())
		return nil
	}
	if err := s.runtime.SetLanguage(tokens[valueIndex]); err != nil {
		return err
	}
	if s.isJSONOutput() {
		return s.printJSON(map[string]any{
			"language": s.runtime.Language(),
			"message":  i18n.T("lang.updated", i18n.DisplayName(s.runtime.Language())),
		})
	}
	s.io.Println(i18n.T("lang.updated", i18n.DisplayName(s.runtime.Language())))
	return nil
}

func (s *Shell) runProfilesList(tokens []string) error {
	if len(tokens) != 3 {
		s.io.Println(s.usageConfigProfiles())
		return nil
	}

	summaries := s.runtime.ProfileSummaries()
	if s.isJSONOutput() {
		return s.printJSON(map[string]any{
			"currentProfile": s.runtime.CurrentProfile(),
			"profiles":       summaries,
		})
	}

	rows := make([][]string, 0, len(summaries))
	for _, summary := range summaries {
		current := ""
		if summary.Current {
			current = "*"
		}
		rows = append(rows, []string{current, summary.Name, summary.APIBaseURL})
	}
	s.io.Println(util.FormatTable(s.profileListHeaders(), rows))
	return nil
}

func (s *Shell) runProfilesUse(tokens []string) error {
	if len(tokens) != 4 {
		s.io.Println(s.usageConfigProfiles())
		return nil
	}

	name := config.NormalizeProfileName(tokens[3])
	if err := s.runtime.UseProfile(name); err != nil {
		return err
	}
	if s.isJSONOutput() {
		return s.printJSON(map[string]any{
			"profile": name,
			"message": s.profileUsedMessage(name),
		})
	}
	s.io.Println(s.profileUsedMessage(name))
	return nil
}

func (s *Shell) runProfilesAdd(tokens []string) error {
	if len(tokens) != 4 {
		s.io.Println(s.usageConfigProfiles())
		return nil
	}

	name := config.NormalizeProfileName(tokens[3])
	added, err := s.runtime.AddProfile(name, s.io)
	if err != nil {
		return err
	}
	if !added {
		return nil
	}
	if s.isJSONOutput() {
		return s.printJSON(map[string]any{
			"profile": name,
			"message": s.profileAddedMessage(name),
		})
	}
	s.io.Println(s.profileAddedMessage(name))
	return nil
}

func (s *Shell) runProfilesRemove(tokens []string) error {
	if len(tokens) != 4 {
		s.io.Println(s.usageConfigProfiles())
		return nil
	}

	name := config.NormalizeProfileName(tokens[3])
	if err := s.runtime.RemoveProfile(name); err != nil {
		return err
	}
	if s.isJSONOutput() {
		return s.printJSON(map[string]any{
			"profile": name,
			"message": s.profileRemovedMessage(name),
		})
	}
	s.io.Println(s.profileRemovedMessage(name))
	return nil
}

func (s *Shell) runVersion(tokens []string) error {
	if len(tokens) != 1 {
		s.io.Println(s.usageVersion())
		return nil
	}

	info := buildinfo.Current()
	if s.isJSONOutput() {
		return s.printJSON(info)
	}
	s.io.Println("version: " + info.Version)
	s.io.Println("commit: " + info.Commit)
	s.io.Println("buildTime: " + info.BuildTime)
	return nil
}

func languageValueIndex(tokens []string) int {
	if len(tokens) > 0 && strings.EqualFold(tokens[0], "config") {
		return 3
	}
	return 2
}

func (s *Shell) prompt() string {
	currentProfile := s.runtime.CurrentProfile()
	if strings.TrimSpace(currentProfile) == "" {
		return "cloudcanal> "
	}
	return "cloudcanal[" + currentProfile + "]> "
}
