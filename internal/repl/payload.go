package repl

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

func decodeBodyOptions[T any](options map[string]string, target *T) error {
	payload, err := readBodyOptions(options)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(payload, target); err != nil {
		return fmt.Errorf("failed to parse request body: %w", err)
	}
	return nil
}

func readBodyOptions(options map[string]string) ([]byte, error) {
	bodyValue, hasBody := popOption(options, "body")
	bodyFile, hasBodyFile := popOption(options, "body-file")

	switch {
	case hasBody && hasBodyFile:
		return nil, errors.New("use either --body or --body-file")
	case !hasBody && !hasBodyFile:
		return nil, errors.New("one of --body or --body-file is required")
	case hasBody:
		value := strings.TrimSpace(bodyValue)
		if value == "" {
			return nil, errors.New("--body can not be empty")
		}
		if strings.HasPrefix(value, "@") {
			return os.ReadFile(strings.TrimSpace(strings.TrimPrefix(value, "@")))
		}
		return []byte(value), nil
	default:
		path := strings.TrimSpace(bodyFile)
		if path == "" {
			return nil, errors.New("--body-file can not be empty")
		}
		return os.ReadFile(path)
	}
}
