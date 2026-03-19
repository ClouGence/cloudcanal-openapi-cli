# cloudcanal-openapi-cli

CloudCanal OpenAPI 的命令行工具，支持：

- 交互式命令行
- 单次命令执行
- `--output json` 机器可读输出
- 安装脚本默认配置 zsh / bash TAB 补全

完整命令说明见 [docs/cloudcanal-cli-usage.md](docs/cloudcanal-cli-usage.md)。

## 快速开始

1. 安装

```bash
curl -fsSL https://raw.githubusercontent.com/Arlowen/cloudcanal-openapi-cli/main/scripts/install.sh | bash
```

2. 启动并完成初始化

```bash
cloudcanal
```

## 常用用法

交互模式：

```bash
cloudcanal
```

单次命令：

```bash
cloudcanal --help
cloudcanal jobs --help
cloudcanal config lang set zh
cloudcanal jobs list
cloudcanal jobs show 123
cloudcanal jobs create --body-file create-job.json
cloudcanal datasources list --type MYSQL
cloudcanal datasources add --body-file add-datasource.json
cloudcanal workers list --cluster-id 2
cloudcanal schemas list-trans-objs-by-meta --src-db demo --src-trans-obj orders
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
- SDK API 对照 CLI 命令：[docs/openapi-sdk-api-reference.md](docs/openapi-sdk-api-reference.md)
- 版本变更记录：[CHANGELOG.md](CHANGELOG.md)
- 机器可读输出：在查询命令后追加 `--output json`
- 交互模式和 shell 均支持 TAB 自动补全

## 卸载

```bash
curl -fsSL https://raw.githubusercontent.com/Arlowen/cloudcanal-openapi-cli/main/scripts/uninstall.sh | bash
```
