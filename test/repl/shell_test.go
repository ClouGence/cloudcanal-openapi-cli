package repl_test

import (
	"cloudcanal-openapi-cli/internal/app"
	"cloudcanal-openapi-cli/internal/cluster"
	"cloudcanal-openapi-cli/internal/config"
	"cloudcanal-openapi-cli/internal/console"
	"cloudcanal-openapi-cli/internal/consolejob"
	"cloudcanal-openapi-cli/internal/datajob"
	"cloudcanal-openapi-cli/internal/datasource"
	"cloudcanal-openapi-cli/internal/i18n"
	"cloudcanal-openapi-cli/internal/jobconfig"
	"cloudcanal-openapi-cli/internal/repl"
	ccschema "cloudcanal-openapi-cli/internal/schema"
	"cloudcanal-openapi-cli/internal/worker"
	"cloudcanal-openapi-cli/test/testsupport"
	"encoding/json"
	"strings"
	"testing"

	"github.com/peterh/liner"
)

func TestShellHandlesHappyPathCommands(t *testing.T) {
	dataJobs := &fakeDataJobs{
		jobs: []datajob.Job{
			{
				DataJobID:     11,
				DataJobName:   "sync-job",
				DataJobType:   "SYNC",
				DataTaskState: "RUNNING",
				SourceDS:      &datajob.Source{InstanceDesc: "src-db"},
				TargetDS:      &datajob.Source{InstanceDesc: "dst-db"},
			},
		},
		job: datajob.Job{
			DataJobID:      11,
			DataJobName:    "sync-job",
			DataJobDesc:    "nightly sync",
			DataJobType:    "SYNC",
			DataTaskState:  "RUNNING",
			CurrTaskStatus: "FULL_RUNNING",
			LifeCycleState: "ACTIVE",
			UserName:       "admin",
			ConsoleJobID:   21,
			SourceDS: &datajob.Source{
				InstanceDesc:   "src-db",
				DataSourceType: "MYSQL",
			},
			TargetDS: &datajob.Source{
				InstanceDesc:   "dst-db",
				DataSourceType: "STARROCKS",
			},
			DataTasks: []datajob.Task{
				{DataTaskID: 101, DataTaskName: "full-task", DataTaskType: "FULL", DataTaskStatus: "RUNNING"},
			},
		},
		schema: datajob.JobSchema{
			SourceSchema:          "src_schema",
			TargetSchema:          "dst_schema",
			MappingConfig:         "{\"rules\":1}",
			DefaultTopic:          "topic-a",
			DefaultTopicPartition: 8,
			SchemaWhiteListLevel:  "TABLE",
		},
		createResult:         datajob.CreateJobResult{JobID: "99", Data: "99"},
		updateIncrePosResult: datajob.UpdateIncrePosResult{Data: "updated"},
	}
	dataSources := &fakeDataSources{
		list: []datasource.DataSource{
			{ID: 7, InstanceID: "cc-mysql-1", DataSourceType: "MYSQL", HostType: "RDS", DeployType: "ALIYUN", LifeCycleState: "ACTIVE", InstanceDesc: "mysql source"},
		},
		item:      datasource.DataSource{ID: 7, InstanceID: "cc-mysql-1", DataSourceType: "MYSQL", HostType: "RDS", DeployType: "ALIYUN", LifeCycleState: "ACTIVE", InstanceDesc: "mysql source"},
		addResult: "ds-77",
	}
	clusters := &fakeClusters{
		list: []cluster.Cluster{
			{ID: 3, ClusterName: "prod-cluster", Region: "cn-hangzhou", CloudOrIDCName: "ALIYUN", WorkerCount: 5, RunningCount: 4, AbnormalCount: 1, OwnerName: "admin"},
		},
	}
	workers := &fakeWorkers{
		list: []worker.Worker{
			{ID: 5, WorkerName: "worker-1", WorkerState: "RUNNING", WorkerType: "FULL", ClusterID: 2, PrivateIP: "10.0.0.5", HealthLevel: "GREEN", WorkerLoad: 0.8},
		},
	}
	consoleJobs := &fakeConsoleJobs{
		job: consolejob.Job{
			ID:           21,
			Label:        "WORKER_INSTALL",
			TaskState:    "RUNNING",
			JobToken:     "abc",
			WorkerName:   "worker-1",
			ResourceType: "WORKER",
			ResourceID:   5,
			TaskVOList: []consolejob.Task{
				{ID: 31, TaskState: "RUNNING", StepName: "Install", Host: "10.0.0.5", ExecuteOrder: 1, Cancelable: true},
			},
		},
	}
	jobConfigs := &fakeJobConfigs{
		specs: []jobconfig.Spec{
			{ID: 1, SpecKind: "SYNC", SpecKindCN: "同步", Spec: "STANDARD", FullMemoryMB: 2048, IncreMemoryMB: 1024, CheckMemoryMB: 512},
		},
		transformResult: jobconfig.TransformJobTypeResponse{Data: json.RawMessage(`{"normalized":"SYNC"}`)},
	}
	schemas := &fakeSchemas{
		items: []ccschema.ApiTransferObjIndexDO{
			{DataJobID: 11, DataJobName: "sync-job", SrcFullTransferObjName: "demo.orders", DstFullTransferObjName: "dw.orders", SrcDsType: "MYSQL", DstDsType: "STARROCKS"},
		},
	}

	runtime := &fakeRuntime{
		cfg:         config.AppConfig{APIBaseURL: "https://cc.example.com", AccessKey: "abcdefghijkl", SecretKey: "qrstuvwxyz1234"},
		dataJobs:    dataJobs,
		dataSources: dataSources,
		clusters:    clusters,
		workers:     workers,
		consoleJobs: consoleJobs,
		jobConfigs:  jobConfigs,
		schemas:     schemas,
	}
	io := testsupport.NewTestConsole(
		"help jobs",
		`jobs list --name sync-job --type SYNC --desc "nightly sync" --source-id 101 --target-id 202`,
		`jobs create --body '{"clusterId":1,"srcDsId":101,"dstDsId":202,"jobType":"SYNC","dataJobDesc":"sdk job"}'`,
		"jobs show 11",
		"jobs schema 11",
		"jobs replay 11 --auto-start --reset-to-created",
		"jobs attach-incre-task 11",
		"jobs detach-incre-task 11",
		`jobs update-incre-pos --body '{"taskId":101,"posType":"MYSQL_LOG_FILE_POS","journalFile":"binlog.000001"}'`,
		"jobs start 11",
		"jobs stop 11",
		"jobs delete 11",
		"datasources list --type MYSQL",
		`datasources add --body '{"type":"MYSQL","host":"127.0.0.1:3306","instanceDesc":"mysql source"}'`,
		"datasources delete 7",
		"datasources show 7",
		"clusters list --name prod",
		"workers list --cluster-id 2",
		"workers start 5",
		"workers stop 5",
		"workers delete 5",
		"workers modify-mem-oversold 5 --percent 120",
		"workers update-alert 5 --phone=true --email=false --im=true --sms=false",
		"consolejobs show 21",
		"job-config specs --type SYNC --initial-sync=true --short-term-sync=false",
		"job-config transform-job-type --source-type MYSQL --target-type STARROCKS",
		"schemas list-trans-objs-by-meta --src-db demo --src-trans-obj orders",
		"config lang show",
		"config lang set zh",
		"config lang set en",
		"config show",
		"exit",
	)

	shell := repl.NewShell(io, runtime)
	if err := shell.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	out := io.Output()
	if dataJobs.lastStartedJobID != 11 || dataJobs.lastStoppedJobID != 11 || dataJobs.lastDeletedJobID != 11 || dataJobs.lastReplayedJobID != 11 {
		t.Fatalf("unexpected job actions: start=%d stop=%d delete=%d replay=%d", dataJobs.lastStartedJobID, dataJobs.lastStoppedJobID, dataJobs.lastDeletedJobID, dataJobs.lastReplayedJobID)
	}
	if workers.lastStartedWorkerID != 5 || workers.lastStoppedWorkerID != 5 || workers.lastDeletedWorkerID != 5 {
		t.Fatalf("unexpected worker actions: start=%d stop=%d delete=%d", workers.lastStartedWorkerID, workers.lastStoppedWorkerID, workers.lastDeletedWorkerID)
	}
	if !dataJobs.lastReplayOptions.AutoStart || !dataJobs.lastReplayOptions.ResetToCreated {
		t.Fatalf("unexpected replay options: %+v", dataJobs.lastReplayOptions)
	}
	if dataJobs.lastCreatedJob.ClusterID != 1 || dataJobs.lastCreatedJob.JobType != "SYNC" || dataJobs.lastCreatedJob.DataJobDesc != "sdk job" {
		t.Fatalf("unexpected create job request: %+v", dataJobs.lastCreatedJob)
	}
	if dataJobs.lastAttachedIncreJobID != 11 || dataJobs.lastDetachedIncreJobID != 11 {
		t.Fatalf("unexpected incre task actions: attach=%d detach=%d", dataJobs.lastAttachedIncreJobID, dataJobs.lastDetachedIncreJobID)
	}
	if dataJobs.lastUpdateIncrePos.TaskID != 101 || dataJobs.lastUpdateIncrePos.PosType != "MYSQL_LOG_FILE_POS" {
		t.Fatalf("unexpected incre pos update: %+v", dataJobs.lastUpdateIncrePos)
	}
	if dataJobs.lastListOptions.DataJobName != "sync-job" || dataJobs.lastListOptions.Desc != "nightly sync" || dataJobs.lastListOptions.SourceInstanceID != 101 || dataJobs.lastListOptions.TargetInstanceID != 202 {
		t.Fatalf("unexpected list options: %+v", dataJobs.lastListOptions)
	}
	if dataSources.lastListOptions.Type != "MYSQL" {
		t.Fatalf("unexpected datasource list options: %+v", dataSources.lastListOptions)
	}
	if dataSources.lastAddOptions.DataSourceAddData.Type != "MYSQL" || dataSources.lastDeletedID != 7 {
		t.Fatalf("unexpected datasource actions: add=%+v delete=%d", dataSources.lastAddOptions, dataSources.lastDeletedID)
	}
	if clusters.lastListOptions.ClusterName != "prod" {
		t.Fatalf("unexpected cluster list options: %+v", clusters.lastListOptions)
	}
	if workers.lastListOptions.ClusterID != 2 {
		t.Fatalf("unexpected worker list options: %+v", workers.lastListOptions)
	}
	if workers.lastMemOverSoldWorkerID != 5 || workers.lastMemOverSoldPercent != 120 {
		t.Fatalf("unexpected mem oversold update: worker=%d percent=%d", workers.lastMemOverSoldWorkerID, workers.lastMemOverSoldPercent)
	}
	if workers.lastAlertWorkerID != 5 || !workers.lastAlertPhone || workers.lastAlertEmail || !workers.lastAlertIM || workers.lastAlertSMS {
		t.Fatalf("unexpected worker alert update: worker=%d phone=%v email=%v im=%v sms=%v", workers.lastAlertWorkerID, workers.lastAlertPhone, workers.lastAlertEmail, workers.lastAlertIM, workers.lastAlertSMS)
	}
	if jobConfigs.lastOptions.DataJobType != "SYNC" || jobConfigs.lastOptions.InitialSync == nil || !*jobConfigs.lastOptions.InitialSync || jobConfigs.lastOptions.ShortTermSync == nil || *jobConfigs.lastOptions.ShortTermSync {
		t.Fatalf("unexpected job config options: %+v", jobConfigs.lastOptions)
	}
	if jobConfigs.lastTransformOption.SourceType != "MYSQL" || jobConfigs.lastTransformOption.TargetType != "STARROCKS" {
		t.Fatalf("unexpected transform job type options: %+v", jobConfigs.lastTransformOption)
	}
	if schemas.lastOptions.SrcDb != "demo" || schemas.lastOptions.SrcTransObj != "orders" {
		t.Fatalf("unexpected schema list options: %+v", schemas.lastOptions)
	}

	for _, want := range []string{
		"jobs commands",
		"--name       Filter by data job name.",
		"sync-job",
		"Job created successfully",
		"Job ID",
		"Job details:",
		"nightly sync",
		"Job schema:",
		"topic-a",
		"Job 11 replay requested successfully",
		"Job 11 incremental task attached successfully",
		"Job 11 incremental task detached successfully",
		"Increment position updated successfully",
		"Job 11 started successfully",
		"Job 11 stopped successfully",
		"Job 11 deleted successfully",
		"mysql source",
		"Data source created successfully",
		"Data source 7 deleted successfully",
		"Data source details:",
		"prod-cluster",
		"worker-1",
		"Worker 5 started successfully",
		"Worker 5 stopped successfully",
		"Worker 5 deleted successfully",
		"Worker 5 memory oversold percentage updated successfully",
		"Worker 5 alert config updated successfully",
		"Console job details:",
		"WORKER_INSTALL",
		"STANDARD",
		"Transform job type result:",
		"normalized",
		"demo.orders",
		"Current language: en",
		"语言已切换为 中文。",
		"Language switched to English.",
		"apiBaseUrl: https://cc.example.com",
		"accessKey: abcd****ijkl",
		"language: en",
		"httpTimeoutSeconds: 10",
		"httpReadMaxRetries: 0",
		"httpReadRetryBackoffMillis: 250",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q in %q", want, out)
		}
	}
	if strings.Contains(out, "secretKey:") {
		t.Fatalf("output unexpectedly contains secret key line: %q", out)
	}
}

