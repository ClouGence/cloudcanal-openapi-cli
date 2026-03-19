package repl

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func splitCommandLine(line string) ([]string, error) {
	var (
		tokens       []string
		current      strings.Builder
		quote        rune
		escaped      bool
		tokenStarted bool
	)

	for _, r := range line {
		switch {
		case escaped:
			current.WriteRune(r)
			tokenStarted = true
			escaped = false
		case r == '\\':
			escaped = true
		case quote != 0:
			if r == quote {
				quote = 0
				continue
			}
			current.WriteRune(r)
			tokenStarted = true
		case r == '"' || r == '\'':
			quote = r
			tokenStarted = true
		case unicode.IsSpace(r):
			if tokenStarted {
				tokens = append(tokens, current.String())
				current.Reset()
				tokenStarted = false
			}
		default:
			current.WriteRune(r)
			tokenStarted = true
		}
	}

	if escaped {
		return nil, fmt.Errorf("unterminated escape sequence")
	}
	if quote != 0 {
		return nil, fmt.Errorf("unterminated quote")
	}
	if tokenStarted {
		tokens = append(tokens, current.String())
	}
	return tokens, nil
}

func parseFlagArgs(tokens []string) (map[string]string, error) {
	options := make(map[string]string, len(tokens))
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		if !strings.HasPrefix(token, "--") {
			return nil, fmt.Errorf("unexpected argument: %s", token)
		}

		raw := strings.TrimPrefix(token, "--")
		if raw == "" {
			return nil, fmt.Errorf("invalid option: %s", token)
		}

		name := raw
		value := "true"
		if strings.Contains(raw, "=") {
			var ok bool
			name, value, ok = strings.Cut(raw, "=")
			if !ok || name == "" {
				return nil, fmt.Errorf("invalid option: %s", token)
			}
		} else if i+1 < len(tokens) && !strings.HasPrefix(tokens[i+1], "--") {
			value = tokens[i+1]
			i++
		}

		if _, exists := options[name]; exists {
			return nil, fmt.Errorf("duplicate option: --%s", name)
		}
		options[name] = value
	}
	return options, nil
}

func popOption(options map[string]string, names ...string) (string, bool) {
	for _, name := range names {
		if value, ok := options[name]; ok {
			delete(options, name)
			return value, true
		}
	}
	return "", false
}

func ensureNoUnknownOptions(options map[string]string) error {
	for name := range options {
		return fmt.Errorf("unknown option: --%s", name)
	}
	return nil
}

func parsePositiveInt64(value string, fieldName string) (int64, error) {
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", fieldName)
	}
	return parsed, nil
}

func parsePositiveInt64Option(options map[string]string, fieldName string, names ...string) (int64, error) {
	value, ok := popOption(options, names...)
	if !ok {
		return 0, nil
	}
	return parsePositiveInt64(value, fieldName)
}

func parseBoolOption(options map[string]string, fieldName string, names ...string) (*bool, error) {
	value, ok := popOption(options, names...)
	if !ok {
		return nil, nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return nil, fmt.Errorf("%s must be a boolean", fieldName)
	}
	return &parsed, nil
}

func formatOptionalInt64(value int64) string {
	if value == 0 {
		return "-"
	}
	return strconv.FormatInt(value, 10)
}

func formatBool(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

func orDash(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}
	return value
}
