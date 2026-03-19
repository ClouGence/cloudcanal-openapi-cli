# cloudcanal-openapi-cli

基于 Go 实现的 CloudCanal OpenAPI CLI，支持交互式使用，也支持单次命令执行。

详细使用说明见 [docs/cloudcanal-cli-usage.md](docs/cloudcanal-cli-usage.md)。

## 快速开始

要求：

- Go 1.25+

构建并测试：

```bash
./scripts/all_build.sh
```

安装到命令行：

```bash
./scripts/install.sh
```

安装脚本会同时安装 zsh / bash 的 TAB 补全文件。

一键安装：

```bash
curl -fsSL https://raw.githubusercontent.com/Arlowen/cloudcanal-openapi-cli/main/scripts/bootstrap_install.sh | bash
```

这个一键安装脚本会下载源码归档并在本机编译，所以仍然需要本机已有 Go 1.25+。
安装后的仓库默认会落在 `~/.local/share/cloudcanal-openapi-cli/repository`。

卸载：

```bash
./scripts/uninstall.sh
```

## 使用方式

详细命令、参数和示例请看上面的使用说明文档。

## 初始化配置

第一次启动会进入初始化向导，配置文件保存到 `~/.cloudcanal/config.json`。配置格式、字段含义和命令参数说明见详细文档。

## 开发

只编译：

```bash
make build
```

只测试：

```bash
make test
```
