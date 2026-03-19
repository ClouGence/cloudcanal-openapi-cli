package i18n

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

const (
	English = "en"
	Chinese = "zh"
)

var (
	languageMu      sync.RWMutex
	currentLanguage = DefaultLanguage()
)

var messages = map[string]map[string]string{
	English: {
		"common.typeHelp":                          "Type 'help' to see available commands.",
		"common.useHelp":                           "Use 'help' to see available commands.",
		"common.errorPrefix":                       "Error: %s",
		"common.fatalErrorPrefix":                  "Fatal error: %s",
		"common.unknownCommand":                    "Unknown command: %s",
		"common.unknownSubcommand":                 "Unknown %s command: %s",
		"common.unknownHelpTopic":                  "Unknown help topic: %s",
		"common.didYouMean":                        "Did you mean: %s",
		"common.supportedLanguages":                "Supported languages: en, zh",
		"runtime.invalidConfig":                    "Existing configuration is invalid: %s",
		"runtime.initCancelled":                    "Initialization cancelled.",
		"runtime.configUpdated":                    "Configuration updated.",
		"config.apiBaseUrlLabel":                   "apiBaseUrl",
		"config.accessKeyLabel":                    "accessKey",
		"config.languageLabel":                     "language",
		"config.apiBaseUrlRequired":                "apiBaseUrl is required",
		"config.accessKeyRequired":                 "accessKey is required",
		"config.secretKeyRequired":                 "secretKey is required",
		"config.languageUnsupported":               "language must be en or zh",
		"config.httpTimeoutLabel":                  "httpTimeoutSeconds",
		"config.httpReadMaxRetriesLabel":           "httpReadMaxRetries",
		"config.httpReadRetryBackoffMillisLabel":   "httpReadRetryBackoffMillis",
		"config.httpTimeoutSecondsInvalid":         "httpTimeoutSeconds must be zero or a positive integer",
		"config.httpReadMaxRetriesInvalid":         "httpReadMaxRetries must be zero or a positive integer",
		"config.httpReadRetryBackoffMillisInvalid": "httpReadRetryBackoffMillis must be zero or a positive integer",
		"config.apiBaseUrlInvalid":                 "apiBaseUrl is not a valid URL",
		"config.apiBaseUrlScheme":                  "apiBaseUrl must start with http:// or https://",
		"config.apiBaseUrlHost":                    "apiBaseUrl must contain a host",
		"config.invalidJSON":                       "configuration file is not valid JSON",
		"wizard.title":                             "CloudCanal CLI initialization",
		"wizard.cancelHint":                        "Type exit at any prompt to cancel.",
		"wizard.apiHostHint":                       "apiHost must be a full URL, for example: https://cc.example.com",
		"wizard.languageHint":                      "Language supports en or zh. Press Enter to use the default language.",
		"wizard.keepCurrent":                       "Press Enter to keep the current value.",
		"wizard.invalidConfig":                     "Invalid configuration: %s",
		"wizard.invalidField":                      "Invalid %s: %s",
		"wizard.checkingConnection":                "Checking OpenAPI connection...",
		"wizard.validationFailed":                  "Configuration validation failed: %s",
		"wizard.reuseValues":                       "Press Enter to reuse the current values, or type new ones.",
		"wizard.savedTo":                           "Configuration saved to %s",
		"parser.unterminatedEscape":                "unterminated escape sequence",
		"parser.unterminatedQuote":                 "unterminated quote",
		"parser.unexpectedArgument":                "unexpected argument: %s",
		"parser.invalidOption":                     "invalid option: %s",
		"parser.duplicateOption":                   "duplicate option: --%s",
		"parser.unknownOption":                     "unknown option: --%s",
		"parser.outputOptionRequiresValue":         "output requires a value: text or json",
		"parser.outputOptionInvalid":               "output must be text or json",
		"parser.optionRequired":                    "%s is required",
		"parser.mustBePositiveInt":                 "%s must be a positive integer",
		"parser.mustBeBoolean":                     "%s must be a boolean",
		"parser.jobId":                             "jobId",
		"parser.workerId":                          "workerId",
		"parser.consoleJobId":                      "consoleJobId",
		"parser.dataSourceId":                      "dataSourceId",
		"parser.dataJobType":                       "dataJobType",
		"parser.sourceInstanceId":                  "sourceInstanceId",
		"parser.targetInstanceId":                  "targetInstanceId",
		"parser.clusterId":                         "clusterId",
		"parser.initialSync":                       "initialSync",
		"parser.shortTermSync":                     "shortTermSync",
		"parser.autoStart":                         "autoStart",
		"parser.resetToCreated":                    "resetToCreated",
		"lang.usage":                               "Usage: config lang show | config lang set <en|zh>",
		"lang.current":                             "Current language: %s",
		"lang.updated":                             "Language switched to %s.",
		"lang.en":                                  "English",
		"lang.zh":                                  "Chinese",
		"util.serverError":                         "server error",
		"util.unknownError":                        "unknown error",
	},
	Chinese: {
		"common.typeHelp":                          "输入 'help' 查看可用命令。",
		"common.useHelp":                           "输入 'help' 查看可用命令。",
		"common.errorPrefix":                       "错误：%s",
		"common.fatalErrorPrefix":                  "致命错误：%s",
		"common.unknownCommand":                    "未知命令：%s",
		"common.unknownSubcommand":                 "未知 %s 命令：%s",
		"common.unknownHelpTopic":                  "未知帮助主题：%s",
		"common.didYouMean":                        "你是不是想输入：%s",
		"common.supportedLanguages":                "支持的语言：en、zh",
		"runtime.invalidConfig":                    "现有配置无效：%s",
		"runtime.initCancelled":                    "已取消初始化。",
		"runtime.configUpdated":                    "配置已更新。",
		"config.apiBaseUrlLabel":                   "apiBaseUrl",
		"config.accessKeyLabel":                    "accessKey",
		"config.languageLabel":                     "language",
		"config.apiBaseUrlRequired":                "apiBaseUrl 不能为空",
		"config.accessKeyRequired":                 "accessKey 不能为空",
		"config.secretKeyRequired":                 "secretKey 不能为空",
		"config.languageUnsupported":               "language 只支持 en 或 zh",
		"config.httpTimeoutLabel":                  "httpTimeoutSeconds",
		"config.httpReadMaxRetriesLabel":           "httpReadMaxRetries",
		"config.httpReadRetryBackoffMillisLabel":   "httpReadRetryBackoffMillis",
		"config.httpTimeoutSecondsInvalid":         "httpTimeoutSeconds 必须是 0 或正整数",
		"config.httpReadMaxRetriesInvalid":         "httpReadMaxRetries 必须是 0 或正整数",
		"config.httpReadRetryBackoffMillisInvalid": "httpReadRetryBackoffMillis 必须是 0 或正整数",
		"config.apiBaseUrlInvalid":                 "apiBaseUrl 不是合法的 URL",
		"config.apiBaseUrlScheme":                  "apiBaseUrl 必须以 http:// 或 https:// 开头",
		"config.apiBaseUrlHost":                    "apiBaseUrl 必须包含主机名",
		"config.invalidJSON":                       "配置文件不是合法的 JSON",
		"wizard.title":                             "CloudCanal CLI 初始化",
		"wizard.cancelHint":                        "任意提示下输入 exit 可取消。",
		"wizard.apiHostHint":                       "apiHost 必须是完整 URL，例如：https://cc.example.com",
		"wizard.languageHint":                      "language 支持 en 或 zh，直接回车使用默认语言。",
		"wizard.keepCurrent":                       "直接回车可保留当前值。",
		"wizard.invalidConfig":                     "配置不合法：%s",
		"wizard.invalidField":                      "%s 输入不合法：%s",
		"wizard.checkingConnection":                "正在检查 OpenAPI 连接...",
		"wizard.validationFailed":                  "配置校验失败：%s",
		"wizard.reuseValues":                       "直接回车可复用当前值，也可以重新输入。",
		"wizard.savedTo":                           "配置已保存到 %s",
		"parser.unterminatedEscape":                "转义符没有正常结束",
		"parser.unterminatedQuote":                 "引号没有正常结束",
		"parser.unexpectedArgument":                "存在未识别的参数：%s",
		"parser.invalidOption":                     "参数格式不合法：%s",
		"parser.duplicateOption":                   "参数重复：--%s",
		"parser.unknownOption":                     "未知参数：--%s",
		"parser.outputOptionRequiresValue":         "output 需要指定 text 或 json",
		"parser.outputOptionInvalid":               "output 只支持 text 或 json",
		"parser.optionRequired":                    "%s 不能为空",
		"parser.mustBePositiveInt":                 "%s 必须是正整数",
		"parser.mustBeBoolean":                     "%s 必须是布尔值",
		"parser.jobId":                             "jobId",
		"parser.workerId":                          "workerId",
		"parser.consoleJobId":                      "consoleJobId",
		"parser.dataSourceId":                      "dataSourceId",
		"parser.dataJobType":                       "dataJobType",
		"parser.sourceInstanceId":                  "sourceInstanceId",
		"parser.targetInstanceId":                  "targetInstanceId",
		"parser.clusterId":                         "clusterId",
		"parser.initialSync":                       "initialSync",
		"parser.shortTermSync":                     "shortTermSync",
		"parser.autoStart":                         "autoStart",
		"parser.resetToCreated":                    "resetToCreated",
		"lang.usage":                               "用法：config lang show | config lang set <en|zh>",
		"lang.current":                             "当前语言：%s",
		"lang.updated":                             "语言已切换为 %s。",
		"lang.en":                                  "英文",
		"lang.zh":                                  "中文",
		"util.serverError":                         "服务端错误",
		"util.unknownError":                        "未知错误",
	},
}

