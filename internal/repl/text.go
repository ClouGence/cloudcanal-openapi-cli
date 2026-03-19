package repl

import (
	"cloudcanal-openapi-cli/internal/i18n"
	"cloudcanal-openapi-cli/internal/util"
	"fmt"
	"strconv"
	"strings"
)

const detailLabelWidth = 24

func (s *Shell) isChinese() bool {
	return i18n.CurrentLanguage() == i18n.Chinese
}

func usageBlock(title string, commands ...string) string {
	lines := make([]string, 0, len(commands)+1)
	lines = append(lines, title)
	for _, command := range commands {
		lines = append(lines, "  "+command)
	}
	return strings.Join(lines, "\n")
}

func (s *Shell) usageConfig() string {
	if s.isChinese() {
		return usageBlock("用法：", "config show", "config init", "config lang show", "config lang set <en|zh>")
	}
	return usageBlock("Usage:", "config show", "config init", "config lang show", "config lang set <en|zh>")
}

func (s *Shell) usageJobsGroup() string {
	if s.isChinese() {
		return usageBlock(
			"用法：",
			"jobs list",
			"jobs create --body-file FILE.json",
			"jobs show <jobId>",
			"jobs schema <jobId>",
			"jobs start <jobId>",
			"jobs stop <jobId>",
			"jobs delete <jobId>",
			"jobs replay <jobId> [--auto-start] [--reset-to-created]",
			"jobs attach-incre-task <jobId>",
			"jobs detach-incre-task <jobId>",
			"jobs update-incre-pos --body-file FILE.json",
		)
	}
	return usageBlock(
		"Usage:",
		"jobs list",
		"jobs create --body-file FILE.json",
		"jobs show <jobId>",
		"jobs schema <jobId>",
		"jobs start <jobId>",
		"jobs stop <jobId>",
		"jobs delete <jobId>",
		"jobs replay <jobId> [--auto-start] [--reset-to-created]",
		"jobs attach-incre-task <jobId>",
		"jobs detach-incre-task <jobId>",
		"jobs update-incre-pos --body-file FILE.json",
	)
}

func (s *Shell) usageJobsList() string {
	if s.isChinese() {
		return "用法：jobs list [--name NAME] [--type TYPE] [--desc DESC] [--source-id ID] [--target-id ID]"
	}
	return "Usage: jobs list [--name NAME] [--type TYPE] [--desc DESC] [--source-id ID] [--target-id ID]"
}

func (s *Shell) usageJobCreate() string {
	if s.isChinese() {
		return "用法：jobs create --body-file FILE.json | --body '{...}'"
	}
	return "Usage: jobs create --body-file FILE.json | --body '{...}'"
}

func (s *Shell) usageJobAction(action string) string {
	if s.isChinese() {
		return fmt.Sprintf("用法：jobs %s <jobId>", action)
	}
	return fmt.Sprintf("Usage: jobs %s <jobId>", action)
}

func (s *Shell) usageJobReplay() string {
	if s.isChinese() {
		return "用法：jobs replay <jobId> [--auto-start] [--reset-to-created]"
	}
	return "Usage: jobs replay <jobId> [--auto-start] [--reset-to-created]"
}

func (s *Shell) usageJobUpdateIncrePos() string {
	if s.isChinese() {
		return "用法：jobs update-incre-pos --body-file FILE.json | --body '{...}'"
	}
	return "Usage: jobs update-incre-pos --body-file FILE.json | --body '{...}'"
}

func (s *Shell) usageDataSources() string {
	if s.isChinese() {
		return usageBlock(
			"用法：",
			"datasources list [--id ID] [--type TYPE] [--deploy-type TYPE] [--host-type TYPE] [--lifecycle STATE]",
			"datasources add --body-file FILE.json [--security-file FILE] [--secret-file FILE]",
			"datasources delete <dataSourceId>",
			"datasources show <dataSourceId>",
		)
	}
	return usageBlock(
		"Usage:",
		"datasources list [--id ID] [--type TYPE] [--deploy-type TYPE] [--host-type TYPE] [--lifecycle STATE]",
		"datasources add --body-file FILE.json [--security-file FILE] [--secret-file FILE]",
		"datasources delete <dataSourceId>",
		"datasources show <dataSourceId>",
	)
}

