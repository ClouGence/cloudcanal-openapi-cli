package repl

import "strings"

type commandTextFunc func(*Shell) string
type commandRunFunc func(*Shell, []string) error

type commandSpec struct {
	name          string
	aliases       []string
	visible       bool
	visibleInHelp bool
	help          commandTextFunc
	usage         commandTextFunc
	run           commandRunFunc
	flags         []flagSpec
	nextArgs      []string
	children      []*commandSpec
}

var (
	boolValues   = []string{"true", "false"}
	outputValues = []string{"text", "json"}
	rootCommands = []*commandSpec{
		{
			name:          "jobs",
			visible:       true,
			visibleInHelp: true,
			help:          (*Shell).helpJobs,
			usage:         (*Shell).usageJobsGroup,
			children: []*commandSpec{
				{name: "list", visible: true, usage: (*Shell).usageJobsList, flags: []flagSpec{{name: "--name"}, {name: "--type"}, {name: "--desc"}, {name: "--source-id"}, {name: "--target-id"}}},
				{name: "create", visible: true, usage: (*Shell).usageJobCreate, flags: []flagSpec{{name: "--body"}, {name: "--body-file"}}},
				{name: "show", visible: true, usage: bindActionUsage("show", (*Shell).usageJobAction)},
				{name: "schema", visible: true, usage: bindActionUsage("schema", (*Shell).usageJobAction)},
				{name: "start", visible: true, usage: bindActionUsage("start", (*Shell).usageJobAction)},
				{name: "stop", visible: true, usage: bindActionUsage("stop", (*Shell).usageJobAction)},
				{name: "delete", visible: true, usage: bindActionUsage("delete", (*Shell).usageJobAction)},
				{name: "replay", visible: true, usage: (*Shell).usageJobReplay, flags: []flagSpec{{name: "--auto-start", values: boolValues}, {name: "--reset-to-created", values: boolValues}}},
				{name: "attach-incre-task", visible: true, usage: bindActionUsage("attach-incre-task", (*Shell).usageJobAction)},
				{name: "detach-incre-task", visible: true, usage: bindActionUsage("detach-incre-task", (*Shell).usageJobAction)},
				{name: "update-incre-pos", visible: true, usage: (*Shell).usageJobUpdateIncrePos, flags: []flagSpec{{name: "--body"}, {name: "--body-file"}}},
			},
		},
		{
			name:          "datasources",
			visible:       true,
			visibleInHelp: true,
			help:          (*Shell).helpDataSources,
			usage:         (*Shell).usageDataSources,
			children: []*commandSpec{
				{name: "list", visible: true, usage: (*Shell).usageDataSourcesList, flags: []flagSpec{{name: "--id"}, {name: "--type"}, {name: "--deploy-type"}, {name: "--host-type"}, {name: "--lifecycle"}}},
				{name: "add", visible: true, usage: (*Shell).usageDataSourceAdd, flags: []flagSpec{{name: "--body"}, {name: "--body-file"}, {name: "--security-file"}, {name: "--secret-file"}}},
				{name: "delete", visible: true, usage: bindActionUsage("delete", (*Shell).usageDataSourceAction)},
				{name: "show", visible: true, usage: (*Shell).usageDataSourceShow},
			},
		},
		{
			name:          "clusters",
			visible:       true,
			visibleInHelp: true,
			help:          (*Shell).helpClusters,
			usage:         (*Shell).usageClusters,
			children: []*commandSpec{
				{name: "list", visible: true, usage: (*Shell).usageClustersList, flags: []flagSpec{{name: "--name"}, {name: "--desc"}, {name: "--cloud"}, {name: "--region"}}},
			},
		},
		{
			name:          "workers",
			visible:       true,
			visibleInHelp: true,
			help:          (*Shell).helpWorkers,
			usage:         (*Shell).usageWorkers,
			children: []*commandSpec{
				{name: "list", visible: true, usage: (*Shell).usageWorkersList, flags: []flagSpec{{name: "--cluster-id"}, {name: "--source-id"}, {name: "--target-id"}}},
				{name: "start", visible: true, usage: bindActionUsage("start", (*Shell).usageWorkerAction)},
				{name: "stop", visible: true, usage: bindActionUsage("stop", (*Shell).usageWorkerAction)},
				{name: "delete", visible: true, usage: bindActionUsage("delete", (*Shell).usageWorkerAction)},
				{name: "modify-mem-oversold", visible: true, usage: (*Shell).usageWorkerModifyMemOverSold, flags: []flagSpec{{name: "--percent"}}},
				{name: "update-alert", visible: true, usage: (*Shell).usageWorkerUpdateAlert, flags: []flagSpec{{name: "--phone", values: boolValues}, {name: "--email", values: boolValues}, {name: "--im", values: boolValues}, {name: "--sms", values: boolValues}}},
			},
		},
		{
			name:          "consolejobs",
			visible:       true,
			visibleInHelp: true,
			help:          (*Shell).helpConsoleJobs,
			usage:         (*Shell).usageConsoleJobs,
			children: []*commandSpec{
				{name: "show", visible: true, usage: (*Shell).usageConsoleJobShow},
			},
		},
		{
			name:          "job-config",
			aliases:       []string{"jobconfig"},
			visible:       true,
			visibleInHelp: true,
			help:          (*Shell).helpJobConfig,
			usage:         (*Shell).usageJobConfig,
			children: []*commandSpec{
				{name: "specs", visible: true, usage: (*Shell).usageJobConfigSpecs, flags: []flagSpec{{name: "--type"}, {name: "--initial-sync", values: boolValues}, {name: "--short-term-sync", values: boolValues}}},
				{name: "transform-job-type", visible: true, usage: (*Shell).usageJobConfigTransform, flags: []flagSpec{{name: "--source-type"}, {name: "--target-type"}}},
			},
		},
		{
			name:          "schemas",
			aliases:       []string{"schema"},
			visible:       true,
			visibleInHelp: true,
			help:          (*Shell).helpSchemas,
			usage:         (*Shell).usageSchemas,
			children: []*commandSpec{
				{name: "list-trans-objs-by-meta", visible: true, usage: (*Shell).usageSchemas, flags: []flagSpec{{name: "--src-db"}, {name: "--src-schema"}, {name: "--src-trans-obj"}, {name: "--dst-db"}, {name: "--dst-schema"}, {name: "--dst-tran-obj"}}},
			},
		},
		{
			name:          "config",
			visible:       true,
			visibleInHelp: true,
			help:          (*Shell).helpConfig,
			usage:         (*Shell).usageConfig,
			children: []*commandSpec{
				{name: "show", visible: true, usage: (*Shell).usageConfigShow},
				{name: "init", visible: true, usage: (*Shell).usageConfigInit},
				newProfilesCommand(),
				newLanguageCommand("lang", true),
			},
		},
		{
			name:          "version",
			visible:       true,
			visibleInHelp: true,
			help:          (*Shell).helpVersion,
			usage:         (*Shell).usageVersion,
		},
		newLanguageCommand("lang", false, "language"),
		{
			name:  "completion",
			help:  (*Shell).helpCompletion,
			usage: (*Shell).usageCompletion,
			children: []*commandSpec{
				{name: "zsh", visible: true, usage: (*Shell).usageCompletion},
				{name: "bash", visible: true, usage: (*Shell).usageCompletion},
			},
		},
		{
			name:    "clear",
			aliases: []string{"cls"},
		},
		{
			name: "__complete",
		},
	}
)

