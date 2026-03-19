package repl

import (
	"cloudcanal-openapi-cli/internal/cluster"
	"cloudcanal-openapi-cli/internal/datasource"
	"cloudcanal-openapi-cli/internal/jobconfig"
	ccschema "cloudcanal-openapi-cli/internal/schema"
	"cloudcanal-openapi-cli/internal/util"
	"cloudcanal-openapi-cli/internal/worker"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func (s *Shell) handleDataSources(tokens []string) error {
	if len(tokens) < 2 {
		s.io.Println(s.usageDataSources())
		return nil
	}

	switch strings.ToLower(tokens[1]) {
	case "list":
		options, err := parseDataSourceListOptions(tokens[2:])
		if err != nil {
			return err
		}
		return s.printDataSources(options)
	case "add":
		options, err := parseDataSourceAddOptions(tokens[2:])
		if err != nil {
			return err
		}
		result, err := s.runtime.DataSources().Add(options)
		if err != nil {
			return err
		}
		return s.printDataSourceAddResult(result)
	case "delete":
		if len(tokens) != 3 {
			s.io.Println(s.usageDataSourceAction("delete"))
			return nil
		}
		dataSourceID, err := parsePositiveInt64(tokens[2], "dataSourceId")
		if err != nil {
			return err
		}
		if err := s.runtime.DataSources().Delete(dataSourceID); err != nil {
			return err
		}
		return s.printActionResult("datasource.deleted", "datasource", "deleted", dataSourceID)
	case "show":
		if len(tokens) != 3 {
			s.io.Println(s.usageDataSourceShow())
			return nil
		}
		dataSourceID, err := parsePositiveInt64(tokens[2], "dataSourceId")
		if err != nil {
			return err
		}
		return s.printDataSource(dataSourceID)
	default:
		s.printUnknownSubcommand("datasources", tokens[1], dataSourceSubcommands, s.usageDataSources())
		return nil
	}
}

func (s *Shell) handleClusters(tokens []string) error {
	if len(tokens) < 2 {
		s.io.Println(s.usageClusters())
		return nil
	}
	if !strings.EqualFold(tokens[1], "list") {
		s.printUnknownSubcommand("clusters", tokens[1], clusterSubcommands, s.usageClusters())
		return nil
	}

	options, err := parseClusterListOptions(tokens[2:])
	if err != nil {
		return err
	}
	return s.printClusters(options)
}

func (s *Shell) handleWorkers(tokens []string) error {
	if len(tokens) < 2 {
		s.io.Println(s.usageWorkers())
		return nil
	}

	switch strings.ToLower(tokens[1]) {
	case "list":
		options, err := parseWorkerListOptions(tokens[2:])
		if err != nil {
			return err
		}
		return s.printWorkers(options)
	case "start", "stop", "delete":
		if len(tokens) != 3 {
			s.io.Println(s.usageWorkerAction(strings.ToLower(tokens[1])))
			return nil
		}
		workerID, err := parsePositiveInt64(tokens[2], "workerId")
		if err != nil {
			return err
		}
		if strings.EqualFold(tokens[1], "start") {
			if err := s.runtime.Workers().Start(workerID); err != nil {
				return err
			}
			return s.printActionResult("worker.started", "worker", "started", workerID)
		}
		if strings.EqualFold(tokens[1], "delete") {
			if err := s.runtime.Workers().Delete(workerID); err != nil {
				return err
			}
			return s.printActionResult("worker.deleted", "worker", "deleted", workerID)
		}
		if err := s.runtime.Workers().Stop(workerID); err != nil {
			return err
		}
		return s.printActionResult("worker.stopped", "worker", "stopped", workerID)
	case "modify-mem-oversold":
		if len(tokens) < 3 {
			s.io.Println(s.usageWorkerModifyMemOverSold())
			return nil
		}
		workerID, err := parsePositiveInt64(tokens[2], "workerId")
		if err != nil {
			return err
		}
		options, err := parseFlagArgs(tokens[3:])
		if err != nil {
			return err
		}
		percentValue, err := parseRequiredPositiveInt64Option(options, "memOverSoldPercent", "percent", "mem-over-sold-percent")
		if err != nil {
			return err
		}
		if err := ensureNoUnknownOptions(options); err != nil {
			return err
		}
		if err := s.runtime.Workers().ModifyMemOverSold(workerID, int(percentValue)); err != nil {
			return err
		}
		return s.printActionResult("worker.memOverSoldUpdated", "worker", "modify-mem-oversold", workerID)
	case "update-alert":
		if len(tokens) < 3 {
			s.io.Println(s.usageWorkerUpdateAlert())
			return nil
		}
		workerID, err := parsePositiveInt64(tokens[2], "workerId")
		if err != nil {
			return err
		}
		options, err := parseFlagArgs(tokens[3:])
		if err != nil {
			return err
		}
		phone, err := parseRequiredBoolOption(options, "phone", "phone")
		if err != nil {
			return err
		}
		email, err := parseRequiredBoolOption(options, "email", "email")
		if err != nil {
			return err
		}
		im, err := parseRequiredBoolOption(options, "im", "im")
		if err != nil {
			return err
		}
		sms, err := parseRequiredBoolOption(options, "sms", "sms")
		if err != nil {
			return err
		}
		if err := ensureNoUnknownOptions(options); err != nil {
			return err
		}
		if err := s.runtime.Workers().UpdateWorkerAlert(workerID, phone, email, im, sms); err != nil {
			return err
		}
		return s.printActionResult("worker.alertUpdated", "worker", "update-alert", workerID)
	default:
		s.printUnknownSubcommand("workers", tokens[1], workerSubcommands, s.usageWorkers())
		return nil
	}
}

func (s *Shell) handleConsoleJobs(tokens []string) error {
	if len(tokens) != 3 {
		s.io.Println(s.usageConsoleJobs())
		return nil
	}
	if !strings.EqualFold(tokens[1], "show") {
		s.printUnknownSubcommand("consolejobs", tokens[1], consoleJobSubcommands, s.usageConsoleJobs())
		return nil
	}

	consoleJobID, err := parsePositiveInt64(tokens[2], "consoleJobId")
	if err != nil {
		return err
	}
	return s.printConsoleJob(consoleJobID)
}

func (s *Shell) handleJobConfig(tokens []string) error {
	if len(tokens) < 2 {
		s.io.Println(s.usageJobConfig())
		return nil
	}
	if !strings.EqualFold(tokens[1], "specs") {
		if !strings.EqualFold(tokens[1], "transform-job-type") {
			s.printUnknownSubcommand("job-config", tokens[1], jobConfigSubcommands, s.usageJobConfig())
			return nil
		}
		options, err := parseTransformJobTypeOptions(tokens[2:])
		if err != nil {
			return err
		}
		return s.printTransformJobType(options)
	}

	options, err := parseListSpecsOptions(tokens[2:])
	if err != nil {
		return err
	}
	return s.printSpecs(options)
}

func (s *Shell) handleSchemas(tokens []string) error {
	if len(tokens) < 2 {
		s.io.Println(s.usageSchemas())
		return nil
	}
	if !strings.EqualFold(tokens[1], "list-trans-objs-by-meta") {
		s.printUnknownSubcommand("schemas", tokens[1], schemaSubcommands, s.usageSchemas())
		return nil
	}
	options, err := parseSchemaListOptions(tokens[2:])
	if err != nil {
		return err
	}
	return s.printTransferObjects(options)
}

func (s *Shell) printDataSources(options datasource.ListOptions) error {
	sources, err := s.runtime.DataSources().List(options)
	if err != nil {
		return err
	}
	if s.isJSONOutput() {
		return s.printJSON(sources)
	}

	headers := []string{s.label("id"), s.label("instance"), s.label("type"), s.label("host"), s.label("deploy"), s.label("lifecycleState"), s.label("description")}
	rows := make([][]string, 0, len(sources))
	for _, source := range sources {
		rows = append(rows, []string{
			strconv.FormatInt(source.ID, 10),
			orDash(source.InstanceID),
			orDash(source.DataSourceType),
			orDash(source.HostType),
			orDash(source.DeployType),
			orDash(source.LifeCycleState),
			orDash(util.MaskSensitiveText(source.InstanceDesc)),
		})
	}

	s.io.Println(util.FormatTable(headers, rows))
	s.io.Println(s.countLabel("datasources", len(sources)))
	return nil
}

func (s *Shell) printDataSource(dataSourceID int64) error {
	source, err := s.runtime.DataSources().Get(dataSourceID)
	if err != nil {
		return err
	}
	if s.isJSONOutput() {
		return s.printJSON(source)
	}

	s.io.Println(s.sectionTitle("datasource.details"))
	s.io.Println(s.line(s.label("id"), strconv.FormatInt(source.ID, 10)))
	s.io.Println(s.line(s.label("instanceId"), orDash(source.InstanceID)))
	s.io.Println(s.line(s.label("description"), orDash(util.MaskSensitiveText(source.InstanceDesc))))
	s.io.Println(s.line(s.label("type"), orDash(source.DataSourceType)))
	s.io.Println(s.line(s.label("host"), orDash(source.HostType)))
	s.io.Println(s.line(s.label("deploy"), orDash(source.DeployType)))
	s.io.Println(s.line(s.label("region"), orDash(source.Region)))
	s.io.Println(s.line(s.label("lifecycle"), orDash(source.LifeCycleState)))
	s.io.Println(s.line(s.label("account"), orDash(source.AccountName)))
	s.io.Println(s.line(s.label("securityType"), orDash(source.SecurityType)))
	s.io.Println(s.line(s.label("consoleJobId"), formatOptionalStringID(string(source.ConsoleJobID))))
	s.io.Println(s.line(s.label("consoleTaskState"), orDash(source.ConsoleTaskState)))
	return nil
}

func (s *Shell) printClusters(options cluster.ListOptions) error {
	clusters, err := s.runtime.Clusters().List(options)
	if err != nil {
		return err
	}
	if s.isJSONOutput() {
		return s.printJSON(clusters)
	}

	headers := []string{s.label("id"), s.label("name"), s.label("region"), s.label("cloud"), s.label("workers"), s.label("running"), s.label("abnormal"), s.label("owner")}
	rows := make([][]string, 0, len(clusters))
	for _, item := range clusters {
		rows = append(rows, []string{
			strconv.FormatInt(item.ID, 10),
			orDash(item.ClusterName),
			orDash(item.Region),
			orDash(item.CloudOrIDCName),
			strconv.Itoa(item.WorkerCount),
			strconv.Itoa(item.RunningCount),
			strconv.Itoa(item.AbnormalCount),
			orDash(item.OwnerName),
		})
	}

	s.io.Println(util.FormatTable(headers, rows))
	s.io.Println(s.countLabel("clusters", len(clusters)))
	return nil
}

func (s *Shell) printWorkers(options worker.ListOptions) error {
	workers, err := s.runtime.Workers().List(options)
	if err != nil {
		return err
	}
	if s.isJSONOutput() {
		return s.printJSON(workers)
	}

	headers := []string{s.label("id"), s.label("name"), s.label("state"), s.label("type"), s.label("cluster"), s.label("privateIP"), s.label("health"), s.label("load")}
	rows := make([][]string, 0, len(workers))
	for _, item := range workers {
		rows = append(rows, []string{
			strconv.FormatInt(item.ID, 10),
			orDash(item.WorkerName),
			orDash(item.WorkerState),
			orDash(item.WorkerType),
			formatOptionalInt64(item.ClusterID),
			orDash(item.PrivateIP),
			orDash(item.HealthLevel),
			formatFloat(item.WorkerLoad),
		})
	}

	s.io.Println(util.FormatTable(headers, rows))
	s.io.Println(s.countLabel("workers", len(workers)))
	return nil
}

func (s *Shell) printConsoleJob(consoleJobID int64) error {
	job, err := s.runtime.ConsoleJobs().Get(consoleJobID)
	if err != nil {
		return err
	}
	if s.isJSONOutput() {
		return s.printJSON(job)
	}

	s.io.Println(s.sectionTitle("consolejob.details"))
	s.io.Println(s.line(s.label("id"), strconv.FormatInt(job.ID, 10)))
	s.io.Println(s.line(s.label("label"), orDash(job.Label)))
	s.io.Println(s.line(s.label("state"), orDash(job.TaskState)))
	s.io.Println(s.line(s.label("jobToken"), orDash(job.JobToken)))
	s.io.Println(s.line(s.label("launcher"), orDash(job.Launcher)))
	s.io.Println(s.line(s.label("dataJobName"), orDash(job.DataJobName)))
	s.io.Println(s.line(s.label("dataJobDesc"), orDash(job.DataJobDesc)))
	s.io.Println(s.line(s.label("workerName"), orDash(job.WorkerName)))
	s.io.Println(s.line(s.label("workerDesc"), orDash(job.WorkerDesc)))
	s.io.Println(s.line(s.label("dataSourceInstance"), orDash(job.DsInstanceID)))
	s.io.Println(s.line(s.label("dataSourceDesc"), orDash(util.MaskSensitiveText(job.DatasourceDesc))))
	s.io.Println(s.line(s.label("resourceType"), orDash(job.ResourceType)))
	s.io.Println(s.line(s.label("resourceId"), formatOptionalInt64(job.ResourceID)))
	s.io.Println(s.line(s.label("tasks"), strconv.Itoa(len(job.TaskVOList))))

	if len(job.TaskVOList) > 0 {
		headers := []string{s.label("taskId"), s.label("state"), s.label("step"), s.label("host"), s.label("order"), s.label("cancelable")}
		rows := make([][]string, 0, len(job.TaskVOList))
		for _, task := range job.TaskVOList {
			rows = append(rows, []string{
				strconv.FormatInt(task.ID, 10),
				orDash(task.TaskState),
				orDash(task.StepName),
				orDash(task.Host),
				strconv.Itoa(task.ExecuteOrder),
				formatBool(task.Cancelable),
			})
		}
		s.io.Println("")
		s.io.Println(util.FormatTable(headers, rows))
	}
	return nil
}

func (s *Shell) printSpecs(options jobconfig.ListSpecsOptions) error {
	specs, err := s.runtime.JobConfigs().ListSpecs(options)
	if err != nil {
		return err
	}
	if s.isJSONOutput() {
		return s.printJSON(specs)
	}

	headers := []string{s.label("id"), s.label("jobType"), s.label("kind"), s.label("spec"), s.label("fullMB"), s.label("increMB"), s.label("checkMB")}
	rows := make([][]string, 0, len(specs))
	for _, spec := range specs {
		rows = append(rows, []string{
			strconv.FormatInt(spec.ID, 10),
			orDash(spec.SpecKind),
			orDash(spec.SpecKindCN),
			orDash(spec.Spec),
			strconv.Itoa(spec.FullMemoryMB),
			strconv.Itoa(spec.IncreMemoryMB),
			strconv.Itoa(spec.CheckMemoryMB),
		})
	}

	s.io.Println(util.FormatTable(headers, rows))
	s.io.Println(s.countLabel("specs", len(specs)))
	return nil
}

func (s *Shell) printTransformJobType(options jobconfig.TransformJobTypeOptions) error {
	result, err := s.runtime.JobConfigs().TransformJobType(options)
	if err != nil {
		return err
	}
	if s.isJSONOutput() {
		if len(result.Data) == 0 {
			return s.printJSON(map[string]any{})
		}
		var payload any
		if err := json.Unmarshal(result.Data, &payload); err != nil {
			return s.printJSON(map[string]any{"data": string(result.Data)})
		}
		return s.printJSON(payload)
	}
	title := "Transform job type result:"
	if s.isChinese() {
		title = "任务类型转换结果："
	}
	s.io.Println(title)
	if len(result.Data) == 0 {
		s.io.Println("{}")
		return nil
	}
	var payload any
	if err := json.Unmarshal(result.Data, &payload); err != nil {
		s.io.Println(string(result.Data))
		return nil
	}
	return s.writeJSON(payload)
}

func (s *Shell) printTransferObjects(options ccschema.ListTransObjsByMetaOptions) error {
	items, err := s.runtime.Schemas().ListTransObjsByMeta(options)
	if err != nil {
		return err
	}
	if s.isJSONOutput() {
		return s.printJSON(items)
	}

	headers := []string{s.label("jobId"), s.label("jobName"), s.label("source"), s.label("target"), s.label("srcType"), s.label("dstType")}
	rows := make([][]string, 0, len(items))
	for _, item := range items {
		rows = append(rows, []string{
			formatOptionalInt64(item.DataJobID),
			orDash(item.DataJobName),
			orDash(item.SrcFullTransferObjName),
			orDash(item.DstFullTransferObjName),
			orDash(item.SrcDsType),
			orDash(item.DstDsType),
		})
	}

	s.io.Println(util.FormatTable(headers, rows))
	s.io.Println(s.countLabel("schemas", len(items)))
	return nil
}

func parseDataSourceAddOptions(tokens []string) (datasource.AddOptions, error) {
	options, err := parseFlagArgs(tokens)
	if err != nil {
		return datasource.AddOptions{}, err
	}

	rawBody, err := readBodyOptions(options)
	if err != nil {
		return datasource.AddOptions{}, err
	}

	var result datasource.AddOptions
	var probe map[string]json.RawMessage
	if err := json.Unmarshal(rawBody, &probe); err == nil {
		if _, ok := probe["dataSourceAddData"]; ok {
			if err := json.Unmarshal(rawBody, &result); err != nil {
				return datasource.AddOptions{}, fmt.Errorf("failed to parse request body: %w", err)
			}
		} else {
			if err := json.Unmarshal(rawBody, &result.DataSourceAddData); err != nil {
				return datasource.AddOptions{}, fmt.Errorf("failed to parse request body: %w", err)
			}
		}
	} else if err := json.Unmarshal(rawBody, &result.DataSourceAddData); err != nil {
		return datasource.AddOptions{}, fmt.Errorf("failed to parse request body: %w", err)
	}

	if value, ok := popOption(options, "security-file"); ok && strings.TrimSpace(value) != "" {
		result.SecurityFilePath = value
	}
	if value, ok := popOption(options, "secret-file"); ok && strings.TrimSpace(value) != "" {
		result.SecretFilePath = value
	}
	if err := ensureNoUnknownOptions(options); err != nil {
		return datasource.AddOptions{}, err
	}
	return result, nil
}

func parseDataSourceListOptions(tokens []string) (datasource.ListOptions, error) {
	options, err := parseFlagArgs(tokens)
	if err != nil {
		return datasource.ListOptions{}, err
	}

	result := datasource.ListOptions{}
	result.DataSourceID, err = parsePositiveInt64Option(options, "dataSourceId", "id", "data-source-id")
	if err != nil {
		return datasource.ListOptions{}, err
	}
	result.Type, _ = popOption(options, "type")
	result.DeployType, _ = popOption(options, "deploy-type")
	result.HostType, _ = popOption(options, "host-type")
	result.LifeCycleState, _ = popOption(options, "lifecycle", "life-cycle-state")
	if err := ensureNoUnknownOptions(options); err != nil {
		return datasource.ListOptions{}, err
	}
	return result, nil
}

func parseClusterListOptions(tokens []string) (cluster.ListOptions, error) {
	options, err := parseFlagArgs(tokens)
	if err != nil {
		return cluster.ListOptions{}, err
	}

	result := cluster.ListOptions{}
	result.ClusterName, _ = popOption(options, "name", "cluster-name")
	result.ClusterDesc, _ = popOption(options, "desc", "cluster-desc")
	result.CloudOrIDCName, _ = popOption(options, "cloud", "cloud-or-idc", "cloud-or-idc-name")
	result.Region, _ = popOption(options, "region")
	if err := ensureNoUnknownOptions(options); err != nil {
		return cluster.ListOptions{}, err
	}
	return result, nil
}

func parseWorkerListOptions(tokens []string) (worker.ListOptions, error) {
	options, err := parseFlagArgs(tokens)
	if err != nil {
		return worker.ListOptions{}, err
	}

	result := worker.ListOptions{}
	result.ClusterID, err = parseRequiredPositiveInt64Option(options, "clusterId", "cluster-id")
	if err != nil {
		return worker.ListOptions{}, err
	}
	result.SourceInstanceID, err = parsePositiveInt64Option(options, "sourceInstanceId", "source-id", "source-instance-id")
	if err != nil {
		return worker.ListOptions{}, err
	}
	result.TargetInstanceID, err = parsePositiveInt64Option(options, "targetInstanceId", "target-id", "target-instance-id")
	if err != nil {
		return worker.ListOptions{}, err
	}
	if err := ensureNoUnknownOptions(options); err != nil {
		return worker.ListOptions{}, err
	}
	return result, nil
}

func parseListSpecsOptions(tokens []string) (jobconfig.ListSpecsOptions, error) {
	options, err := parseFlagArgs(tokens)
	if err != nil {
		return jobconfig.ListSpecsOptions{}, err
	}

	result := jobconfig.ListSpecsOptions{}
	result.DataJobType, err = parseRequiredStringOption(options, "dataJobType", "type", "data-job-type")
	if err != nil {
		return jobconfig.ListSpecsOptions{}, err
	}
	result.InitialSync, err = parseBoolOption(options, "initialSync", "initial-sync")
	if err != nil {
		return jobconfig.ListSpecsOptions{}, err
	}
	result.ShortTermSync, err = parseBoolOption(options, "shortTermSync", "short-term-sync")
	if err != nil {
		return jobconfig.ListSpecsOptions{}, err
	}
	if err := ensureNoUnknownOptions(options); err != nil {
		return jobconfig.ListSpecsOptions{}, err
	}
	return result, nil
}

func parseTransformJobTypeOptions(tokens []string) (jobconfig.TransformJobTypeOptions, error) {
	options, err := parseFlagArgs(tokens)
	if err != nil {
		return jobconfig.TransformJobTypeOptions{}, err
	}

	result := jobconfig.TransformJobTypeOptions{}
	result.SourceType, err = parseRequiredStringOption(options, "sourceType", "source-type")
	if err != nil {
		return jobconfig.TransformJobTypeOptions{}, err
	}
	result.TargetType, err = parseRequiredStringOption(options, "targetType", "target-type")
	if err != nil {
		return jobconfig.TransformJobTypeOptions{}, err
	}
	if err := ensureNoUnknownOptions(options); err != nil {
		return jobconfig.TransformJobTypeOptions{}, err
	}
	return result, nil
}

func parseSchemaListOptions(tokens []string) (ccschema.ListTransObjsByMetaOptions, error) {
	options, err := parseFlagArgs(tokens)
	if err != nil {
		return ccschema.ListTransObjsByMetaOptions{}, err
	}

	result := ccschema.ListTransObjsByMetaOptions{}
	result.SrcDb, _ = popOption(options, "src-db")
	result.SrcSchema, _ = popOption(options, "src-schema")
	result.SrcTransObj, _ = popOption(options, "src-trans-obj")
	result.DstDb, _ = popOption(options, "dst-db")
	result.DstSchema, _ = popOption(options, "dst-schema")
	result.DstTranObj, _ = popOption(options, "dst-tran-obj")
	if err := ensureNoUnknownOptions(options); err != nil {
		return ccschema.ListTransObjsByMetaOptions{}, err
	}
	return result, nil
}

func formatFloat(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

func formatOptionalStringID(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || trimmed == "0" {
		return "-"
	}
	return trimmed
}

func (s *Shell) printDataSourceAddResult(data string) error {
	message := "Data source created successfully"
	if s.isChinese() {
		message = "数据源创建成功"
	}
	if s.isJSONOutput() {
		return s.printJSON(map[string]any{
			"resource": "datasource",
			"action":   "created",
			"data":     data,
			"message":  message,
		})
	}
	s.io.Println(message)
	if strings.TrimSpace(data) != "" {
		s.io.Println(s.line(s.label("result"), data))
	}
	return nil
}