func (s *Shell) usageDataSourcesList() string {
	if s.isChinese() {
		return "用法：datasources list [--id ID] [--type TYPE] [--deploy-type TYPE] [--host-type TYPE] [--lifecycle STATE]"
	}
	return "Usage: datasources list [--id ID] [--type TYPE] [--deploy-type TYPE] [--host-type TYPE] [--lifecycle STATE]"
}

func (s *Shell) usageDataSourceShow() string {
	if s.isChinese() {
		return "用法：datasources show <dataSourceId>"
	}
	return "Usage: datasources show <dataSourceId>"
}

func (s *Shell) usageDataSourceAdd() string {
	if s.isChinese() {
		return "用法：datasources add --body-file FILE.json [--security-file FILE] [--secret-file FILE]"
	}
	return "Usage: datasources add --body-file FILE.json [--security-file FILE] [--secret-file FILE]"
}

func (s *Shell) usageDataSourceAction(action string) string {
	if s.isChinese() {
		return fmt.Sprintf("用法：datasources %s <dataSourceId>", action)
	}
	return fmt.Sprintf("Usage: datasources %s <dataSourceId>", action)
}

func (s *Shell) usageClusters() string {
	if s.isChinese() {
		return "用法：clusters list [--name NAME] [--desc DESC] [--cloud CLOUD] [--region REGION]"
	}
	return "Usage: clusters list [--name NAME] [--desc DESC] [--cloud CLOUD] [--region REGION]"
}

func (s *Shell) usageClustersList() string {
	if s.isChinese() {
		return "用法：clusters list [--name NAME] [--desc DESC] [--cloud CLOUD] [--region REGION]"
	}
	return "Usage: clusters list [--name NAME] [--desc DESC] [--cloud CLOUD] [--region REGION]"
}

func (s *Shell) usageWorkers() string {
	if s.isChinese() {
		return usageBlock(
			"用法：",
			"workers list --cluster-id ID [--source-id ID] [--target-id ID]",
			"workers start <workerId>",
			"workers stop <workerId>",
			"workers delete <workerId>",
			"workers modify-mem-oversold <workerId> --percent N",
			"workers update-alert <workerId> --phone=true|false --email=true|false --im=true|false --sms=true|false",
		)
	}
	return usageBlock(
		"Usage:",
		"workers list --cluster-id ID [--source-id ID] [--target-id ID]",
		"workers start <workerId>",
		"workers stop <workerId>",
		"workers delete <workerId>",
		"workers modify-mem-oversold <workerId> --percent N",
		"workers update-alert <workerId> --phone=true|false --email=true|false --im=true|false --sms=true|false",
	)
}

func (s *Shell) usageWorkersList() string {
	if s.isChinese() {
		return "用法：workers list --cluster-id ID [--source-id ID] [--target-id ID]"
	}
	return "Usage: workers list --cluster-id ID [--source-id ID] [--target-id ID]"
}

func (s *Shell) usageWorkerAction(action string) string {
	if s.isChinese() {
		return fmt.Sprintf("用法：workers %s <workerId>", action)
	}
	return fmt.Sprintf("Usage: workers %s <workerId>", action)
}

func (s *Shell) usageWorkerModifyMemOverSold() string {
	if s.isChinese() {
		return "用法：workers modify-mem-oversold <workerId> --percent N"
	}
	return "Usage: workers modify-mem-oversold <workerId> --percent N"
}

func (s *Shell) usageWorkerUpdateAlert() string {
	if s.isChinese() {
		return "用法：workers update-alert <workerId> --phone=true|false --email=true|false --im=true|false --sms=true|false"
	}
	return "Usage: workers update-alert <workerId> --phone=true|false --email=true|false --im=true|false --sms=true|false"
}

func (s *Shell) usageConsoleJobs() string {
	if s.isChinese() {
		return "用法：consolejobs show <consoleJobId>"
	}
	return "Usage: consolejobs show <consoleJobId>"
}

func (s *Shell) usageConsoleJobShow() string {
	if s.isChinese() {
		return "用法：consolejobs show <consoleJobId>"
	}
	return "Usage: consolejobs show <consoleJobId>"
}