func bindActionUsage(action string, fn func(*Shell, string) string) commandTextFunc {
	return func(s *Shell) string {
		return fn(s, action)
	}
}

func newLanguageCommand(name string, visible bool, aliases ...string) *commandSpec {
	return &commandSpec{
		name:    name,
		aliases: aliases,
		visible: visible,
		help:    (*Shell).helpLanguage,
		usage:   (*Shell).usageConfigLang,
		children: []*commandSpec{
			{name: "show", visible: true},
			{name: "set", visible: true, nextArgs: []string{"en", "zh"}},
		},
	}
}

func newProfilesCommand() *commandSpec {
	return &commandSpec{
		name:    "profiles",
		visible: true,
		help:    (*Shell).helpProfiles,
		usage:   (*Shell).usageConfigProfiles,
		children: []*commandSpec{
			{name: "list", visible: true, usage: (*Shell).usageConfigProfiles},
			{name: "use", visible: true, usage: (*Shell).usageConfigProfiles},
			{name: "add", visible: true, usage: (*Shell).usageConfigProfiles},
			{name: "remove", visible: true, usage: (*Shell).usageConfigProfiles},
		},
	}
}

func init() {
	mustSetCommandRun("jobs", (*Shell).handleJobs)
	mustSetCommandRun("datasources", (*Shell).handleDataSources)
	mustSetCommandRun("clusters", (*Shell).handleClusters)
	mustSetCommandRun("workers", (*Shell).handleWorkers)
	mustSetCommandRun("consolejobs", (*Shell).handleConsoleJobs)
	mustSetCommandRun("job-config", (*Shell).handleJobConfig)
	mustSetCommandRun("schemas", (*Shell).handleSchemas)
	mustSetCommandRun("config", (*Shell).handleConfig)
	mustSetCommandRun("lang", (*Shell).handleLang)
	mustSetCommandRun("version", (*Shell).handleVersion)
	mustSetCommandRun("completion", (*Shell).handleCompletion)
	mustSetCommandRun("clear", runClearScreen)
	mustSetCommandRun("__complete", runHiddenCompletion)

	mustSetSubcommandRun("jobs", "list", (*Shell).runJobsList)
	mustSetSubcommandRun("jobs", "create", (*Shell).runJobsCreate)
	mustSetSubcommandRun("jobs", "show", (*Shell).runJobsShow)
	mustSetSubcommandRun("jobs", "schema", (*Shell).runJobsSchema)
	mustSetSubcommandRun("jobs", "start", (*Shell).runJobsStart)
	mustSetSubcommandRun("jobs", "stop", (*Shell).runJobsStop)
	mustSetSubcommandRun("jobs", "delete", (*Shell).runJobsDelete)
	mustSetSubcommandRun("jobs", "replay", (*Shell).runJobsReplay)
	mustSetSubcommandRun("jobs", "attach-incre-task", (*Shell).runJobsAttachIncreTask)
	mustSetSubcommandRun("jobs", "detach-incre-task", (*Shell).runJobsDetachIncreTask)
	mustSetSubcommandRun("jobs", "update-incre-pos", (*Shell).runJobsUpdateIncrePos)

	mustSetSubcommandRun("datasources", "list", (*Shell).runDataSourcesList)
	mustSetSubcommandRun("datasources", "add", (*Shell).runDataSourcesAdd)
	mustSetSubcommandRun("datasources", "delete", (*Shell).runDataSourcesDelete)
	mustSetSubcommandRun("datasources", "show", (*Shell).runDataSourcesShow)

	mustSetSubcommandRun("clusters", "list", (*Shell).runClustersList)

	mustSetSubcommandRun("workers", "list", (*Shell).runWorkersList)
	mustSetSubcommandRun("workers", "start", (*Shell).runWorkersStart)
	mustSetSubcommandRun("workers", "stop", (*Shell).runWorkersStop)
	mustSetSubcommandRun("workers", "delete", (*Shell).runWorkersDelete)
	mustSetSubcommandRun("workers", "modify-mem-oversold", (*Shell).runWorkersModifyMemOversold)
	mustSetSubcommandRun("workers", "update-alert", (*Shell).runWorkersUpdateAlert)

	mustSetSubcommandRun("consolejobs", "show", (*Shell).runConsoleJobsShow)

	mustSetSubcommandRun("job-config", "specs", (*Shell).runJobConfigSpecs)
	mustSetSubcommandRun("job-config", "transform-job-type", (*Shell).runJobConfigTransformJobType)

	mustSetSubcommandRun("schemas", "list-trans-objs-by-meta", (*Shell).runSchemasListTransferObjects)

	mustSetCommandRunPath([]string{"config", "show"}, (*Shell).runConfigShow)
	mustSetCommandRunPath([]string{"config", "init"}, (*Shell).runConfigInit)
	mustSetCommandRunPath([]string{"config", "profiles", "list"}, (*Shell).runProfilesList)
	mustSetCommandRunPath([]string{"config", "profiles", "use"}, (*Shell).runProfilesUse)
	mustSetCommandRunPath([]string{"config", "profiles", "add"}, (*Shell).runProfilesAdd)
	mustSetCommandRunPath([]string{"config", "profiles", "remove"}, (*Shell).runProfilesRemove)
	mustSetCommandRunPath([]string{"config", "lang", "show"}, (*Shell).runLanguageShow)
	mustSetCommandRunPath([]string{"config", "lang", "set"}, (*Shell).runLanguageSet)

	mustSetCommandRunPath([]string{"lang", "show"}, (*Shell).runLanguageShow)
	mustSetCommandRunPath([]string{"lang", "set"}, (*Shell).runLanguageSet)

	mustSetCommandRunPath([]string{"version"}, (*Shell).runVersion)

	mustSetCommandRunPath([]string{"completion", "zsh"}, (*Shell).runCompletionZsh)
	mustSetCommandRunPath([]string{"completion", "bash"}, (*Shell).runCompletionBash)
}

