package repl

import (
	"cloudcanal-openapi-cli/internal/i18n"
	"fmt"
	"strings"
)

func (s *Shell) printHelp(args []string) {
	topic := ""
	if len(args) > 0 {
		topic = strings.ToLower(args[0])
	}

	switch topic {
	case "", "overview":
		s.io.Println(s.helpOverview())
	case "jobs":
		s.io.Println(s.helpJobs())
	case "datasources":
		s.io.Println(s.helpDataSources())
	case "clusters":
		s.io.Println(s.helpClusters())
	case "workers":
		s.io.Println(s.helpWorkers())
	case "consolejobs":
		s.io.Println(s.helpConsoleJobs())
	case "job-config", "jobconfig":
		s.io.Println(s.helpJobConfig())
	case "config":
		s.io.Println(s.helpConfig())
	case "lang", "language":
		s.io.Println(s.helpLanguage())
	case "completion":
		s.io.Println(s.helpCompletion())
	default:
		s.io.Println(s.helpOverview())
	}
}

func (s *Shell) helpOverview() string {
	if s.isChinese() {
		return strings.TrimSpace(`
CloudCanal CLI 帮助

帮助主题：
  help jobs         查看数据任务命令和参数说明
  help datasources  查看数据源命令和参数说明
  help clusters     查看集群命令和参数说明
  help workers      查看机器命令和参数说明
  help consolejobs  查看 ConsoleJob 命令说明
  help job-config   查看数据任务规格命令说明
  help config       查看配置命令说明
  help lang         查看语言切换命令说明

常用命令：
  jobs list         列出数据任务
  datasources list  列出数据源
  clusters list     列出集群
  workers list      列出机器
  consolejobs show  查看 ConsoleJob 详情
  job-config specs  查看任务规格
  config show       查看当前配置
  config init       重新执行初始化向导
  lang set zh       切换为中文日志
  lang set en       切换为英文日志

交互提示：
  TAB               自动补全命令和参数
  exit              退出交互模式

详细使用文档：
  docs/cloudcanal-cli-usage.md
`)
	}

	return strings.TrimSpace(`
CloudCanal CLI help

Help topics:
  help jobs         Show data job commands and filter meanings
  help datasources  Show datasource commands and filter meanings
  help clusters     Show cluster commands and filter meanings
  help workers      Show worker commands and filter meanings
  help consolejobs  Show console job commands
  help job-config   Show data job spec commands
  help config       Show configuration commands
  help lang         Show language switch commands

Common commands:
  jobs list         List data jobs
  datasources list  List data sources
  clusters list     List clusters
  workers list      List workers
  consolejobs show  Show console job details
  job-config specs  List data job specs
  config show       Show current config
  config init       Re-run the initialization wizard
  lang set zh       Switch CLI messages to Chinese
  lang set en       Switch CLI messages to English

REPL tips:
  TAB               Complete commands and options
  exit              Leave interactive mode

Detailed guide:
  docs/cloudcanal-cli-usage.md
`)
}