func (s *Shell) usageJobConfig() string {
	if s.isChinese() {
		return usageBlock(
			"用法：",
			"job-config specs --type TYPE [--initial-sync=true|false] [--short-term-sync=true|false]",
			"job-config transform-job-type --source-type TYPE --target-type TYPE",
		)
	}
	return usageBlock(
		"Usage:",
		"job-config specs --type TYPE [--initial-sync=true|false] [--short-term-sync=true|false]",
		"job-config transform-job-type --source-type TYPE --target-type TYPE",
	)
}

func (s *Shell) usageJobConfigSpecs() string {
	if s.isChinese() {
		return "用法：job-config specs --type TYPE [--initial-sync=true|false] [--short-term-sync=true|false]"
	}
	return "Usage: job-config specs --type TYPE [--initial-sync=true|false] [--short-term-sync=true|false]"
}

func (s *Shell) usageJobConfigTransform() string {
	if s.isChinese() {
		return "用法：job-config transform-job-type --source-type TYPE --target-type TYPE"
	}
	return "Usage: job-config transform-job-type --source-type TYPE --target-type TYPE"
}

func (s *Shell) usageSchemas() string {
	if s.isChinese() {
		return "用法：schemas list-trans-objs-by-meta [--src-db NAME] [--src-schema NAME] [--src-trans-obj NAME] [--dst-db NAME] [--dst-schema NAME] [--dst-tran-obj NAME]"
	}
	return "Usage: schemas list-trans-objs-by-meta [--src-db NAME] [--src-schema NAME] [--src-trans-obj NAME] [--dst-db NAME] [--dst-schema NAME] [--dst-tran-obj NAME]"
}

func (s *Shell) usageConfigShow() string {
	if s.isChinese() {
		return "用法：config show"
	}
	return "Usage: config show"
}

func (s *Shell) usageConfigInit() string {
	if s.isChinese() {
		return "用法：config init"
	}
	return "Usage: config init"
}

func (s *Shell) usageConfigLang() string {
	if s.isChinese() {
		return usageBlock("用法：", "config lang show", "config lang set <en|zh>")
	}
	return usageBlock("Usage:", "config lang show", "config lang set <en|zh>")
}

func (s *Shell) usageCompletion() string {
	if s.isChinese() {
		return "用法：completion <zsh|bash> [commandName]"
	}
	return "Usage: completion <zsh|bash> [commandName]"
}

func (s *Shell) actionMessage(kind string, id int64) string {
	if s.isChinese() {
		switch kind {
		case "job.started":
			return fmt.Sprintf("任务 %d 已启动", id)
		case "job.stopped":
			return fmt.Sprintf("任务 %d 已停止", id)
		case "job.deleted":
			return fmt.Sprintf("任务 %d 已删除", id)
		case "job.replayed":
			return fmt.Sprintf("任务 %d 已提交重放请求", id)
		case "job.increAttached":
			return fmt.Sprintf("任务 %d 已绑定增量任务", id)
		case "job.increDetached":
			return fmt.Sprintf("任务 %d 已解绑增量任务", id)
		case "datasource.deleted":
			return fmt.Sprintf("数据源 %d 已删除", id)
		case "worker.started":
			return fmt.Sprintf("机器 %d 已启动", id)
		case "worker.stopped":
			return fmt.Sprintf("机器 %d 已停止", id)
		case "worker.deleted":
			return fmt.Sprintf("机器 %d 已删除", id)
		case "worker.memOverSoldUpdated":
			return fmt.Sprintf("机器 %d 的内存超卖比例已更新", id)
		case "worker.alertUpdated":
			return fmt.Sprintf("机器 %d 的告警配置已更新", id)
		}
	}
	switch kind {
	case "job.started":
		return fmt.Sprintf("Job %d started successfully", id)
	case "job.stopped":
		return fmt.Sprintf("Job %d stopped successfully", id)
	case "job.deleted":
		return fmt.Sprintf("Job %d deleted successfully", id)
	case "job.replayed":
		return fmt.Sprintf("Job %d replay requested successfully", id)
	case "job.increAttached":
		return fmt.Sprintf("Job %d incremental task attached successfully", id)
	case "job.increDetached":
		return fmt.Sprintf("Job %d incremental task detached successfully", id)
	case "datasource.deleted":
		return fmt.Sprintf("Data source %d deleted successfully", id)
	case "worker.started":
		return fmt.Sprintf("Worker %d started successfully", id)
	case "worker.stopped":
		return fmt.Sprintf("Worker %d stopped successfully", id)
	case "worker.deleted":
		return fmt.Sprintf("Worker %d deleted successfully", id)
	case "worker.memOverSoldUpdated":
		return fmt.Sprintf("Worker %d memory oversold percentage updated successfully", id)
	case "worker.alertUpdated":
		return fmt.Sprintf("Worker %d alert config updated successfully", id)
	default:
		return ""
	}
}