func mustSetCommandRun(name string, run commandRunFunc) {
	mustSetCommandRunPath([]string{name}, run)
}

func mustSetSubcommandRun(root string, child string, run commandRunFunc) {
	mustSetCommandRunPath([]string{root, child}, run)
}

func mustSetCommandRunPath(path []string, run commandRunFunc) {
	spec, consumed := findCommandPath(path)
	if spec == nil || consumed != len(path) {
		panic("command not found: " + strings.Join(path, " "))
	}
	spec.run = run
}

func runClearScreen(shell *Shell, _ []string) error {
	shell.io.ClearScreen()
	return nil
}

func runHiddenCompletion(shell *Shell, tokens []string) error {
	shell.printHiddenCompletions(tokens[1:])
	return nil
}

func (s *Shell) dispatchRegisteredCommand(tokens []string) error {
	if len(tokens) == 0 {
		return nil
	}
	spec, consumed := findCommandPath(tokens)
	if spec == nil {
		return nil
	}
	parent := findCommandParent(tokens, consumed)

	if len(tokens) == consumed {
		if len(spec.children) > 0 {
			s.io.Println(commandUsageText(s, spec, parent))
			return nil
		}
		if spec.run != nil {
			return spec.run(s, tokens)
		}
		s.io.Println(commandUsageText(s, spec, parent))
		return nil
	}

	if len(spec.children) > 0 {
		s.printUnknownSubcommand(commandPathName(tokens, consumed), tokens[consumed], visibleCommandNames(spec.children), commandUsageText(s, spec, parent))
		return nil
	}
	if spec.run != nil {
		return spec.run(s, tokens)
	}
	s.io.Println(commandUsageText(s, spec, parent))
	return nil
}

