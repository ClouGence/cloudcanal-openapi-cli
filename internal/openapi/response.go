package openapi

import (
	"encoding/json"
	"errors"
	"strings"
)

type Response struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func EnsureSuccess(resp Response, fallback string) error {
	if resp.Code == "1" {
		return nil
	}
	if message := normalizeResponseMessage(resp.Msg); message != "" {
		return errors.New(message)
	}
	return errors.New(fallback)
}

func normalizeResponseMessage(message string) string {
	message = strings.TrimSpace(message)
	if message == "" {
		return ""
	}

	var items []string
	if strings.HasPrefix(message, "[") && strings.HasSuffix(message, "]") {
		if err := json.Unmarshal([]byte(message), &items); err == nil {
			normalized := make([]string, 0, len(items))
			for _, item := range items {
				if cleaned := normalizeValidationMessage(item); cleaned != "" {
					normalized = append(normalized, cleaned)
				}
			}
			return strings.Join(normalized, "; ")
		}
	}

	return normalizeValidationMessage(message)
}

func normalizeValidationMessage(message string) string {
	fields := strings.Fields(strings.TrimSpace(message))
	if len(fields) >= 2 && fields[0] == fields[1] {
		fields = append(fields[:1], fields[2:]...)
	}
	if len(fields) == 0 {
		return ""
	}
	return strings.Join(fields, " ")
}
