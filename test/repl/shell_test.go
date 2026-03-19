package repl_test

import (
	"cloudcanal-openapi-cli/internal/app"
	"cloudcanal-openapi-cli/internal/cluster"
	"cloudcanal-openapi-cli/internal/config"
	"cloudcanal-openapi-cli/internal/console"
	"cloudcanal-openapi-cli/internal/consolejob"
	"cloudcanal-openapi-cli/internal/datajob"
	"cloudcanal-openapi-cli/internal/datasource"
	"cloudcanal-openapi-cli/internal/jobconfig"
	"cloudcanal-openapi-cli/internal/repl"
	"cloudcanal-openapi-cli/internal/worker"
	"cloudcanal-openapi-cli/test/testsupport"
	"strings"
	"testing"
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
	}
	dataSources := &fakeDataSources{
		list: []datasource.DataSource{
			{ID: 7, InstanceID: "cc-mysql-1", DataSourceType: "MYSQL", HostType: "RDS", DeployType: "ALIYUN", LifeCycleState: "ACTIVE", InstanceDesc: "mysql source"},
		},
		item: datasource.DataSource{ID: 7, InstanceID: "cc-mysql-1", DataSourceType: "MYSQL", HostType: "RDS", DeployType: "ALIYUN", LifeCycleState: "ACTIVE", InstanceDesc: "mysql source"},
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
	}

	runtime := &fakeRuntime{
		cfg:         config.AppConfig{APIBaseURL: "https://cc.example.com", AccessKey: "abcdefghijkl", SecretKey: "qrstuvwxyz1234"},
		dataJobs:    dataJobs,
		dataSources: dataSources,
		clusters:    clusters,
		workers:     workers,
		consoleJobs: consoleJobs,
		jobConfigs:  jobConfigs,
	}
	io := testsupport.NewTestConsole(
		"help",
		`jobs list --name sync-job --type SYNC --desc "nightly sync" --source-id 101 --target-id 202`,
		"jobs show 11",
		"jobs schema 11",
		"jobs replay 11 --auto-start --reset-to-created",
		"jobs start 11",
		"jobs stop 11",
		"jobs delete 11",
		"datasources list --type MYSQL",
		"datasources show 7",
		"clusters list --name prod",
		"workers list --cluster-id 2",
		"workers start 5",
		"workers stop 5",
		"consolejobs show 21",
		"job-config specs --type SYNC --initial-sync=true --short-term-sync=false",
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
	if workers.lastStartedWorkerID != 5 || workers.lastStoppedWorkerID != 5 {
		t.Fatalf("unexpected worker actions: start=%d stop=%d", workers.lastStartedWorkerID, workers.lastStoppedWorkerID)
	}
	if !dataJobs.lastReplayOptions.AutoStart || !dataJobs.lastReplayOptions.ResetToCreated {
		t.Fatalf("unexpected replay options: %+v", dataJobs.lastReplayOptions)
	}
	if dataJobs.lastListOptions.DataJobName != "sync-job" || dataJobs.lastListOptions.Desc != "nightly sync" || dataJobs.lastListOptions.SourceInstanceID != 101 || dataJobs.lastListOptions.TargetInstanceID != 202 {
		t.Fatalf("unexpected list options: %+v", dataJobs.lastListOptions)
	}
	if dataSources.lastListOptions.Type != "MYSQL" {
		t.Fatalf("unexpected datasource list options: %+v", dataSources.lastListOptions)
	}
	if clusters.lastListOptions.ClusterName != "prod" {
		t.Fatalf("unexpected cluster list options: %+v", clusters.lastListOptions)
	}
	if workers.lastListOptions.ClusterID != 2 {
		t.Fatalf("unexpected worker list options: %+v", workers.lastListOptions)
	}
	if jobConfigs.lastOptions.DataJobType != "SYNC" || jobConfigs.lastOptions.InitialSync == nil || !*jobConfigs.lastOptions.InitialSync || jobConfigs.lastOptions.ShortTermSync == nil || *jobConfigs.lastOptions.ShortTermSync {
		t.Fatalf("unexpected job config options: %+v", jobConfigs.lastOptions)
	}

	for _, want := range []string{
		"Available commands:",
		"sync-job",
		"Job details:",
		"nightly sync",
		"Job schema:",
		"topic-a",
		"Job 11 replay requested successfully",
		"Job 11 started successfully",
		"Job 11 stopped successfully",
		"Job 11 deleted successfully",
		"mysql source",
		"Data source details:",
		"prod-cluster",
		"worker-1",
		"Worker 5 started successfully",
		"Worker 5 stopped successfully",
		"Console job details:",
		"WORKER_INSTALL",
		"STANDARD",
		"apiBaseUrl: https://cc.example.com",
		"accessKey: abcd****ijkl",
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
		"job-config specs --initial-sync=maybe",
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

type fakeRuntime struct {
	cfg               config.AppConfig
	dataJobs          datajob.Operations
	dataSources       datasource.Operations
	clusters          cluster.Operations
	workers           worker.Operations
	consoleJobs       consolejob.Operations
	jobConfigs        jobconfig.Operations
	reinitializeCalls int
	reinitializeValue bool
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

func (f *fakeRuntime) Reinitialize(io console.IO) (bool, error) {
	f.reinitializeCalls++
	return f.reinitializeValue, nil
}

var _ app.RuntimeContext = (*fakeRuntime)(nil)

type fakeDataJobs struct {
	jobs              []datajob.Job
	job               datajob.Job
	schema            datajob.JobSchema
	lastListOptions   datajob.ListOptions
	lastStartedJobID  int64
	lastStoppedJobID  int64
	lastDeletedJobID  int64
	lastReplayedJobID int64
	lastReplayOptions datajob.ReplayOptions
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

type fakeDataSources struct {
	list            []datasource.DataSource
	item            datasource.DataSource
	lastListOptions datasource.ListOptions
	lastGetID       int64
}

func (f *fakeDataSources) List(options datasource.ListOptions) ([]datasource.DataSource, error) {
	f.lastListOptions = options
	return f.list, nil
}

func (f *fakeDataSources) Get(dataSourceID int64) (datasource.DataSource, error) {
	f.lastGetID = dataSourceID
	return f.item, nil
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
	list                []worker.Worker
	lastListOptions     worker.ListOptions
	lastStartedWorkerID int64
	lastStoppedWorkerID int64
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

type fakeConsoleJobs struct {
	job       consolejob.Job
	lastGetID int64
}

func (f *fakeConsoleJobs) Get(consoleJobID int64) (consolejob.Job, error) {
	f.lastGetID = consoleJobID
	return f.job, nil
}

type fakeJobConfigs struct {
	specs       []jobconfig.Spec
	lastOptions jobconfig.ListSpecsOptions
}

func (f *fakeJobConfigs) ListSpecs(options jobconfig.ListSpecsOptions) ([]jobconfig.Spec, error) {
	f.lastOptions = options
	return f.specs, nil
}