func (s *Shell) sectionTitle(key string) string {
	if s.isChinese() {
		switch key {
		case "job.details":
			return "任务详情："
		case "job.schema":
			return "任务 Schema："
		case "job.mappingConfig":
			return "映射配置："
		case "datasource.details":
			return "数据源详情："
		case "consolejob.details":
			return "ConsoleJob 详情："
		}
	}
	switch key {
	case "job.details":
		return "Job details:"
	case "job.schema":
		return "Job schema:"
	case "job.mappingConfig":
		return "Mapping Config:"
	case "datasource.details":
		return "Data source details:"
	case "consolejob.details":
		return "Console job details:"
	default:
		return ""
	}
}

func (s *Shell) countLabel(kind string, count int) string {
	if s.isChinese() {
		switch kind {
		case "jobs":
			return fmt.Sprintf("%d 个任务", count)
		case "datasources":
			return fmt.Sprintf("%d 个数据源", count)
		case "clusters":
			return fmt.Sprintf("%d 个集群", count)
		case "workers":
			return fmt.Sprintf("%d 台机器", count)
		case "specs":
			return fmt.Sprintf("%d 条规格", count)
		case "schemas":
			return fmt.Sprintf("%d 个对象", count)
		}
	}
	switch kind {
	case "jobs":
		return fmt.Sprintf("%d jobs", count)
	case "datasources":
		return fmt.Sprintf("%d data sources", count)
	case "clusters":
		return fmt.Sprintf("%d clusters", count)
	case "workers":
		return fmt.Sprintf("%d workers", count)
	case "specs":
		return fmt.Sprintf("%d specs", count)
	case "schemas":
		return fmt.Sprintf("%d transfer objects", count)
	default:
		return strconv.Itoa(count)
	}
}