func TestShellReportsInvalidCommandsWithoutExiting(t *testing.T) {
	runtime := &fakeRuntime{
		cfg:               config.AppConfig{APIBaseURL: "https://cc.example.com", AccessKey: "abcdefghijkl", SecretKey: "qrstuvwxyz1234"},
		dataJobs:          &fakeDataJobs{},
		dataSources:       &fakeDataSources{},
		clusters:          &fakeClusters{},
		workers:           &fakeWorkers{},
		consoleJobs:       &fakeConsoleJobs{},
		jobConfigs:        &fakeJobConfigs{},
		reinitializeValue: true,
	}
	io := testsupport.NewTestConsole(
		"jobs start abc",
		"jobs replay 11 --bad",
		`jobs list --desc "unterminated`,
		"workers list",
		"job-config specs",
		"job-config specs --type SYNC --initial-sync=maybe",
		"unknown",
		"config init",
		"exit",
	)

	shell := repl.NewShell(io, runtime)
	if err := shell.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	out := io.Output()
	if runtime.reinitializeCalls != 1 {
		t.Fatalf("reinitializeCalls = %d, want 1", runtime.reinitializeCalls)
	}
	for _, want := range []string{
		"jobId must be a positive integer",
		"unknown option: --bad",
		"unterminated quote",
		"clusterId is required",
		"dataJobType is required",
		"initialSync must be a boolean",
		"Unknown command: unknown",
		"Configuration updated.",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q in %q", want, out)
		}
	}
}