func visibleTopLevelCommands() []string {
	return append([]string{"help"}, visibleCommandNames(rootCommands)...)
}

func visibleHelpTopics() []string {
	topics := make([]string, 0, len(rootCommands))
	for _, spec := range rootCommands {
		if spec.visibleInHelp {
			topics = append(topics, spec.name)
		}
	}
	return topics
}

func visibleCommandNames(specs []*commandSpec) []string {
	names := make([]string, 0, len(specs))
	for _, spec := range specs {
		if spec.visible {
			names = append(names, spec.name)
		}
	}
	return names
}

func findRootCommand(token string) *commandSpec {
	return findCommand(rootCommands, token)
}

func findChildCommand(parent *commandSpec, token string) *commandSpec {
	if parent == nil {
		return nil
	}
	return findCommand(parent.children, token)
}

func subcommandCandidates(path ...string) []string {
	spec, consumed := findCommandPath(path)
	if spec == nil || consumed != len(path) {
		return nil
	}
	return visibleCommandNames(spec.children)
}

func findCommandPath(tokens []string) (*commandSpec, int) {
	if len(tokens) == 0 {
		return nil, 0
	}

	spec := findRootCommand(tokens[0])
	if spec == nil {
		return nil, 0
	}

	consumed := 1
	for consumed < len(tokens) {
		child := findChildCommand(spec, tokens[consumed])
		if child == nil {
			break
		}
		spec = child
		consumed++
	}
	return spec, consumed
}

func findCommandParent(tokens []string, consumed int) *commandSpec {
	if consumed <= 1 {
		return nil
	}
	parent, _ := findCommandPath(tokens[:consumed-1])
	return parent
}

func commandPathName(tokens []string, consumed int) string {
	names := make([]string, 0, consumed)
	spec := findRootCommand(tokens[0])
	if spec == nil {
		return ""
	}
	names = append(names, spec.name)
	for i := 1; i < consumed; i++ {
		spec = findChildCommand(spec, tokens[i])
		if spec == nil {
			break
		}
		names = append(names, spec.name)
	}
	return strings.Join(names, " ")
}

func findCommand(specs []*commandSpec, token string) *commandSpec {
	lower := strings.ToLower(strings.TrimSpace(token))
	if lower == "" {
		return nil
	}

	for _, spec := range specs {
		if strings.EqualFold(spec.name, lower) {
			return spec
		}
		for _, alias := range spec.aliases {
			if strings.EqualFold(alias, lower) {
				return spec
			}
		}
	}
	return nil
}

func commandHelpText(shell *Shell, spec *commandSpec) string {
	if spec == nil {
		return shell.helpOverview()
	}
	if spec.help != nil {
		return spec.help(shell)
	}
	if spec.usage != nil {
		return spec.usage(shell)
	}
	return shell.helpOverview()
}

func commandUsageText(shell *Shell, spec *commandSpec, parent *commandSpec) string {
	if spec == nil {
		return shell.helpOverview()
	}
	if spec.usage != nil {
		return spec.usage(shell)
	}
	return commandUsageOrHelpText(shell, spec, parent)
}

func canRenderHelp(spec *commandSpec) bool {
	if spec == nil {
		return false
	}
	return spec.visibleInHelp || spec.help != nil || spec.usage != nil
}

func commandUsageOrHelpText(shell *Shell, spec *commandSpec, parent *commandSpec) string {
	if spec == nil {
		return commandHelpText(shell, parent)
	}
	if len(spec.children) > 0 && spec.help != nil {
		return spec.help(shell)
	}
	if spec.usage != nil {
		return spec.usage(shell)
	}
	if spec.help != nil {
		return spec.help(shell)
	}
	return commandHelpText(shell, parent)
}