func (s *Shell) helpJobs() string {
	if s.isChinese() {
		return strings.TrimSpace(`
jobs 命令

jobs list [--name NAME] [--type TYPE] [--desc DESC] [--source-id ID] [--target-id ID]
  列出数据任务。
  --name       按任务名称过滤。
  --type       按任务类型过滤。
  --desc       按任务描述关键字过滤。
  --source-id  按源数据源实例 ID 过滤。
  --target-id  按目标数据源实例 ID 过滤。
  示例：cloudcanal jobs list --type SYNC --desc "nightly sync"

jobs show <jobId>
  查看单个任务详情。

jobs schema <jobId>
  查看任务的 schema 和映射配置。

jobs start <jobId>
  启动任务。

jobs stop <jobId>
  停止任务。

jobs delete <jobId>
  删除任务。

jobs replay <jobId> [--auto-start] [--reset-to-created]
  重放任务。
  --auto-start        重放后自动启动。
  --reset-to-created  重放前先重置到 CREATED 状态。
`)
	}

	return strings.TrimSpace(`
jobs commands

jobs list [--name NAME] [--type TYPE] [--desc DESC] [--source-id ID] [--target-id ID]
  List data jobs.
  --name       Filter by data job name.
  --type       Filter by data job type.
  --desc       Filter by description text.
  --source-id  Filter by source datasource instance ID.
  --target-id  Filter by target datasource instance ID.
  Example: cloudcanal jobs list --type SYNC --desc "nightly sync"

jobs show <jobId>
  Show a single job in detail.

jobs schema <jobId>
  Show schema and mapping config for a job.

jobs start <jobId>
  Start a job.

jobs stop <jobId>
  Stop a job.

jobs delete <jobId>
  Delete a job.

jobs replay <jobId> [--auto-start] [--reset-to-created]
  Replay a job.
  --auto-start        Start the job automatically after replay.
  --reset-to-created  Reset the job to CREATED before replay.
`)
}

func (s *Shell) helpDataSources() string {
	if s.isChinese() {
		return strings.TrimSpace(`
datasources 命令

datasources list [--id ID] [--type TYPE] [--deploy-type TYPE] [--host-type TYPE] [--lifecycle STATE]
  列出数据源。
  --id           按数据源 ID 过滤。
  --type         按数据源类型过滤，例如 MYSQL、STARROCKS。
  --deploy-type  按部署类型过滤，例如 ALIYUN、IDC。
  --host-type    按主机类型过滤，例如 RDS。
  --lifecycle    按生命周期状态过滤。
  示例：cloudcanal datasources list --type MYSQL --deploy-type ALIYUN

datasources show <dataSourceId>
  查看单个数据源详情。
`)
	}

	return strings.TrimSpace(`
datasources commands

datasources list [--id ID] [--type TYPE] [--deploy-type TYPE] [--host-type TYPE] [--lifecycle STATE]
  List data sources.
  --id           Filter by datasource ID.
  --type         Filter by datasource type, such as MYSQL or STARROCKS.
  --deploy-type  Filter by deploy type, such as ALIYUN or IDC.
  --host-type    Filter by host type, such as RDS.
  --lifecycle    Filter by lifecycle state.
  Example: cloudcanal datasources list --type MYSQL --deploy-type ALIYUN

datasources show <dataSourceId>
  Show a single datasource in detail.
`)
}

func (s *Shell) helpClusters() string {
	if s.isChinese() {
		return strings.TrimSpace(`
clusters 命令

clusters list [--name NAME] [--desc DESC] [--cloud CLOUD] [--region REGION]
  列出集群。
  --name    按集群名称模糊过滤。
  --desc    按集群描述模糊过滤。
  --cloud   按云厂商或 IDC 名称过滤。
  --region  按区域过滤。
  示例：cloudcanal clusters list --name prod --region cn-hangzhou
`)
	}

	return strings.TrimSpace(`
clusters commands

clusters list [--name NAME] [--desc DESC] [--cloud CLOUD] [--region REGION]
  List clusters.
  --name    Filter by cluster name.
  --desc    Filter by cluster description.
  --cloud   Filter by cloud vendor or IDC name.
  --region  Filter by region.
  Example: cloudcanal clusters list --name prod --region cn-hangzhou
`)
}

func (s *Shell) helpWorkers() string {
	if s.isChinese() {
		return strings.TrimSpace(`
workers 命令

workers list [--cluster-id ID] [--source-id ID] [--target-id ID]
  列出机器。
  --cluster-id  按集群 ID 过滤。
  --source-id   按源数据源实例 ID 过滤。
  --target-id   按目标数据源实例 ID 过滤。
  示例：cloudcanal workers list --cluster-id 2

workers start <workerId>
  启动机器。

workers stop <workerId>
  停止机器。
`)
	}

	return strings.TrimSpace(`
workers commands

workers list [--cluster-id ID] [--source-id ID] [--target-id ID]
  List workers.
  --cluster-id  Filter by cluster ID.
  --source-id   Filter by source datasource instance ID.
  --target-id   Filter by target datasource instance ID.
  Example: cloudcanal workers list --cluster-id 2

workers start <workerId>
  Start a worker.

workers stop <workerId>
  Stop a worker.
`)
}