func TestShellExecutesArgsWithoutInteractiveLoop(t *testing.T) {
	dataJobs := &fakeDataJobs{
		jobs: []datajob.Job{
			{
				DataJobID:     22,
				DataJobName:   "batch-job",
				DataJobType:   "CHECK",
				DataTaskState: "CREATED",
				SourceDS:      &datajob.Source{InstanceDesc: "src-check"},
				TargetDS:      &datajob.Source{InstanceDesc: "dst-check"},
			},
		},
	}
	runtime := &fakeRuntime{
		cfg:         config.AppConfig{APIBaseURL: "https://cc.example.com", AccessKey: "abcdefghijkl", SecretKey: "qrstuvwxyz1234"},
		dataJobs:    dataJobs,
		dataSources: &fakeDataSources{},
		clusters:    &fakeClusters{},
		workers:     &fakeWorkers{},
		consoleJobs: &fakeConsoleJobs{},
		jobConfigs:  &fakeJobConfigs{},
	}
	io := testsupport.NewTestConsole()

	shell := repl.NewShell(io, runtime)
	if err := shell.ExecuteArgs([]string{"jobs", "list", "--type", "CHECK"}); err != nil {
		t.Fatalf("ExecuteArgs() error = %v", err)
	}

	out := io.Output()
	for _, want := range []string{
		"batch-job",
		"1 jobs",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q in %q", want, out)
		}
	}
	if dataJobs.lastListOptions.DataJobType != "CHECK" {
		t.Fatalf("lastListOptions = %+v, want type CHECK", dataJobs.lastListOptions)
	}
	if strings.Contains(out, "cloudcanal>") || strings.Contains(out, "Type 'help'") {
		t.Fatalf("output unexpectedly contains interactive text: %q", out)
	}
}

