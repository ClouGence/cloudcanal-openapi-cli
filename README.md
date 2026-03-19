# cloudcanal-openapi-cli

基于 Go 实现的 CloudCanal OpenAPI CLI，支持交互式使用，也支持单次命令执行。

详细使用说明见 [docs/cloudcanal-cli-usage.md](docs/cloudcanal-cli-usage.md)。

## 快速开始

要求：

- 日常使用：`curl`、`tar`
- 本地源码开发：Go 1.25+

构建并测试：

```bash
./scripts/all_build.sh
```

源码方式安装到命令行：

```bash
./scripts/install.sh
```

安装脚本会同时安装 zsh / bash 的 TAB 补全文件。

一键安装：

```bash
curl -fsSL https://raw.githubusercontent.com/Arlowen/cloudcanal-openapi-cli/main/scripts/bootstrap_install.sh | bash
```

这个一键安装脚本会从 GitHub Releases 下载预编译二进制，不需要本机安装 Go。
默认会把二进制安装到 `~/.local/share/cloudcanal-openapi-cli/bin/cloudcanal`。

一键卸载：

```bash
curl -fsSL https://raw.githubusercontent.com/Arlowen/cloudcanal-openapi-cli/main/scripts/bootstrap_uninstall.sh | bash
```

卸载：

```bash
./scripts/uninstall.sh
```

## 使用方式

详细命令、参数和示例请看上面的使用说明文档。

## 初始化配置

第一次启动会进入初始化向导，配置文件保存到 `~/.cloudcanal/config.json`。配置格式、字段含义和命令参数说明见详细文档。

## 开发

发布：

- 推送 tag，例如 `v0.1.0`
- GitHub Actions 会自动构建并发布 `darwin/linux + amd64/arm64` 的 release 资产

只编译：

```bash
make build
```

只测试：

```bash
make test
```