func (s *Shell) label(key string) string {
	if s.isChinese() {
		switch key {
		case "id":
			return "ID"
		case "jobId":
			return "任务 ID"
		case "jobName":
			return "任务名"
		case "name":
			return "名称"
		case "type":
			return "类型"
		case "state":
			return "状态"
		case "source":
			return "源端"
		case "target":
			return "目标端"
		case "description":
			return "描述"
		case "currentTaskStatus":
			return "当前任务状态"
		case "lifecycle":
			return "生命周期"
		case "user":
			return "用户"
		case "consoleJobId":
			return "Console Job ID"
		case "consoleTaskState":
			return "Console 任务状态"
		case "sourceSchema":
			return "源 Schema"
		case "targetSchema":
			return "目标 Schema"
		case "tasks":
			return "任务数"
		case "hasException":
			return "是否有异常"
		case "taskId":
			return "任务 ID"
		case "workerIP":
			return "Worker IP"
		case "defaultTopic":
			return "默认 Topic"
		case "defaultTopicPartition":
			return "默认 Topic 分区"
		case "schemaWhitelistLevel":
			return "Schema 白名单级别"
		case "srcSchemaLessFormat":
			return "源端无 Schema 格式"
		case "dstSchemaLessFormat":
			return "目标端无 Schema 格式"
		case "instance":
			return "实例"
		case "instanceId":
			return "实例 ID"
		case "host":
			return "主机类型"
		case "deploy":
			return "部署类型"
		case "region":
			return "区域"
		case "account":
			return "账号"
		case "securityType":
			return "安全类型"
		case "label":
			return "标签"
		case "jobToken":
			return "任务 Token"
		case "launcher":
			return "发起人"
		case "dataJobName":
			return "任务名称"
		case "dataJobDesc":
			return "任务描述"
		case "workerName":
			return "机器名称"
		case "workerDesc":
			return "机器描述"
		case "dataSourceInstance":
			return "数据源实例"
		case "dataSourceDesc":
			return "数据源描述"
		case "resourceType":
			return "资源类型"
		case "resourceId":
			return "资源 ID"
		case "step":
			return "步骤"
		case "order":
			return "顺序"
		case "cancelable":
			return "可取消"
		case "kind":
			return "类别"
		case "spec":
			return "规格"
		case "fullMB":
			return "全量 MB"
		case "increMB":
			return "增量 MB"
		case "checkMB":
			return "校验 MB"
		case "workers":
			return "机器数"
		case "running":
			return "运行中"
		case "abnormal":
			return "异常"
		case "cloud":
			return "云厂商"
		case "cluster":
			return "集群"
		case "owner":
			return "负责人"
		case "health":
			return "健康度"
		case "load":
			return "负载"
		case "privateIP":
			return "私网 IP"
		case "jobType":
			return "任务类型"
		case "lifecycleState":
			return "生命周期"
		case "srcType":
			return "源类型"
		case "dstType":
			return "目标类型"
		case "result":
			return "结果"
		}
	}

	switch key {
	case "id":
		return "ID"
	case "jobId":
		return "Job ID"
	case "jobName":
		return "Job Name"
	case "name":
		return "Name"
	case "type":
		return "Type"
	case "state":
		return "State"
	case "source":
		return "Source"
	case "target":
		return "Target"
	case "description":
		return "Description"
	case "currentTaskStatus":
		return "Current Task Status"
	case "lifecycle":
		return "Lifecycle"
	case "user":
		return "User"
	case "consoleJobId":
		return "Console Job ID"
	case "consoleTaskState":
		return "Console Task State"
	case "sourceSchema":
		return "Source Schema"
	case "targetSchema":
		return "Target Schema"
	case "tasks":
		return "Tasks"
	case "hasException":
		return "Has Exception"
	case "taskId":
		return "Task ID"
	case "workerIP":
		return "Worker IP"
	case "defaultTopic":
		return "Default Topic"
	case "defaultTopicPartition":
		return "Default Topic Partition"
	case "schemaWhitelistLevel":
		return "Schema Whitelist Level"
	case "srcSchemaLessFormat":
		return "Source Schema Less Format"
	case "dstSchemaLessFormat":
		return "Target Schema Less Format"
	case "instance":
		return "Instance"
	case "instanceId":
		return "Instance ID"
	case "host":
		return "Host"
	case "deploy":
		return "Deploy"
	case "region":
		return "Region"
	case "account":
		return "Account"
	case "securityType":
		return "Security Type"
	case "label":
		return "Label"
	case "jobToken":
		return "Job Token"
	case "launcher":
		return "Launcher"
	case "dataJobName":
		return "Data Job Name"
	case "dataJobDesc":
		return "Data Job Desc"
	case "workerName":
		return "Worker Name"
	case "workerDesc":
		return "Worker Desc"
	case "dataSourceInstance":
		return "Data Source Instance"
	case "dataSourceDesc":
		return "Data Source Desc"
	case "resourceType":
		return "Resource Type"
	case "resourceId":
		return "Resource ID"
	case "step":
		return "Step"
	case "order":
		return "Order"
	case "cancelable":
		return "Cancelable"
	case "kind":
		return "Kind"
	case "spec":
		return "Spec"
	case "fullMB":
		return "Full MB"
	case "increMB":
		return "Incre MB"
	case "checkMB":
		return "Check MB"
	case "workers":
		return "Workers"
	case "running":
		return "Running"
	case "abnormal":
		return "Abnormal"
	case "cloud":
		return "Cloud"
	case "cluster":
		return "Cluster"
	case "owner":
		return "Owner"
	case "health":
		return "Health"
	case "load":
		return "Load"
	case "privateIP":
		return "Private IP"
	case "jobType":
		return "Job Type"
	case "lifecycleState":
		return "Lifecycle"
	case "srcType":
		return "Source Type"
	case "dstType":
		return "Target Type"
	case "result":
		return "Result"
	default:
		return key
	}
}

func (s *Shell) line(label string, value string) string {
	return "  " + util.PadDisplayRight(label, detailLabelWidth) + " : " + strings.TrimSpace(value)
}