func TestShellShowsGroupedUsageOnSeparateLines(t *testing.T) {
	runtime := &fakeRuntime{
		cfg:         config.AppConfig{APIBaseURL: "https://cc.example.com", AccessKey: "abcdefghijkl", SecretKey: "qrstuvwxyz1234"},
		dataJobs:    &fakeDataJobs{},
		dataSources: &fakeDataSources{},
		clusters:    &fakeClusters{},
		workers:     &fakeWorkers{},
		consoleJobs: &fakeConsoleJobs{},
		jobConfigs:  &fakeJobConfigs{},
	}
	io := testsupport.NewTestConsole()

	shell := repl.NewShell(io, runtime)
	if err := shell.ExecuteArgs([]string{"jobs"}); err != nil {
		t.Fatalf("ExecuteArgs(jobs) error = %v", err)
	}

	out := io.Output()
	for _, want := range []string{
		"Usage:",
		"  jobs list",
		"  jobs show <jobId>",
		"  jobs schema <jobId>",
		"  jobs replay <jobId> [--auto-start] [--reset-to-created]",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q in %q", want, out)
		}
	}
}

func TestShellClearsScreenWithAliases(t *testing.T) {
	runtime := &fakeRuntime{
		cfg:         config.AppConfig{APIBaseURL: "https://cc.example.com", AccessKey: "abcdefghijkl", SecretKey: "qrstuvwxyz1234"},
		dataJobs:    &fakeDataJobs{},
		dataSources: &fakeDataSources{},
		clusters:    &fakeClusters{},
		workers:     &fakeWorkers{},
		consoleJobs: &fakeConsoleJobs{},
		jobConfigs:  &fakeJobConfigs{},
	}

	io := testsupport.NewTestConsole("clear", "cls", "exit")
	shell := repl.NewShell(io, runtime)
	if err := shell.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if got := strings.Count(io.Output(), "\x1b[H\x1b[2J"); got != 2 {
		t.Fatalf("clear sequence count = %d, want 2 in %q", got, io.Output())
	}

	io = testsupport.NewTestConsole()
	shell = repl.NewShell(io, runtime)
	if err := shell.ExecuteArgs([]string{"clear"}); err != nil {
		t.Fatalf("ExecuteArgs(clear) error = %v", err)
	}
	if !strings.Contains(io.Output(), "\x1b[H\x1b[2J") {
		t.Fatalf("output missing clear sequence in %q", io.Output())
	}
}