func DefaultLanguage() string {
	return English
}

func NormalizeLanguage(language string) string {
	switch strings.ToLower(strings.TrimSpace(language)) {
	case "", English, "en-us", "en_us":
		return English
	case Chinese, "zh-cn", "zh_cn", "zh-hans", "cn":
		return Chinese
	default:
		return ""
	}
}

func CurrentLanguage() string {
	languageMu.RLock()
	defer languageMu.RUnlock()
	return currentLanguage
}

func SetLanguage(language string) error {
	normalized := NormalizeLanguage(language)
	if normalized == "" {
		return errors.New(T("config.languageUnsupported"))
	}
	languageMu.Lock()
	currentLanguage = normalized
	languageMu.Unlock()
	return nil
}

func T(key string, args ...any) string {
	return TFor(CurrentLanguage(), key, args...)
}

func TFor(language string, key string, args ...any) string {
	normalized := NormalizeLanguage(language)
	if normalized == "" {
		normalized = DefaultLanguage()
	}

	template := lookup(normalized, key)
	if len(args) == 0 {
		return template
	}
	return fmt.Sprintf(template, args...)
}

func DisplayName(language string) string {
	normalized := NormalizeLanguage(language)
	switch normalized {
	case Chinese:
		return TFor(normalized, "lang.zh")
	default:
		return TFor(normalized, "lang.en")
	}
}

func lookup(language string, key string) string {
	if value, ok := messages[language][key]; ok {
		return value
	}
	if value, ok := messages[DefaultLanguage()][key]; ok {
		return value
	}
	return key
}
