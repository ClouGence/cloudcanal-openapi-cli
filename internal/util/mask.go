package util

import (
	"regexp"
	"strings"
)

var sensitiveAssignmentPattern = regexp.MustCompile(`(?i)\b(access[-_]?key|secret[-_]?key|password|passwd|pwd|token|api[-_]?key)\s*=\s*([^\s,)]+)`)
var slashCredentialPattern = regexp.MustCompile(`(^|[\s(])([^,\s)]+?\s*/\s*)([^,\s)]+)`)
var singleTokenSecretPattern = regexp.MustCompile(`\(((?i:sk-|ak-|token-)[^)\s]+)\)`)

func MaskSecret(value string) string {
	if value == "" {
		return "-"
	}
	if len(value) <= 8 {
		return strings.Repeat("*", len(value))
	}
	return value[:4] + strings.Repeat("*", len(value)-8) + value[len(value)-4:]
}

func MaskSensitiveText(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	masked := sensitiveAssignmentPattern.ReplaceAllStringFunc(value, func(match string) string {
		index := strings.Index(match, "=")
		if index < 0 {
			return match
		}
		key := match[:index+1]
		secret := strings.TrimSpace(match[index+1:])
		return key + MaskSecret(secret)
	})

	masked = slashCredentialPattern.ReplaceAllStringFunc(masked, func(match string) string {
		parts := slashCredentialPattern.FindStringSubmatch(match)
		if len(parts) != 4 {
			return match
		}
		return parts[1] + parts[2] + MaskSecret(parts[3])
	})

	masked = singleTokenSecretPattern.ReplaceAllStringFunc(masked, func(match string) string {
		parts := singleTokenSecretPattern.FindStringSubmatch(match)
		if len(parts) != 2 {
			return match
		}
		return "(" + MaskSecret(parts[1]) + ")"
	})

	return masked
}