func TestShellIgnoresPromptAbortInInteractiveMode(t *testing.T) {
	runtime := &fakeRuntime{
		cfg:         config.AppConfig{APIBaseURL: "https://cc.example.com", AccessKey: "abcdefghijkl", SecretKey: "qrstuvwxyz1234"},
		dataJobs:    &fakeDataJobs{},
		dataSources: &fakeDataSources{},
		clusters:    &fakeClusters{},
		workers:     &fakeWorkers{},
		consoleJobs: &fakeConsoleJobs{},
		jobConfigs:  &fakeJobConfigs{},
	}
	io := &abortingConsole{}

	shell := repl.NewShell(io, runtime)
	if err := shell.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	out := io.Output()
	if strings.Contains(out, "prompt aborted") || strings.Contains(out, "Fatal error") {
		t.Fatalf("output should not contain prompt abort error in %q", out)
	}
	if got := strings.Count(out, "cloudcanal> "); got != 1 {
		t.Fatalf("prompt count = %d, want 1 in %q", got, out)
	}
}

func TestShellHelpOverviewHidesInternalCommands(t *testing.T) {
	runtime := &fakeRuntime{
		cfg:         config.AppConfig{APIBaseURL: "https://cc.example.com", AccessKey: "abcdefghijkl", SecretKey: "qrstuvwxyz1234"},
		dataJobs:    &fakeDataJobs{},
		dataSources: &fakeDataSources{},
		clusters:    &fakeClusters{},
		workers:     &fakeWorkers{},
		consoleJobs: &fakeConsoleJobs{},
		jobConfigs:  &fakeJobConfigs{},
	}
	io := testsupport.NewTestConsole()

	shell := repl.NewShell(io, runtime)
	if err := shell.ExecuteArgs([]string{"help"}); err != nil {
		t.Fatalf("ExecuteArgs(help) error = %v", err)
	}

	out := io.Output()
	for _, want := range []string{
		"jobs list",
		"datasources list",
		"config init",
		"config lang show",
		"TAB               Complete commands and options",
		"Ctrl+C            Exit interactive mode",
		"exit              Leave interactive mode",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q in %q", want, out)
		}
	}
	for _, hidden := range []string{
		"completion zsh",
		"help completion",
		"help lang",
		"clear             Clear the current screen",
		"cls               Alias of clear",
		"quit",
	} {
		if strings.Contains(out, hidden) {
			t.Fatalf("output should not contain %q in %q", hidden, out)
		}
	}
}

func TestShellSupportsHelpFlags(t *testing.T) {
	runtime := &fakeRuntime{
		cfg:         config.AppConfig{APIBaseURL: "https://cc.example.com", AccessKey: "abcdefghijkl", SecretKey: "qrstuvwxyz1234"},
		dataJobs:    &fakeDataJobs{},
		dataSources: &fakeDataSources{},
		clusters:    &fakeClusters{},
		workers:     &fakeWorkers{},
		consoleJobs: &fakeConsoleJobs{},
		jobConfigs:  &fakeJobConfigs{},
	}
	io := testsupport.NewTestConsole()

	shell := repl.NewShell(io, runtime)
	if err := shell.ExecuteArgs([]string{"--help"}); err != nil {
		t.Fatalf("ExecuteArgs(--help) error = %v", err)
	}
	if !strings.Contains(io.Output(), "CloudCanal CLI help") {
		t.Fatalf("output missing top-level help in %q", io.Output())
	}

	io = testsupport.NewTestConsole()
	shell = repl.NewShell(io, runtime)
	if err := shell.ExecuteArgs([]string{"jobs", "--help"}); err != nil {
		t.Fatalf("ExecuteArgs(jobs --help) error = %v", err)
	}
	if !strings.Contains(io.Output(), "jobs commands") {
		t.Fatalf("output missing jobs help in %q", io.Output())
	}

	io = testsupport.NewTestConsole()
	shell = repl.NewShell(io, runtime)
	if err := shell.ExecuteArgs([]string{"jobs", "list", "--help"}); err != nil {
		t.Fatalf("ExecuteArgs(jobs list --help) error = %v", err)
	}
	if !strings.Contains(io.Output(), "Usage: jobs list [--name NAME] [--type TYPE] [--desc DESC] [--source-id ID] [--target-id ID]") {
		t.Fatalf("output missing jobs list usage in %q", io.Output())
	}

	io = testsupport.NewTestConsole()
	shell = repl.NewShell(io, runtime)
	if err := shell.ExecuteArgs([]string{"config", "lang", "--help"}); err != nil {
		t.Fatalf("ExecuteArgs(config lang --help) error = %v", err)
	}
	if !strings.Contains(io.Output(), "config lang commands") {
		t.Fatalf("output missing config lang help in %q", io.Output())
	}
}

