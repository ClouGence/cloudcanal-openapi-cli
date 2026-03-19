package repl

import (
	"cloudcanal-openapi-cli/internal/cluster"
	"cloudcanal-openapi-cli/internal/datasource"
	"cloudcanal-openapi-cli/internal/jobconfig"
	"cloudcanal-openapi-cli/internal/util"
	"cloudcanal-openapi-cli/internal/worker"
	"fmt"
	"strconv"
	"strings"
)

func (s *Shell) handleDataSources(tokens []string) error {
	if len(tokens) < 2 {
		s.io.Println("Usage: datasources list [--id ID] [--type TYPE] [--deploy-type TYPE] [--host-type TYPE] [--lifecycle STATE] | datasources show <dataSourceId>")
		return nil
	}

	switch strings.ToLower(tokens[1]) {
	case "list":
		options, err := parseDataSourceListOptions(tokens[2:])
		if err != nil {
			return err
		}
		return s.printDataSources(options)
	case "show":
		if len(tokens) != 3 {
			s.io.Println("Usage: datasources show <dataSourceId>")
			return nil
		}
		dataSourceID, err := parsePositiveInt64(tokens[2], "dataSourceId")
		if err != nil {
			return err
		}
		return s.printDataSource(dataSourceID)
	default:
		s.io.Println("Usage: datasources list [--id ID] [--type TYPE] [--deploy-type TYPE] [--host-type TYPE] [--lifecycle STATE] | datasources show <dataSourceId>")
		return nil
	}
}

func (s *Shell) handleClusters(tokens []string) error {
	if len(tokens) < 2 || !strings.EqualFold(tokens[1], "list") {
		s.io.Println("Usage: clusters list [--name NAME] [--desc DESC] [--cloud CLOUD] [--region REGION]")
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
		s.io.Println("Usage: workers list [--cluster-id ID] [--source-id ID] [--target-id ID] | workers start <workerId> | workers stop <workerId>")
		return nil
	}

	switch strings.ToLower(tokens[1]) {
	case "list":
		options, err := parseWorkerListOptions(tokens[2:])
		if err != nil {
			return err
		}
		return s.printWorkers(options)
	case "start", "stop":
		if len(tokens) != 3 {
			s.io.Println("Usage: workers " + strings.ToLower(tokens[1]) + " <workerId>")
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
			s.io.Println(fmt.Sprintf("Worker %d started successfully", workerID))
			return nil
		}
		if err := s.runtime.Workers().Stop(workerID); err != nil {
			return err
		}
		s.io.Println(fmt.Sprintf("Worker %d stopped successfully", workerID))
		return nil
	default:
		s.io.Println("Usage: workers list [--cluster-id ID] [--source-id ID] [--target-id ID] | workers start <workerId> | workers stop <workerId>")
		return nil
	}
}

func (s *Shell) handleConsoleJobs(tokens []string) error {
	if len(tokens) != 3 || !strings.EqualFold(tokens[1], "show") {
		s.io.Println("Usage: consolejobs show <consoleJobId>")
		return nil
	}

	consoleJobID, err := parsePositiveInt64(tokens[2], "consoleJobId")
	if err != nil {
		return err
	}
	return s.printConsoleJob(consoleJobID)
}

func (s *Shell) handleJobConfig(tokens []string) error {
	if len(tokens) < 2 || !strings.EqualFold(tokens[1], "specs") {
		s.io.Println("Usage: job-config specs [--type TYPE] [--initial-sync=true|false] [--short-term-sync=true|false]")
		return nil
	}

	options, err := parseListSpecsOptions(tokens[2:])
	if err != nil {
		return err
	}
	return s.printSpecs(options)
}

func (s *Shell) printDataSources(options datasource.ListOptions) error {
	sources, err := s.runtime.DataSources().List(options)
	if err != nil {
		return err
	}

	headers := []string{"ID", "Instance", "Type", "Host", "Deploy", "Lifecycle", "Description"}
	rows := make([][]string, 0, len(sources))
	for _, source := range sources {
		rows = append(rows, []string{
			strconv.FormatInt(source.ID, 10),
			orDash(source.InstanceID),
			orDash(source.DataSourceType),
			orDash(source.HostType),
			orDash(source.DeployType),
			orDash(source.LifeCycleState),
			orDash(source.InstanceDesc),
		})
	}

	s.io.Println(util.FormatTable(headers, rows))
	s.io.Println(fmt.Sprintf("%d data sources", len(sources)))
	return nil
}

