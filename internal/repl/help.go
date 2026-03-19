package repl

import (
	"cloudcanal-openapi-cli/internal/i18n"
	"fmt"
	"strings"
)

const detailedGuideURL = "https://github.com/Arlowen/cloudcanal-openapi-cli/blob/main/docs/cloudcanal-cli-usage.md"

func (s *Shell) printHelp(args []string) {
	s.io.Println(s.renderHelp(args))
}

func RenderHelp(args []string) string {
	return (&Shell{}).renderHelp(args)
}

func (s *Shell) renderHelp(args []string) string {
	topic := ""
	if len(args) > 0 {
		topic = strings.ToLower(args[0])
	}

	switch topic {
	case "", "overview":
		return s.helpOverview()
	case "jobs":
		return s.helpJobs()
	case "datasources":
		return s.helpDataSources()
	case "clusters":
		return s.helpClusters()
	case "workers":
		return s.helpWorkers()
	case "consolejobs":
		return s.helpConsoleJobs()
	case "job-config", "jobconfig":
		return s.helpJobConfig()
	case "schemas", "schema":
		return s.helpSchemas()
	case "config":
		return s.helpConfig()
	case "lang", "language":
		return s.helpLanguage()
	case "completion":
		return s.helpCompletion()
	default:
		return s.unknownHelpText(topic)
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
  help schemas      查看 Schema 查询命令说明
  help config       查看配置命令说明

常用命令：
  jobs list         列出数据任务
  jobs create       创建数据任务
  datasources list  列出数据源
  datasources add   创建数据源
  clusters list     列出集群
  workers list      列出机器
  consolejobs show  查看 ConsoleJob 详情
  job-config specs  查看任务规格
  schemas list-trans-objs-by-meta 查看映射对象
  config show       查看当前配置
  config init       重新执行初始化向导
  config lang show  查看当前语言
  config lang set zh 切换为中文日志
  config lang set en 切换为英文日志

交互提示：
  TAB               自动补全命令和参数
  Ctrl+C            退出交互模式
  exit              退出交互模式

详细使用文档：
  ` + detailedGuideURL + `
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
  help schemas      Show schema lookup commands
  help config       Show configuration commands

Common commands:
  jobs list         List data jobs
  jobs create       Create a data job
  datasources list  List data sources
  datasources add   Create a data source
  clusters list     List clusters
  workers list      List workers
  consolejobs show  Show console job details
  job-config specs  List data job specs
  schemas list-trans-objs-by-meta List transfer objects by metadata
  config show       Show current config
  config init       Re-run the initialization wizard
  config lang show  Show current language
  config lang set zh Switch CLI messages to Chinese
  config lang set en Switch CLI messages to English

REPL tips:
  TAB               Complete commands and options
  Ctrl+C            Exit interactive mode
  exit              Leave interactive mode

Detailed guide:
  ` + detailedGuideURL + `
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

jobs create --body-file FILE.json
  按 SDK 的 AddJobRequest 字段创建任务。
  也支持 --body '{"..."}' 直接传 JSON。

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

jobs attach-incre-task <jobId>
  绑定增量任务。

jobs detach-incre-task <jobId>
  解绑增量任务。

jobs update-incre-pos --body-file FILE.json
  按 SDK 的 UpdateIncrePosRequest 字段更新增量位点。
  也支持 --body '{"..."}' 直接传 JSON。
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

jobs create --body-file FILE.json
  Create a job with the SDK AddJobRequest JSON fields.
  --body '{"..."}' is also supported.

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

jobs attach-incre-task <jobId>
  Attach the incremental task.

jobs detach-incre-task <jobId>
  Detach the incremental task.

jobs update-incre-pos --body-file FILE.json
  Update incremental position with the SDK UpdateIncrePosRequest JSON fields.
  --body '{"..."}' is also supported.
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

datasources add --body-file FILE.json [--security-file FILE] [--secret-file FILE]
  创建数据源。请求体支持两种形式：
  1. 直接传 ApiDsAddData JSON
  2. 传 {"dataSourceAddData":{...},"securityFilePath":"...","secretFilePath":"..."} 包装 JSON

datasources delete <dataSourceId>
  删除数据源。
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

datasources add --body-file FILE.json [--security-file FILE] [--secret-file FILE]
  Create a data source. The body can be either the ApiDsAddData JSON itself
  or a wrapper object containing dataSourceAddData plus optional file paths.

datasources delete <dataSourceId>
  Delete a data source.
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

workers list --cluster-id ID [--source-id ID] [--target-id ID]
  列出机器。
  --cluster-id  必填，按集群 ID 过滤。
  --source-id   按源数据源实例 ID 过滤。
  --target-id   按目标数据源实例 ID 过滤。
  示例：cloudcanal workers list --cluster-id 2

workers start <workerId>
  启动机器。

workers stop <workerId>
  停止机器。

workers delete <workerId>
  删除机器。

workers modify-mem-oversold <workerId> --percent N
  修改内存超卖百分比。

workers update-alert <workerId> --phone=true|false --email=true|false --im=true|false --sms=true|false
  更新机器告警开关。
`)
	}

	return strings.TrimSpace(`
workers commands

workers list --cluster-id ID [--source-id ID] [--target-id ID]
  List workers.
  --cluster-id  Required. Filter by cluster ID.
  --source-id   Filter by source datasource instance ID.
  --target-id   Filter by target datasource instance ID.
  Example: cloudcanal workers list --cluster-id 2

workers start <workerId>
  Start a worker.

workers stop <workerId>
  Stop a worker.

workers delete <workerId>
  Delete a worker.

workers modify-mem-oversold <workerId> --percent N
  Update the memory oversold percentage.

workers update-alert <workerId> --phone=true|false --email=true|false --im=true|false --sms=true|false
  Update worker alert channels.
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

job-config specs --type TYPE [--initial-sync=true|false] [--short-term-sync=true|false]
  列出数据任务规格。
  --type               必填，按任务类型过滤。
  --initial-sync       是否要求初始同步。
  --short-term-sync    是否要求短期同步。
  示例：cloudcanal job-config specs --type SYNC --initial-sync=true

job-config transform-job-type --source-type TYPE --target-type TYPE
  根据源端和目标端类型转换任务类型。
`)
	}

	return strings.TrimSpace(`
job-config commands

job-config specs --type TYPE [--initial-sync=true|false] [--short-term-sync=true|false]
  List data job specs.
  --type               Required. Filter by data job type.
  --initial-sync       Whether initial sync is required.
  --short-term-sync    Whether short-term sync is required.
  Example: cloudcanal job-config specs --type SYNC --initial-sync=true

job-config transform-job-type --source-type TYPE --target-type TYPE
  Transform the job type based on source and target types.
`)
}

func (s *Shell) helpSchemas() string {
	if s.isChinese() {
		return strings.TrimSpace(`
schemas 命令

schemas list-trans-objs-by-meta [参数]
  按源端/目标端元信息查询传输对象。
  --src-db         源端库名
  --src-schema     源端 schema
  --src-trans-obj  源端对象名
  --dst-db         目标库名
  --dst-schema     目标 schema
  --dst-tran-obj   目标对象名
  示例：cloudcanal schemas list-trans-objs-by-meta --src-db demo --src-trans-obj orders
`)
	}

	return strings.TrimSpace(`
schemas commands

schemas list-trans-objs-by-meta [flags]
  List transfer objects by source and target metadata.
  --src-db         Source database
  --src-schema     Source schema
  --src-trans-obj  Source transfer object
  --dst-db         Destination database
  --dst-schema     Destination schema
  --dst-tran-obj   Destination transfer object
  Example: cloudcanal schemas list-trans-objs-by-meta --src-db demo --src-trans-obj orders
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

config lang show
  查看当前 CLI 文案语言。

config lang set <en|zh>
  立即切换 CLI 文案语言并持久化到配置文件。
`)
	}

	return strings.TrimSpace(`
config commands

config show
  Show current config, including apiBaseUrl, masked accessKey, and current language.

config init
  Re-run the initialization wizard to update API URL, credentials, and language.

config lang show
  Show the current CLI message language.

config lang set <en|zh>
  Switch the CLI message language immediately and persist it to config.
`)
}

func (s *Shell) helpLanguage() string {
	if s.isChinese() {
		return fmt.Sprintf(strings.TrimSpace(`
config lang 命令

%s
  查看当前 CLI 文案语言。

config lang set <en|zh>
  立即切换 CLI 文案语言并持久化到配置文件。
  en  表示英文
  zh  表示中文
`), i18n.T("lang.usage"))
	}

	return fmt.Sprintf(strings.TrimSpace(`
config lang commands

%s
  Show the current CLI message language.

config lang set <en|zh>
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