func TestShellSuggestsClosestCommands(t *testing.T) {
	runtime := &fakeRuntime{
		cfg:         config.AppConfig{APIBaseURL: "https://cc.example.com", AccessKey: "abcdefghijkl", SecretKey: "qrstuvwxyz1234"},
		dataJobs:    &fakeDataJobs{},
		dataSources: &fakeDataSources{},
		clusters:    &fakeClusters{},
		workers:     &fakeWorkers{},
		consoleJobs: &fakeConsoleJobs{},
		jobConfigs:  &fakeJobConfigs{},
	}
	io := testsupport.NewTestConsole()

	shell := repl.NewShell(io, runtime)
	if err := shell.ExecuteArgs([]string{"jbos", "list"}); err != nil {
		t.Fatalf("ExecuteArgs(jbos list) error = %v", err)
	}
	if !strings.Contains(io.Output(), "Did you mean: jobs") {
		t.Fatalf("output missing root suggestion in %q", io.Output())
	}

	io = testsupport.NewTestConsole()
	shell = repl.NewShell(io, runtime)
	if err := shell.ExecuteArgs([]string{"jobs", "shwo"}); err != nil {
		t.Fatalf("ExecuteArgs(jobs shwo) error = %v", err)
	}
	out := io.Output()
	if !strings.Contains(out, "Unknown jobs command: shwo") || !strings.Contains(out, "Did you mean: jobs show") {
		t.Fatalf("output missing jobs suggestion in %q", out)
	}
}

func TestShellUnknownHelpTopicShowsSuggestion(t *testing.T) {
	runtime := &fakeRuntime{
		cfg:         config.AppConfig{APIBaseURL: "https://cc.example.com", AccessKey: "abcdefghijkl", SecretKey: "qrstuvwxyz1234"},
		dataJobs:    &fakeDataJobs{},
		dataSources: &fakeDataSources{},
		clusters:    &fakeClusters{},
		workers:     &fakeWorkers{},
		consoleJobs: &fakeConsoleJobs{},
		jobConfigs:  &fakeJobConfigs{},
	}
	io := testsupport.NewTestConsole()

	shell := repl.NewShell(io, runtime)
	if err := shell.ExecuteArgs([]string{"help", "clustrs"}); err != nil {
		t.Fatalf("ExecuteArgs(help clustrs) error = %v", err)
	}

	out := io.Output()
	if !strings.Contains(out, "Unknown help topic: clustrs") || !strings.Contains(out, "Did you mean: help clusters") {
		t.Fatalf("output missing help suggestion in %q", out)
	}
}

func TestShellSwitchesLanguageForFollowUpOutput(t *testing.T) {
	runtime := &fakeRuntime{
		cfg:         config.AppConfig{APIBaseURL: "https://cc.example.com", AccessKey: "abcdefghijkl", SecretKey: "qrstuvwxyz1234"},
		dataJobs:    &fakeDataJobs{},
		dataSources: &fakeDataSources{},
		clusters:    &fakeClusters{},
		workers:     &fakeWorkers{},
		consoleJobs: &fakeConsoleJobs{},
		jobConfigs:  &fakeJobConfigs{},
	}
	io := testsupport.NewTestConsole("config lang set zh", "help config", "exit")

	shell := repl.NewShell(io, runtime)
	if err := shell.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	out := io.Output()
	for _, want := range []string{
		"语言已切换为 中文。",
		"config 命令",
		"config lang set <en|zh>",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q in %q", want, out)
		}
	}
}

func TestShellOutputsJSONForCommandsAndErrors(t *testing.T) {
	dataJobs := &fakeDataJobs{
		jobs: []datajob.Job{
			{
				DataJobID:     22,
				DataJobName:   "batch-job",
				DataJobType:   "CHECK",
				DataTaskState: "CREATED",
			},
		},
	}
	runtime := &fakeRuntime{
		cfg:         config.AppConfig{APIBaseURL: "https://cc.example.com", AccessKey: "abcdefghijkl", SecretKey: "qrstuvwxyz1234"},
		dataJobs:    dataJobs,
		dataSources: &fakeDataSources{},
		clusters:    &fakeClusters{},
		workers:     &fakeWorkers{},
		consoleJobs: &fakeConsoleJobs{},
		jobConfigs:  &fakeJobConfigs{},
	}

	io := testsupport.NewTestConsole()
	shell := repl.NewShell(io, runtime)
	if err := shell.ExecuteArgs([]string{"jobs", "list", "--type", "CHECK", "--output", "json"}); err != nil {
		t.Fatalf("ExecuteArgs(json jobs list) error = %v", err)
	}

	var jobs []map[string]any
	if err := json.Unmarshal([]byte(io.Output()), &jobs); err != nil {
		t.Fatalf("json.Unmarshal(jobs output) error = %v, output = %q", err, io.Output())
	}
	if len(jobs) != 1 || jobs[0]["dataJobName"] != "batch-job" {
		t.Fatalf("jobs json = %#v, want batch-job", jobs)
	}

	io = testsupport.NewTestConsole()
	shell = repl.NewShell(io, runtime)
	err := shell.ExecuteArgs([]string{"jobs", "start", "abc", "--output", "json"})
	if err == nil {
		t.Fatal("ExecuteArgs(json jobs start) error = nil, want non-nil")
	}
	shell.PrintFatalError(err)

	var payload map[string]any
	if err := json.Unmarshal([]byte(io.Output()), &payload); err != nil {
		t.Fatalf("json.Unmarshal(error output) error = %v, output = %q", err, io.Output())
	}
	if payload["error"] != "jobId must be a positive integer" {
		t.Fatalf("error payload = %#v, want jobId validation", payload)
	}
	if payload["fatal"] != true {
		t.Fatalf("error payload = %#v, want fatal=true", payload)
	}
}

