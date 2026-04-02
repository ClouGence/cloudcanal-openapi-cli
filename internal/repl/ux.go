package repl

import (
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/i18n"
	"strings"
)

func isHelpToken(token string) bool {
	switch strings.ToLower(strings.TrimSpace(token)) {
	case "help", "-h", "--help":
		return true
	default:
		return false
	}
}

func suggestCandidate(input string, candidates []string) string {
	query := strings.ToLower(strings.TrimSpace(input))
	if query == "" || len(candidates) == 0 {
		return ""
	}

	best := ""
	bestDistance := -1
	for _, candidate := range candidates {
		candidateLower := strings.ToLower(candidate)
		distance := levenshteinDistance(query, candidateLower)
		if strings.HasPrefix(candidateLower, query) || strings.HasPrefix(query, candidateLower) {
			distance = min(distance, 1)
		}
		if bestDistance == -1 || distance < bestDistance || (distance == bestDistance && candidateLower < strings.ToLower(best)) {
			best = candidate
			bestDistance = distance
		}
	}

	maxDistance := 2
	if len(query) >= 7 {
		maxDistance = 3
	}
	if bestDistance > maxDistance {
		return ""
	}
	return best
}

func levenshteinDistance(a string, b string) int {
	if a == b {
		return 0
	}
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	prev := make([]int, len(b)+1)
	curr := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}

	for i := 1; i <= len(a); i++ {
		curr[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			curr[j] = min(
				min(curr[j-1]+1, prev[j]+1),
				prev[j-1]+cost,
			)
		}
		prev, curr = curr, prev
	}

	return prev[len(b)]
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s *Shell) unknownHelpText(topic string) string {
	lines := []string{
		i18n.T("common.errorPrefix", i18n.T("common.unknownHelpTopic", topic)),
	}
	if suggestion := suggestCandidate(topic, visibleHelpTopics()); suggestion != "" {
		lines = append(lines, i18n.T("common.didYouMean", "help "+suggestion))
	}
	lines = append(lines, "", s.helpOverview())
	return strings.Join(lines, "\n")
}

func (s *Shell) printUnknownCommand(command string) {
	s.io.Println(i18n.T("common.unknownCommand", command))
	if suggestion := suggestCandidate(command, visibleTopLevelCommands()); suggestion != "" {
		s.io.Println(i18n.T("common.didYouMean", suggestion))
	}
	s.io.Println(i18n.T("common.useHelp"))
}

func (s *Shell) printUnknownSubcommand(group string, command string, candidates []string, usage string) {
	s.io.Println(i18n.T("common.unknownSubcommand", group, command))
	if suggestion := suggestCandidate(command, candidates); suggestion != "" {
		s.io.Println(i18n.T("common.didYouMean", group+" "+suggestion))
	}
	s.io.Println(usage)
}

func RenderCommandHelp(tokens []string) (string, bool) {
	if len(tokens) == 0 {
		return "", false
	}

	shell := &Shell{}
	if strings.EqualFold(tokens[0], "help") {
		return shell.renderHelp(tokens[1:]), true
	}
	if isHelpToken(tokens[0]) {
		return shell.renderHelp(nil), true
	}

	if !isHelpToken(tokens[len(tokens)-1]) {
		return "", false
	}

	path := tokens[:len(tokens)-1]
	spec, consumed := findCommandPath(path)
	if spec == nil || consumed != len(path) {
		return "", false
	}
	parent := findCommandParent(path, consumed)
	return commandUsageOrHelpText(shell, spec, parent), true
}
