package repl

import (
	"fmt"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

type flagSpec struct {
	name   string
	values []string
}

const CompletionEnvVar = "CLOUDCANAL_INTERNAL_COMPLETE"

var (
	visibleTopLevelCommands = []string{
		"help",
		"jobs",
		"datasources",
		"clusters",
		"workers",
		"consolejobs",
		"job-config",
		"schemas",
		"config",
	}
	visibleReplOnlyCommands = []string{"exit"}
	visibleHelpTopics       = []string{
		"jobs",
		"datasources",
		"clusters",
		"workers",
		"consolejobs",
		"job-config",
		"schemas",
		"config",
	}
	boolValues   = []string{"true", "false"}
	outputValues = []string{"text", "json"}
)

func RenderCompletionScript(args []string) (string, error) {
	if len(args) == 0 || len(args) > 2 {
		return "", fmt.Errorf("usage: completion <zsh|bash> [command-name]")
	}

	commandName := "cloudcanal"
	if len(args) == 2 && strings.TrimSpace(args[1]) != "" {
		commandName = strings.TrimSpace(args[1])
	}

	switch strings.ToLower(args[0]) {
	case "zsh":
		return renderZshCompletionScript(commandName), nil
	case "bash":
		return renderBashCompletionScript(commandName), nil
	default:
		return "", fmt.Errorf("unsupported shell: %s", args[0])
	}
}

func CompletionCandidates(args []string, replMode bool) []string {
	context, prefix := completionContextFromArgs(args)
	return completeContext(context, prefix, replMode)
}

func (s *Shell) completeLine(line string) []string {
	state := parseCompletionLine(line)
	candidates := completeContext(state.context, state.prefix, true)
	base := line[:state.tokenStart]
	results := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		results = append(results, base+candidate)
	}
	return results
}

func (s *Shell) handleCompletion(tokens []string) error {
	if len(tokens) < 2 || len(tokens) > 3 {
		s.io.Println(s.usageCompletion())
		return nil
	}

	script, err := RenderCompletionScript(tokens[1:])
	if err != nil {
		return err
	}
	s.io.Println(script)
	return nil
}

func (s *Shell) printHiddenCompletions(args []string) {
	for _, candidate := range CompletionCandidates(args, false) {
		s.io.Println(candidate)
	}
}

