package repl

import (
	"cloudcanal-openapi-cli/internal/i18n"
	"cloudcanal-openapi-cli/internal/util"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type outputFormat string

const (
	outputText outputFormat = "text"
	outputJSON outputFormat = "json"
)

type commandError struct {
	format outputFormat
	err    error
}

func (e *commandError) Error() string {
	return e.err.Error()
}

func (e *commandError) Unwrap() error {
	return e.err
}

func wrapCommandError(err error, format outputFormat) error {
	if err == nil {
		return nil
	}
	return &commandError{format: format, err: err}
}

func outputFormatFromError(err error) outputFormat {
	var commandErr *commandError
	if errors.As(err, &commandErr) {
		return commandErr.format
	}
	return outputText
}

func extractOutputFormat(tokens []string) ([]string, outputFormat, error) {
	format := outputText
	filtered := make([]string, 0, len(tokens))
	seen := false

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		switch {
		case token == "--output":
			if i+1 >= len(tokens) {
				return nil, outputText, errors.New(i18n.T("parser.outputOptionRequiresValue"))
			}
			if seen {
				return nil, outputText, errors.New(i18n.T("parser.duplicateOption", "output"))
			}
			parsed, err := parseOutputFormat(tokens[i+1])
			if err != nil {
				return nil, outputText, err
			}
			format = parsed
			seen = true
			i++
		case strings.HasPrefix(token, "--output="):
			if seen {
				return nil, outputText, errors.New(i18n.T("parser.duplicateOption", "output"))
			}
			_, value, _ := strings.Cut(token, "=")
			parsed, err := parseOutputFormat(value)
			if err != nil {
				return nil, outputText, err
			}
			format = parsed
			seen = true
		default:
			filtered = append(filtered, token)
		}
	}

	return filtered, format, nil
}

func parseOutputFormat(value string) (outputFormat, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "text":
		return outputText, nil
	case "json":
		return outputJSON, nil
	default:
		return outputText, errors.New(i18n.T("parser.outputOptionInvalid"))
	}
}

func (s *Shell) isJSONOutput() bool {
	return s.outputFormat == outputJSON
}

func (s *Shell) printJSON(value any) error {
	if err := s.writeJSON(value); err != nil {
		return wrapCommandError(err, s.outputFormat)
	}
	return nil
}

func (s *Shell) writeJSON(value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode JSON output: %w", err)
	}
	s.io.Println(string(data))
	return nil
}

func (s *Shell) PrintError(err error) {
	s.printError(err, false)
}

func (s *Shell) PrintFatalError(err error) {
	s.printError(err, true)
}

func (s *Shell) printError(err error, fatal bool) {
	message := util.SummarizeError(err)
	if outputFormatFromError(err) == outputJSON {
		payload := map[string]any{"error": message}
		if fatal {
			payload["fatal"] = true
		}
		if jsonErr := s.writeJSON(payload); jsonErr == nil {
			return
		}
	}

	if fatal {
		s.io.Println(i18n.T("common.fatalErrorPrefix", message))
		return
	}
	s.io.Println(i18n.T("common.errorPrefix", message))
}

func (s *Shell) printActionResult(kind string, resource string, action string, id int64) error {
	message := s.actionMessage(kind, id)
	if s.isJSONOutput() {
		return s.printJSON(map[string]any{
			"resource": resource,
			"action":   action,
			"id":       id,
			"message":  message,
		})
	}
	s.io.Println(message)
	return nil
}
