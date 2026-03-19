package repl

import (
	"cloudcanal-openapi-cli/internal/i18n"
	"strings"
)

var (
	jobsSubcommands       = []string{"list", "create", "show", "schema", "start", "stop", "delete", "replay", "attach-incre-task", "detach-incre-task", "update-incre-pos"}
	dataSourceSubcommands = []string{"list", "add", "delete", "show"}
	clusterSubcommands    = []string{"list"}
	workerSubcommands     = []string{"list", "start", "stop", "delete", "modify-mem-oversold", "update-alert"}
	consoleJobSubcommands = []string{"show"}
	jobConfigSubcommands  = []string{"specs", "transform-job-type"}
	schemaSubcommands     = []string{"list-trans-objs-by-meta"}
	configSubcommands     = []string{"show", "init", "lang"}
	langSubcommands       = []string{"show", "set"}
	completionSubcommands = []string{"zsh", "bash"}
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
	if suggestion := suggestCandidate(topic, visibleHelpTopics); suggestion != "" {
		lines = append(lines, i18n.T("common.didYouMean", "help "+suggestion))
	}
	lines = append(lines, "", s.helpOverview())
	return strings.Join(lines, "\n")
}

func (s *Shell) printUnknownCommand(command string) {
	s.io.Println(i18n.T("common.unknownCommand", command))
	if suggestion := suggestCandidate(command, visibleTopLevelCommands); suggestion != "" {
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

	root := strings.ToLower(tokens[0])
	if len(tokens) >= 2 && isHelpToken(tokens[1]) {
		switch root {
		case "jobs":
			return shell.helpJobs(), true
		case "datasources":
			return shell.helpDataSources(), true
		case "clusters":
			return shell.helpClusters(), true
		case "workers":
			return shell.helpWorkers(), true
		case "consolejobs":
			return shell.helpConsoleJobs(), true
		case "job-config", "jobconfig":
			return shell.helpJobConfig(), true
		case "schemas", "schema":
			return shell.helpSchemas(), true
		case "config":
			return shell.helpConfig(), true
		case "lang", "language":
			return shell.helpLanguage(), true
		case "completion":
			return shell.helpCompletion(), true
		}
	}

	if len(tokens) < 3 || !isHelpToken(tokens[2]) {
		return "", false
	}

	switch root {
	case "jobs":
		switch strings.ToLower(tokens[1]) {
		case "list":
			return shell.usageJobsList(), true
		case "create":
			return shell.usageJobCreate(), true
		case "show", "schema", "start", "stop", "delete":
			return shell.usageJobAction(strings.ToLower(tokens[1])), true
		case "replay":
			return shell.usageJobReplay(), true
		case "attach-incre-task", "detach-incre-task":
			return shell.usageJobAction(strings.ToLower(tokens[1])), true
		case "update-incre-pos":
			return shell.usageJobUpdateIncrePos(), true
		default:
			return shell.helpJobs(), true
		}
	case "datasources":
		switch strings.ToLower(tokens[1]) {
		case "list":
			return shell.usageDataSourcesList(), true
		case "add":
			return shell.usageDataSourceAdd(), true
		case "delete":
			return shell.usageDataSourceAction("delete"), true
		case "show":
			return shell.usageDataSourceShow(), true
		default:
			return shell.helpDataSources(), true
		}
	case "clusters":
		if strings.EqualFold(tokens[1], "list") {
			return shell.usageClustersList(), true
		}
		return shell.helpClusters(), true
	case "workers":
		switch strings.ToLower(tokens[1]) {
		case "list":
			return shell.usageWorkersList(), true
		case "start", "stop", "delete":
			return shell.usageWorkerAction(strings.ToLower(tokens[1])), true
		case "modify-mem-oversold":
			return shell.usageWorkerModifyMemOverSold(), true
		case "update-alert":
			return shell.usageWorkerUpdateAlert(), true
		default:
			return shell.helpWorkers(), true
		}
	case "consolejobs":
		if strings.EqualFold(tokens[1], "show") {
			return shell.usageConsoleJobShow(), true
		}
		return shell.helpConsoleJobs(), true
	case "job-config", "jobconfig":
		if strings.EqualFold(tokens[1], "specs") {
			return shell.usageJobConfigSpecs(), true
		}
		if strings.EqualFold(tokens[1], "transform-job-type") {
			return shell.usageJobConfigTransform(), true
		}
		return shell.helpJobConfig(), true
	case "schemas", "schema":
		if strings.EqualFold(tokens[1], "list-trans-objs-by-meta") {
			return shell.usageSchemas(), true
		}
		return shell.helpSchemas(), true
	case "config":
		switch strings.ToLower(tokens[1]) {
		case "show":
			return shell.usageConfigShow(), true
		case "init":
			return shell.usageConfigInit(), true
		case "lang":
			return shell.helpLanguage(), true
		default:
			return shell.helpConfig(), true
		}
	case "lang", "language":
		return shell.helpLanguage(), true
	case "completion":
		if strings.EqualFold(tokens[1], "zsh") || strings.EqualFold(tokens[1], "bash") {
			return shell.usageCompletion(), true
		}
		return shell.helpCompletion(), true
	default:
		return "", false
	}
}