func (s *Shell) helpConsoleJobs() string {
	if s.isChinese() {
		return strings.TrimSpace(`
consolejobs 命令

consolejobs show <consoleJobId>
  查看单个 ConsoleJob 详情，包括任务步骤、状态和资源信息。
`)
	}

	return strings.TrimSpace(`
consolejobs commands

consolejobs show <consoleJobId>
  Show a single console job, including task steps, state, and resource info.
`)
}

func (s *Shell) helpJobConfig() string {
	if s.isChinese() {
		return strings.TrimSpace(`
job-config 命令

job-config specs [--type TYPE] [--initial-sync=true|false] [--short-term-sync=true|false]
  列出数据任务规格。
  --type               按任务类型过滤。
  --initial-sync       是否要求初始同步。
  --short-term-sync    是否要求短期同步。
  示例：cloudcanal job-config specs --type SYNC --initial-sync=true
`)
	}

	return strings.TrimSpace(`
job-config commands

job-config specs [--type TYPE] [--initial-sync=true|false] [--short-term-sync=true|false]
  List data job specs.
  --type               Filter by data job type.
  --initial-sync       Whether initial sync is required.
  --short-term-sync    Whether short-term sync is required.
  Example: cloudcanal job-config specs --type SYNC --initial-sync=true
`)
}

func (s *Shell) helpConfig() string {
	if s.isChinese() {
		return strings.TrimSpace(`
config 命令

config show
  查看当前配置，包括 apiBaseUrl、accessKey 掩码和当前 language。

config init
  重新进入初始化向导，更新 API 地址、密钥和 language。
`)
	}

	return strings.TrimSpace(`
config commands

config show
  Show current config, including apiBaseUrl, masked accessKey, and current language.

config init
  Re-run the initialization wizard to update API URL, credentials, and language.
`)
}

func (s *Shell) helpLanguage() string {
	if s.isChinese() {
		return fmt.Sprintf(strings.TrimSpace(`
lang 命令

%s
  查看当前 CLI 文案语言。

lang set <en|zh>
  立即切换 CLI 文案语言并持久化到配置文件。
  en  表示英文
  zh  表示中文
`), i18n.T("lang.usage"))
	}

	return fmt.Sprintf(strings.TrimSpace(`
lang commands

%s
  Show the current CLI message language.

lang set <en|zh>
  Switch the CLI message language immediately and persist it to config.
  en  English
  zh  Chinese
`), i18n.T("lang.usage"))
}

func (s *Shell) helpCompletion() string {
	if s.isChinese() {
		return strings.TrimSpace(`
TAB 补全（高级）

completion <zsh|bash> [commandName]
  手动输出 shell TAB 补全脚本。
  zsh   生成 zsh 补全脚本
  bash  生成 bash 补全脚本
  commandName 可选，用于指定安装后的命令名。

说明：
  REPL 模式下如果终端支持行编辑，TAB 会自动补全命令和参数。
  安装脚本会默认安装 zsh 和 bash 补全文件，通常不需要手动执行这个命令。
`)
	}

	return strings.TrimSpace(`
TAB completion (advanced)

completion <zsh|bash> [commandName]
  Print a shell TAB completion script manually.
  zsh   Generate the zsh completion script.
  bash  Generate the bash completion script.
  commandName is optional and overrides the installed command name.

Notes:
  In REPL mode, TAB completes commands and options when the terminal supports line editing.
  The install script installs zsh and bash completion files by default, so you rarely need this command.
`)
}