func (s *Shell) printDataSource(dataSourceID int64) error {
	source, err := s.runtime.DataSources().Get(dataSourceID)
	if err != nil {
		return err
	}

	s.io.Println("Data source details:")
	s.io.Println("  ID: " + strconv.FormatInt(source.ID, 10))
	s.io.Println("  Instance ID: " + orDash(source.InstanceID))
	s.io.Println("  Description: " + orDash(source.InstanceDesc))
	s.io.Println("  Type: " + orDash(source.DataSourceType))
	s.io.Println("  Host Type: " + orDash(source.HostType))
	s.io.Println("  Deploy Type: " + orDash(source.DeployType))
	s.io.Println("  Region: " + orDash(source.Region))
	s.io.Println("  Lifecycle: " + orDash(source.LifeCycleState))
	s.io.Println("  Account: " + orDash(source.AccountName))
	s.io.Println("  Security Type: " + orDash(source.SecurityType))
	s.io.Println("  Console Job ID: " + orDash(source.ConsoleJobID))
	s.io.Println("  Console Task State: " + orDash(source.ConsoleTaskState))
	return nil
}

func (s *Shell) printClusters(options cluster.ListOptions) error {
	clusters, err := s.runtime.Clusters().List(options)
	if err != nil {
		return err
	}

	headers := []string{"ID", "Name", "Region", "Cloud", "Workers", "Running", "Abnormal", "Owner"}
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
	s.io.Println(fmt.Sprintf("%d clusters", len(clusters)))
	return nil
}

func (s *Shell) printWorkers(options worker.ListOptions) error {
	workers, err := s.runtime.Workers().List(options)
	if err != nil {
		return err
	}

	headers := []string{"ID", "Name", "State", "Type", "Cluster", "Private IP", "Health", "Load"}
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
	s.io.Println(fmt.Sprintf("%d workers", len(workers)))
	return nil
}

func (s *Shell) printConsoleJob(consoleJobID int64) error {
	job, err := s.runtime.ConsoleJobs().Get(consoleJobID)
	if err != nil {
		return err
	}

	s.io.Println("Console job details:")
	s.io.Println("  ID: " + strconv.FormatInt(job.ID, 10))
	s.io.Println("  Label: " + orDash(job.Label))
	s.io.Println("  State: " + orDash(job.TaskState))
	s.io.Println("  Job Token: " + orDash(job.JobToken))
	s.io.Println("  Launcher: " + orDash(job.Launcher))
	s.io.Println("  Data Job Name: " + orDash(job.DataJobName))
	s.io.Println("  Data Job Desc: " + orDash(job.DataJobDesc))
	s.io.Println("  Worker Name: " + orDash(job.WorkerName))
	s.io.Println("  Worker Desc: " + orDash(job.WorkerDesc))
	s.io.Println("  Data Source Instance: " + orDash(job.DsInstanceID))
	s.io.Println("  Data Source Desc: " + orDash(job.DatasourceDesc))
	s.io.Println("  Resource Type: " + orDash(job.ResourceType))
	s.io.Println("  Resource ID: " + formatOptionalInt64(job.ResourceID))
	s.io.Println("  Tasks: " + strconv.Itoa(len(job.TaskVOList)))

	if len(job.TaskVOList) > 0 {
		headers := []string{"Task ID", "State", "Step", "Host", "Order", "Cancelable"}
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

	headers := []string{"ID", "Job Type", "Kind", "Spec", "Full MB", "Incre MB", "Check MB"}
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
	s.io.Println(fmt.Sprintf("%d specs", len(specs)))
	return nil
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
	result.ClusterID, err = parsePositiveInt64Option(options, "clusterId", "cluster-id")
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
	result.DataJobType, _ = popOption(options, "type", "data-job-type")
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

func formatFloat(value float64) string {
	return fmt.Sprintf("%.2f", value)
}
