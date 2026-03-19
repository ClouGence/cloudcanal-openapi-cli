# cloudcanal-cli 使用说明

`cloudcanal-cli` 是 CloudCanal OpenAPI 的命令行工具，既支持交互式使用，也支持单次命令执行。

## 启动方式

交互模式：

```bash
cloudcanal
```

单次命令模式：

```bash
cloudcanal jobs list
cloudcanal datasources list --type MYSQL
cloudcanal jobs list --type SYNC --output json
```

如果还没有安装到系统命令，也可以直接执行本地二进制：

```bash
./bin/cloudcanal jobs list
```

一键安装：

```bash
curl -fsSL https://raw.githubusercontent.com/Arlowen/cloudcanal-openapi-cli/main/scripts/bootstrap_install.sh | bash
```

说明：

- 当前一键安装会从 GitHub Releases 下载预编译二进制
- 下载后会自动校验 release 里的 `checksums.txt`
- 不需要本机安装 `Go`
- 默认会把二进制安装到 `~/.local/share/cloudcanal-openapi-cli/bin/cloudcanal`
- 之后会自动完成命令、PATH 和补全安装

一键卸载：

```bash
curl -fsSL https://raw.githubusercontent.com/Arlowen/cloudcanal-openapi-cli/main/scripts/bootstrap_uninstall.sh | bash
```

## 初始化配置

首次启动会进入初始化向导，配置文件默认保存到：

```text
~/.cloudcanal/config.json
```

配置格式：

```json
{
  "apiBaseUrl": "https://cc.example.com",
  "accessKey": "your-ak",
  "secretKey": "your-sk",
  "language": "en",
  "httpTimeoutSeconds": 15,
  "httpReadMaxRetries": 2,
  "httpReadRetryBackoffMillis": 300
}
```

说明：

- `apiBaseUrl` 必须包含 `http://` 或 `https://`
- `accessKey` 是访问密钥 ID
- `secretKey` 是访问密钥 Secret，不会在 `config show` 中明文展示
- `language` 是 CLI 文案语言，支持 `en` 和 `zh`
- `httpTimeoutSeconds` 是单次 HTTP 请求超时秒数，默认 `10`
- `httpReadMaxRetries` 是只读请求的最大重试次数，默认 `0`
- `httpReadRetryBackoffMillis` 是只读请求的首次退避毫秒数，默认 `250`

## 基本命令

`help`

显示帮助入口。也支持：

- `help jobs`
- `help datasources`
- `help clusters`
- `help workers`
- `help consolejobs`
- `help job-config`
- `help config`
- `help lang`

`clear` / `cls`

清空当前终端屏幕，适合在交互模式下快速整理输出内容。

`config show`

显示当前配置，`accessKey` 会做掩码处理，同时会显示当前 `language`。

`config init`

重新执行初始化向导，更新配置。

`lang show`

显示当前 CLI 文案语言。

`lang set <en|zh>`

切换 CLI 文案语言，并持久化到配置文件。

`completion <zsh|bash> [commandName]`

输出 shell TAB 补全脚本。安装脚本默认会安装 zsh 和 bash 的补全文件；如果你想手动安装，也可以直接执行：

```bash
cloudcanal completion zsh
cloudcanal completion bash
```

`exit` / `quit`

退出交互模式。

## Jobs

`jobs list [参数]`

列出数据任务，支持以下参数：

- `--name <name>`: 按任务名称过滤
- `--type <type>`: 按任务类型过滤
- `--desc <desc>`: 按任务描述过滤
- `--source-id <id>`: 按源数据源 ID 过滤
- `--target-id <id>`: 按目标数据源 ID 过滤
- `--output <text|json>`: 输出文本表格或 JSON

示例：

```bash
cloudcanal jobs list --type SYNC --name demo
cloudcanal jobs list --desc "nightly sync"
cloudcanal jobs list --type SYNC --output json
```

`jobs show <jobId>`

查看任务详情。

`jobs schema <jobId>`

查看任务 schema 信息。

`jobs start <jobId>`

启动任务。

`jobs stop <jobId>`

停止任务。

`jobs delete <jobId>`

删除任务。

`jobs replay <jobId> [--auto-start] [--reset-to-created]`

重放任务。

- `--auto-start`: 重放后自动启动
- `--reset-to-created`: 重放前重置到创建状态

示例：

```bash
cloudcanal jobs replay 123 --auto-start --reset-to-created
```

## DataSource

`datasources list [参数]`

列出数据源，支持以下参数：

- `--id <id>`: 按数据源 ID 过滤
- `--type <type>`: 按数据源类型过滤
- `--deploy-type <type>`: 按部署类型过滤
- `--host-type <type>`: 按主机类型过滤
- `--lifecycle <state>`: 按生命周期状态过滤

示例：

```bash
cloudcanal datasources list --type MYSQL --deploy-type ALIYUN
```

`datasources show <dataSourceId>`

查看单个数据源详情。

## Cluster

`clusters list [参数]`

列出集群，支持以下参数：

- `--name <name>`: 按集群名模糊过滤
- `--desc <desc>`: 按集群描述模糊过滤
- `--cloud <name>`: 按云厂商或 IDC 名称过滤
- `--region <region>`: 按区域过滤

示例：

```bash
cloudcanal clusters list --name prod --region cn-hangzhou
```

## Worker

`workers list [参数]`

列出机器，支持以下参数：

- `--cluster-id <id>`: 按集群 ID 过滤
- `--source-id <id>`: 按源数据源 ID 过滤
- `--target-id <id>`: 按目标数据源 ID 过滤

示例：

```bash
cloudcanal workers list --cluster-id 2
```

`workers start <workerId>`

启动机器。

`workers stop <workerId>`

停止机器。

## ConsoleJob

`consolejobs show <consoleJobId>`

查看控制台异步任务详情。

## DataJob 配置

`job-config specs [参数]`

列出数据任务配置规格，支持以下参数：

- `--type <type>`: 按数据任务类型过滤
- `--initial-sync=<true|false>`: 是否初始同步
- `--short-term-sync=<true|false>`: 是否短期同步

示例：

```bash
cloudcanal job-config specs --type SYNC --initial-sync=true
```

## 使用建议

- 带空格的参数值请使用引号包裹，例如 `--desc "nightly sync"`
- 可以在查询类命令后追加 `--output json` 获取机器可读结果
- 交互模式下如果终端支持行编辑，可以直接使用 `TAB` 补全命令、子命令和常见参数
- 可以先执行 `cloudcanal help` 查看帮助主题，再执行 `cloudcanal help jobs` 这类子帮助查看参数含义
- 如果想切换中文或英文文案，可执行 `cloudcanal lang set zh` 或 `cloudcanal lang set en`
- 如果命令执行失败，优先检查 `apiBaseUrl`、`accessKey`、`secretKey` 是否正确