func completeContext(context []string, prefix string, replMode bool) []string {
	if len(context) == 0 {
		if name, valuePrefix, ok := splitInlineFlag(prefix); ok && name == "--output" {
			return prependInlineFlag(name, matchCandidates(outputValues, valuePrefix))
		}
		candidates := append([]string{}, visibleTopLevelCommands...)
		if replMode {
			candidates = append(candidates, visibleReplOnlyCommands...)
		}
		if prefix == "" || strings.HasPrefix(prefix, "--") {
			candidates = append(candidates, "--help", "--output")
		}
		return matchCandidates(candidates, prefix)
	}

	root := strings.ToLower(context[0])
	switch root {
	case "help":
		return matchCandidates(visibleHelpTopics, prefix)
	case "completion":
		if len(context) == 1 {
			return matchCandidates(append(append([]string{}, completionSubcommands...), "--help"), prefix)
		}
		return nil
	case "lang", "language":
		if len(context) == 1 {
			return matchCandidates(append(append([]string{}, langSubcommands...), "--help"), prefix)
		}
		if strings.EqualFold(context[1], "set") {
			return matchCandidates([]string{"en", "zh"}, prefix)
		}
		return nil
	case "config":
		if len(context) == 1 {
			return matchCandidates(append(append([]string{}, configSubcommands...), "--help"), prefix)
		}
		if strings.EqualFold(context[1], "lang") {
			if len(context) == 2 {
				return matchCandidates(append(append([]string{}, langSubcommands...), "--help"), prefix)
			}
			if len(context) == 3 && strings.EqualFold(context[2], "set") {
				return matchCandidates([]string{"en", "zh"}, prefix)
			}
		}
		return nil
	case "jobs":
		if len(context) == 1 {
			return matchCandidates(append(append([]string{}, jobsSubcommands...), "--help"), prefix)
		}
		switch strings.ToLower(context[1]) {
		case "list":
			return completeFlags(context[2:], prefix, withGlobalFlags([]flagSpec{
				{name: "--name"},
				{name: "--type"},
				{name: "--desc"},
				{name: "--source-id"},
				{name: "--target-id"},
			}))
		case "create", "update-incre-pos":
			return completeFlags(context[2:], prefix, withGlobalFlags([]flagSpec{
				{name: "--body"},
				{name: "--body-file"},
			}))
		case "replay":
			return completeFlags(context[2:], prefix, withGlobalFlags([]flagSpec{
				{name: "--auto-start", values: boolValues},
				{name: "--reset-to-created", values: boolValues},
			}))
		}
		return nil
	case "datasources":
		if len(context) == 1 {
			return matchCandidates(append(append([]string{}, dataSourceSubcommands...), "--help"), prefix)
		}
		if strings.EqualFold(context[1], "list") {
			return completeFlags(context[2:], prefix, withGlobalFlags([]flagSpec{
				{name: "--id"},
				{name: "--type"},
				{name: "--deploy-type"},
				{name: "--host-type"},
				{name: "--lifecycle"},
			}))
		}
		if strings.EqualFold(context[1], "add") {
			return completeFlags(context[2:], prefix, withGlobalFlags([]flagSpec{
				{name: "--body"},
				{name: "--body-file"},
				{name: "--security-file"},
				{name: "--secret-file"},
			}))
		}
		return nil
	case "clusters":
		if len(context) == 1 {
			return matchCandidates(append(append([]string{}, clusterSubcommands...), "--help"), prefix)
		}
		if strings.EqualFold(context[1], "list") {
			return completeFlags(context[2:], prefix, withGlobalFlags([]flagSpec{
				{name: "--name"},
				{name: "--desc"},
				{name: "--cloud"},
				{name: "--region"},
			}))
		}
		return nil
	case "workers":
		if len(context) == 1 {
			return matchCandidates(append(append([]string{}, workerSubcommands...), "--help"), prefix)
		}
		if strings.EqualFold(context[1], "list") {
			return completeFlags(context[2:], prefix, withGlobalFlags([]flagSpec{
				{name: "--cluster-id"},
				{name: "--source-id"},
				{name: "--target-id"},
			}))
		}
		if strings.EqualFold(context[1], "modify-mem-oversold") {
			return completeFlags(context[2:], prefix, withGlobalFlags([]flagSpec{
				{name: "--percent"},
			}))
		}
		if strings.EqualFold(context[1], "update-alert") {
			return completeFlags(context[2:], prefix, withGlobalFlags([]flagSpec{
				{name: "--phone", values: boolValues},
				{name: "--email", values: boolValues},
				{name: "--im", values: boolValues},
				{name: "--sms", values: boolValues},
			}))
		}
		return nil
	case "consolejobs":
		if len(context) == 1 {
			return matchCandidates(append(append([]string{}, consoleJobSubcommands...), "--help"), prefix)
		}
		return nil
	case "job-config", "jobconfig":
		if len(context) == 1 {
			return matchCandidates(append(append([]string{}, jobConfigSubcommands...), "--help"), prefix)
		}
		if strings.EqualFold(context[1], "specs") {
			return completeFlags(context[2:], prefix, withGlobalFlags([]flagSpec{
				{name: "--type"},
				{name: "--initial-sync", values: boolValues},
				{name: "--short-term-sync", values: boolValues},
			}))
		}
		if strings.EqualFold(context[1], "transform-job-type") {
			return completeFlags(context[2:], prefix, withGlobalFlags([]flagSpec{
				{name: "--source-type"},
				{name: "--target-type"},
			}))
		}
		return nil
	case "schemas", "schema":
		if len(context) == 1 {
			return matchCandidates(append(append([]string{}, schemaSubcommands...), "--help"), prefix)
		}
		if strings.EqualFold(context[1], "list-trans-objs-by-meta") {
			return completeFlags(context[2:], prefix, withGlobalFlags([]flagSpec{
				{name: "--src-db"},
				{name: "--src-schema"},
				{name: "--src-trans-obj"},
				{name: "--dst-db"},
				{name: "--dst-schema"},
				{name: "--dst-tran-obj"},
			}))
		}
		return nil
	default:
		return nil
	}
}

func completeFlags(args []string, prefix string, specs []flagSpec) []string {
	if len(args) > 0 {
		if values, handled := valuesForPreviousFlag(args[len(args)-1], prefix, specs); handled {
			return values
		}
	}

	if name, valuePrefix, ok := splitInlineFlag(prefix); ok {
		for _, spec := range specs {
			if spec.name == name {
				return prependInlineFlag(name, matchCandidates(spec.values, valuePrefix))
			}
		}
		return nil
	}

	if prefix == "" || strings.HasPrefix(prefix, "--") {
		used := usedFlags(args)
		candidates := make([]string, 0, len(specs))
		for _, spec := range specs {
			if !used[spec.name] {
				candidates = append(candidates, spec.name)
			}
		}
		return matchCandidates(candidates, prefix)
	}

	return nil
}

