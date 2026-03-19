# cloudcanal-openapi-cli

CloudCanal OpenAPI 的命令行工具，支持：

- 交互式命令行
- 单次命令执行
- `--output json` 机器可读输出
- zsh / bash TAB 补全

完整命令说明见 [docs/cloudcanal-cli-usage.md](docs/cloudcanal-cli-usage.md)。

## 快速开始

1. 安装

```bash
curl -fsSL https://raw.githubusercontent.com/Arlowen/cloudcanal-openapi-cli/main/scripts/bootstrap_install.sh | bash
```

2. 启动并完成初始化

```bash
cloudcanal
```

首次启动会进入初始化向导。

默认安装目录是 `~/.cloudcanal-cli`，其中二进制位于 `~/.cloudcanal-cli/bin/cloudcanal`，补全文件位于 `~/.cloudcanal-cli/completions`。

## 常用用法

交互模式：

```bash
cloudcanal
```

单次命令：

```bash
cloudcanal jobs list
cloudcanal jobs show 123
cloudcanal datasources list --type MYSQL
cloudcanal workers list --cluster-id 2
```

JSON 输出：

```bash
cloudcanal jobs list --type SYNC --output json
```

## 配置

配置文件默认保存在：

```text
~/.cloudcanal-cli/config.json
```

最小配置示例：

```json
{
  "apiBaseUrl": "https://cc.example.com",
  "accessKey": "your-ak",
  "secretKey": "your-sk",
  "language": "en"
}
```

如果你需要调整网络行为，也可以追加这些可选项：

```json
{
  "httpTimeoutSeconds": 15,
  "httpReadMaxRetries": 2,
  "httpReadRetryBackoffMillis": 300
}
```

## 文档入口

- 安装、初始化、命令参数、示例：[docs/cloudcanal-cli-usage.md](docs/cloudcanal-cli-usage.md)
- 机器可读输出：在查询命令后追加 `--output json`
- 补全脚本：`cloudcanal completion zsh` / `cloudcanal completion bash`

## 卸载

```bash
curl -fsSL https://raw.githubusercontent.com/Arlowen/cloudcanal-openapi-cli/main/scripts/bootstrap_uninstall.sh | bash
```