type fakeRuntime struct {
	cfg               config.AppConfig
	dataJobs          datajob.Operations
	dataSources       datasource.Operations
	clusters          cluster.Operations
	workers           worker.Operations
	consoleJobs       consolejob.Operations
	jobConfigs        jobconfig.Operations
	schemas           ccschema.Operations
	reinitializeCalls int
	reinitializeValue bool
}

type abortingConsole struct {
	output  strings.Builder
	aborted bool
}

func (a *abortingConsole) ReadLine(prompt string) (string, error) {
	a.output.WriteString(prompt)
	if !a.aborted {
		a.aborted = true
		return "", liner.ErrPromptAborted
	}
	return "exit", nil
}

func (a *abortingConsole) ReadSecret(prompt string) (string, error) {
	a.output.WriteString(prompt)
	return "", liner.ErrPromptAborted
}

func (a *abortingConsole) Println(text string) {
	a.output.WriteString(text)
	a.output.WriteString("\n")
}

func (a *abortingConsole) ClearScreen() {
	a.output.WriteString("\033[H\033[2J")
}

func (a *abortingConsole) Output() string {
	return a.output.String()
}

func (f *fakeRuntime) Config() config.AppConfig {
	return f.cfg
}

func (f *fakeRuntime) DataJobs() datajob.Operations {
	return f.dataJobs
}

func (f *fakeRuntime) DataSources() datasource.Operations {
	return f.dataSources
}

func (f *fakeRuntime) Clusters() cluster.Operations {
	return f.clusters
}

func (f *fakeRuntime) Workers() worker.Operations {
	return f.workers
}

func (f *fakeRuntime) ConsoleJobs() consolejob.Operations {
	return f.consoleJobs
}

func (f *fakeRuntime) JobConfigs() jobconfig.Operations {
	return f.jobConfigs
}

func (f *fakeRuntime) Schemas() ccschema.Operations {
	return f.schemas
}

func (f *fakeRuntime) Reinitialize(io console.IO) (bool, error) {
	f.reinitializeCalls++
	return f.reinitializeValue, nil
}

func (f *fakeRuntime) SetLanguage(language string) error {
	f.cfg.Language = language
	_ = i18n.SetLanguage(language)
	return nil
}

var _ app.RuntimeContext = (*fakeRuntime)(nil)

type fakeDataJobs struct {
	jobs                   []datajob.Job
	job                    datajob.Job
	schema                 datajob.JobSchema
	createResult           datajob.CreateJobResult
	updateIncrePosResult   datajob.UpdateIncrePosResult
	lastListOptions        datajob.ListOptions
	lastCreatedJob         datajob.CreateJobRequest
	lastStartedJobID       int64
	lastStoppedJobID       int64
	lastDeletedJobID       int64
	lastReplayedJobID      int64
	lastAttachedIncreJobID int64
	lastDetachedIncreJobID int64
	lastUpdateIncrePos     datajob.UpdateIncrePosRequest
	lastReplayOptions      datajob.ReplayOptions
}

func (f *fakeDataJobs) ListJobs(options datajob.ListOptions) ([]datajob.Job, error) {
	f.lastListOptions = options
	return f.jobs, nil
}

func (f *fakeDataJobs) GetJob(jobID int64) (datajob.Job, error) {
	return f.job, nil
}

func (f *fakeDataJobs) GetJobSchema(jobID int64) (datajob.JobSchema, error) {
	return f.schema, nil
}

func (f *fakeDataJobs) CreateJob(request datajob.CreateJobRequest) (datajob.CreateJobResult, error) {
	f.lastCreatedJob = request
	return f.createResult, nil
}

func (f *fakeDataJobs) StartJob(jobID int64) error {
	f.lastStartedJobID = jobID
	return nil
}

func (f *fakeDataJobs) StopJob(jobID int64) error {
	f.lastStoppedJobID = jobID
	return nil
}