func withGlobalFlags(specs []flagSpec) []flagSpec {
	combined := make([]flagSpec, 0, len(specs)+2)
	combined = append(combined, specs...)
	combined = append(combined, flagSpec{name: "--help"})
	combined = append(combined, flagSpec{name: "--output", values: outputValues})
	return combined
}

func valuesForPreviousFlag(previousToken string, prefix string, specs []flagSpec) ([]string, bool) {
	if strings.HasPrefix(prefix, "--") {
		return nil, false
	}
	for _, spec := range specs {
		if spec.name == previousToken {
			if len(spec.values) == 0 {
				return nil, true
			}
			return matchCandidates(spec.values, prefix), true
		}
	}
	return nil, false
}

func usedFlags(args []string) map[string]bool {
	used := make(map[string]bool, len(args))
	for _, arg := range args {
		if !strings.HasPrefix(arg, "--") {
			continue
		}
		name := arg
		if head, _, ok := strings.Cut(arg, "="); ok {
			name = head
		}
		used[name] = true
	}
	return used
}

func splitInlineFlag(prefix string) (string, string, bool) {
	if !strings.HasPrefix(prefix, "--") || !strings.Contains(prefix, "=") {
		return "", "", false
	}
	name, valuePrefix, ok := strings.Cut(prefix, "=")
	if !ok {
		return "", "", false
	}
	return name, valuePrefix, true
}

func prependInlineFlag(name string, values []string) []string {
	results := make([]string, 0, len(values))
	for _, value := range values {
		results = append(results, name+"="+value)
	}
	return results
}

func matchCandidates(candidates []string, prefix string) []string {
	if len(candidates) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(candidates))
	results := make([]string, 0, len(candidates))
	lowerPrefix := strings.ToLower(prefix)
	for _, candidate := range candidates {
		if prefix != "" && !strings.HasPrefix(strings.ToLower(candidate), lowerPrefix) {
			continue
		}
		if _, exists := seen[candidate]; exists {
			continue
		}
		seen[candidate] = struct{}{}
		results = append(results, candidate)
	}
	slices.Sort(results)
	return results
}

func completionContextFromArgs(args []string) ([]string, string) {
	if len(args) == 0 {
		return nil, ""
	}
	context := append([]string(nil), args[:len(args)-1]...)
	return context, args[len(args)-1]
}

type completionState struct {
	context    []string
	prefix     string
	tokenStart int
}

func parseCompletionLine(line string) completionState {
	var (
		context    []string
		current    strings.Builder
		quote      rune
		escaped    bool
		tokenStart = -1
	)

	for index, r := range line {
		switch {
		case escaped:
			if tokenStart < 0 {
				tokenStart = index
			}
			current.WriteRune(r)
			escaped = false
		case r == '\\':
			escaped = true
		case quote != 0:
			if r == quote {
				quote = 0
				continue
			}
			current.WriteRune(r)
		case r == '"' || r == '\'':
			if tokenStart < 0 {
				tokenStart = index + utf8.RuneLen(r)
			}
			quote = r
		case unicode.IsSpace(r):
			if tokenStart >= 0 {
				context = append(context, current.String())
				current.Reset()
				tokenStart = -1
			}
		default:
			if tokenStart < 0 {
				tokenStart = index
			}
			current.WriteRune(r)
		}
	}

	if tokenStart < 0 {
		return completionState{
			context:    context,
			prefix:     "",
			tokenStart: len(line),
		}
	}

	return completionState{
		context:    context,
		prefix:     current.String(),
		tokenStart: tokenStart,
	}
}

func renderZshCompletionScript(commandName string) string {
	return fmt.Sprintf(`#compdef %s

_%s() {
  local -a args completions
  local i

  args=()
  for ((i=2; i<CURRENT; i++)); do
    args+=("${words[i]}")
  done
  args+=("$PREFIX")

  completions=("${(@f)$(%s=1 "%s" "${args[@]}")}")
  if (( ${#completions[@]} > 0 )); then
    compadd -Q -- ${completions[@]}
  fi
}

compdef _%s %s
`, commandName, commandName, CompletionEnvVar, commandName, commandName, commandName)
}

func renderBashCompletionScript(commandName string) string {
	return fmt.Sprintf(`_%s_completion() {
  local -a args
  local i
  local cur

  args=()
  for ((i=1; i<COMP_CWORD; i++)); do
    args+=("${COMP_WORDS[i]}")
  done
  cur="${COMP_WORDS[COMP_CWORD]}"
  args+=("$cur")

  local IFS=$'\n'
  COMPREPLY=($(%s=1 "%s" "${args[@]}"))
}

complete -F _%s_completion %s
`, commandName, CompletionEnvVar, commandName, commandName, commandName)
}