func (f *fakeDataJobs) DeleteJob(jobID int64) error {
	f.lastDeletedJobID = jobID
	return nil
}

func (f *fakeDataJobs) ReplayJob(jobID int64, options datajob.ReplayOptions) error {
	f.lastReplayedJobID = jobID
	f.lastReplayOptions = options
	return nil
}

func (f *fakeDataJobs) AttachIncreJob(jobID int64) error {
	f.lastAttachedIncreJobID = jobID
	return nil
}

func (f *fakeDataJobs) DetachIncreJob(jobID int64) error {
	f.lastDetachedIncreJobID = jobID
	return nil
}

func (f *fakeDataJobs) UpdateIncrePos(request datajob.UpdateIncrePosRequest) (datajob.UpdateIncrePosResult, error) {
	f.lastUpdateIncrePos = request
	return f.updateIncrePosResult, nil
}

type fakeDataSources struct {
	list            []datasource.DataSource
	item            datasource.DataSource
	addResult       string
	lastListOptions datasource.ListOptions
	lastAddOptions  datasource.AddOptions
	lastGetID       int64
	lastDeletedID   int64
}

func (f *fakeDataSources) List(options datasource.ListOptions) ([]datasource.DataSource, error) {
	f.lastListOptions = options
	return f.list, nil
}

func (f *fakeDataSources) Get(dataSourceID int64) (datasource.DataSource, error) {
	f.lastGetID = dataSourceID
	return f.item, nil
}

func (f *fakeDataSources) Add(options datasource.AddOptions) (string, error) {
	f.lastAddOptions = options
	return f.addResult, nil
}

func (f *fakeDataSources) Delete(dataSourceID int64) error {
	f.lastDeletedID = dataSourceID
	return nil
}

type fakeClusters struct {
	list            []cluster.Cluster
	lastListOptions cluster.ListOptions
}

func (f *fakeClusters) List(options cluster.ListOptions) ([]cluster.Cluster, error) {
	f.lastListOptions = options
	return f.list, nil
}

type fakeWorkers struct {
	list                    []worker.Worker
	lastListOptions         worker.ListOptions
	lastStartedWorkerID     int64
	lastStoppedWorkerID     int64
	lastDeletedWorkerID     int64
	lastMemOverSoldWorkerID int64
	lastMemOverSoldPercent  int
	lastAlertWorkerID       int64
	lastAlertPhone          bool
	lastAlertEmail          bool
	lastAlertIM             bool
	lastAlertSMS            bool
}

func (f *fakeWorkers) List(options worker.ListOptions) ([]worker.Worker, error) {
	f.lastListOptions = options
	return f.list, nil
}

func (f *fakeWorkers) Start(workerID int64) error {
	f.lastStartedWorkerID = workerID
	return nil
}

func (f *fakeWorkers) Stop(workerID int64) error {
	f.lastStoppedWorkerID = workerID
	return nil
}

func (f *fakeWorkers) Delete(workerID int64) error {
	f.lastDeletedWorkerID = workerID
	return nil
}

func (f *fakeWorkers) ModifyMemOverSold(workerID int64, memOverSoldPercent int) error {
	f.lastMemOverSoldWorkerID = workerID
	f.lastMemOverSoldPercent = memOverSoldPercent
	return nil
}

func (f *fakeWorkers) UpdateWorkerAlert(workerID int64, phone, email, im, sms bool) error {
	f.lastAlertWorkerID = workerID
	f.lastAlertPhone = phone
	f.lastAlertEmail = email
	f.lastAlertIM = im
	f.lastAlertSMS = sms
	return nil
}

type fakeConsoleJobs struct {
	job       consolejob.Job
	lastGetID int64
}

func (f *fakeConsoleJobs) Get(consoleJobID int64) (consolejob.Job, error) {
	f.lastGetID = consoleJobID
	return f.job, nil
}

type fakeJobConfigs struct {
	specs               []jobconfig.Spec
	transformResult     jobconfig.TransformJobTypeResponse
	lastOptions         jobconfig.ListSpecsOptions
	lastTransformOption jobconfig.TransformJobTypeOptions
}

func (f *fakeJobConfigs) ListSpecs(options jobconfig.ListSpecsOptions) ([]jobconfig.Spec, error) {
	f.lastOptions = options
	return f.specs, nil
}

func (f *fakeJobConfigs) TransformJobType(options jobconfig.TransformJobTypeOptions) (jobconfig.TransformJobTypeResponse, error) {
	f.lastTransformOption = options
	return f.transformResult, nil
}

type fakeSchemas struct {
	items       []ccschema.ApiTransferObjIndexDO
	lastOptions ccschema.ListTransObjsByMetaOptions
}

func (f *fakeSchemas) ListTransObjsByMeta(options ccschema.ListTransObjsByMetaOptions) ([]ccschema.ApiTransferObjIndexDO, error) {
	f.lastOptions = options
	return f.items, nil
}
