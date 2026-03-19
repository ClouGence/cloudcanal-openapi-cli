package repl_test

import (
	"cloudcanal-openapi-cli/internal/cluster"
	"cloudcanal-openapi-cli/internal/config"
	"cloudcanal-openapi-cli/internal/consolejob"
	"cloudcanal-openapi-cli/internal/datajob"
	"cloudcanal-openapi-cli/internal/datasource"
	"cloudcanal-openapi-cli/internal/jobconfig"
	"cloudcanal-openapi-cli/internal/repl"
	ccschema "cloudcanal-openapi-cli/internal/schema"
	"cloudcanal-openapi-cli/internal/worker"
	"cloudcanal-openapi-cli/test/testsupport"
	"strings"
	"testing"

	"github.com/mattn/go-runewidth"
)

func TestTabularCommandsAlignMixedWidthOutput(t *testing.T) {
	runtime := newAlignmentRuntime()

	testCases := []struct {
		name string
		args []string
	}{
		{name: "jobs list", args: []string{"jobs", "list"}},
		{name: "jobs show", args: []string{"jobs", "show", "11"}},
		{name: "datasources list", args: []string{"datasources", "list"}},
		{name: "clusters list", args: []string{"clusters", "list"}},
		{name: "workers list", args: []string{"workers", "list", "--cluster-id", "1"}},
		{name: "consolejobs show", args: []string{"consolejobs", "show", "21"}},
		{name: "job-config specs", args: []string{"job-config", "specs", "--type", "SYNC"}},
		{name: "schemas list", args: []string{"schemas", "list-trans-objs-by-meta", "--src-db", "demo"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			io := testsupport.NewTestConsole()
			shell := repl.NewShell(io, runtime)
			if err := shell.ExecuteArgs(tc.args); err != nil {
				t.Fatalf("ExecuteArgs(%v) error = %v", tc.args, err)
			}
			assertTableBlocksAligned(t, io.Output())
		})
	}
}

func TestDetailCommandsAlignLabelColumn(t *testing.T) {
	runtime := newAlignmentRuntime()

	testCases := []struct {
		name string
		args []string
	}{
		{name: "jobs show", args: []string{"jobs", "show", "11"}},
		{name: "jobs schema", args: []string{"jobs", "schema", "11"}},
		{name: "datasources show", args: []string{"datasources", "show", "7"}},
		{name: "consolejobs show", args: []string{"consolejobs", "show", "21"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			io := testsupport.NewTestConsole()
			shell := repl.NewShell(io, runtime)
			if err := shell.ExecuteArgs(tc.args); err != nil {
				t.Fatalf("ExecuteArgs(%v) error = %v", tc.args, err)
			}
			assertLabelColumnAligned(t, io.Output())
		})
	}
}

func newAlignmentRuntime() *fakeRuntime {
	return &fakeRuntime{
		cfg: config.AppConfig{
			APIBaseURL: "https://cc.example.com",
			AccessKey:  "abcdefghijkl",
			SecretKey:  "qrstuvwxyz1234",
			Language:   "zh",
		},
		dataJobs: &fakeDataJobs{
			jobs: []datajob.Job{
				{
					DataJobID:     11,
					DataJobName:   "同步任务A",
					DataJobType:   "SYNC",
					DataTaskState: "RUNNING",
					SourceDS:      &datajob.Source{InstanceDesc: "源MySQL-杭州"},
					TargetDS:      &datajob.Source{InstanceDesc: "目标StarRocks-上海"},
				},
			},
			job: datajob.Job{
				DataJobID:      11,
				DataJobName:    "同步任务A",
				DataJobDesc:    "跨区域 nightly sync",
				DataJobType:    "SYNC",
				DataTaskState:  "RUNNING",
				CurrTaskStatus: "FULL_RUNNING",
				LifeCycleState: "ACTIVE",
				UserName:       "admin",
				ConsoleJobID:   21,
				SourceDS: &datajob.Source{
					InstanceDesc:   "源MySQL-杭州",
					DataSourceType: "MYSQL",
					HostType:       "RDS",
					Region:         "cn-hangzhou",
				},
				TargetDS: &datajob.Source{
					InstanceDesc:   "目标StarRocks-上海",
					DataSourceType: "STARROCKS",
					HostType:       "EMR",
					Region:         "cn-shanghai",
				},
				DataTasks: []datajob.Task{
					{DataTaskID: 101, DataTaskName: "全量同步任务A", DataTaskType: "FULL", DataTaskStatus: "RUNNING", WorkerIP: "198.18.0.11"},
				},
			},
			schema: datajob.JobSchema{
				SourceSchema:          "src_schema",
				TargetSchema:          "dst_schema",
				MappingConfig:         "{\"rules\":1}",
				DefaultTopic:          "topic-同步",
				DefaultTopicPartition: 8,
				SchemaWhiteListLevel:  "TABLE",
			},
		},
		dataSources: &fakeDataSources{
			list: []datasource.DataSource{
				{ID: 7, InstanceID: "cc-mysql-1", DataSourceType: "MYSQL", HostType: "RDS", DeployType: "ALIYUN", LifeCycleState: "ACTIVE", InstanceDesc: "杭州主库-同步源"},
			},
			item: datasource.DataSource{
				ID:               7,
				InstanceID:       "cc-mysql-1",
				DataSourceType:   "MYSQL",
				HostType:         "RDS",
				DeployType:       "ALIYUN",
				Region:           "cn-hangzhou",
				LifeCycleState:   "ACTIVE",
				InstanceDesc:     "杭州主库-同步源",
				AccountName:      "root",
				SecurityType:     "password",
				ConsoleJobID:     "21",
				ConsoleTaskState: "RUNNING",
			},
		},
		clusters: &fakeClusters{
			list: []cluster.Cluster{
				{ID: 1, ClusterName: "生产集群-A", Region: "cn-hangzhou", CloudOrIDCName: "阿里云杭州", WorkerCount: 12, RunningCount: 11, AbnormalCount: 1, OwnerName: "运维团队"},
			},
		},
		workers: &fakeWorkers{
			list: []worker.Worker{
				{ID: 2, WorkerName: "worker71kah6ac07o", WorkerState: "ONLINE", WorkerType: "VM", ClusterID: 1, PrivateIP: "198.18.0.1", HealthLevel: "健康", WorkerLoad: 3.43},
			},
		},
		consoleJobs: &fakeConsoleJobs{
			job: consolejob.Job{
				ID:             21,
				Label:          "WORKER_INSTALL",
				TaskState:      "RUNNING",
				JobToken:       "token-abc",
				Launcher:       "admin",
				DataJobName:    "同步任务A",
				DataJobDesc:    "跨区域 nightly sync",
				WorkerName:     "worker71kah6ac07o",
				WorkerDesc:     "杭州工作节点A",
				DsInstanceID:   "cc-mysql-1",
				DatasourceDesc: "杭州主库-同步源",
				ResourceType:   "WORKER",
				ResourceID:     2,
				TaskVOList: []consolejob.Task{
					{ID: 31, TaskState: "RUNNING", StepName: "安装 Agent", Host: "198.18.0.1", ExecuteOrder: 1, Cancelable: true},
				},
			},
		},
		jobConfigs: &fakeJobConfigs{
			specs: []jobconfig.Spec{
				{ID: 1, SpecKind: "SYNC", SpecKindCN: "同步规格", Spec: "STANDARD", FullMemoryMB: 2048, IncreMemoryMB: 1024, CheckMemoryMB: 512},
			},
		},
		schemas: &fakeSchemas{
			items: []ccschema.ApiTransferObjIndexDO{
				{DataJobID: 11, DataJobName: "同步任务A", SrcFullTransferObjName: "源库.orders", DstFullTransferObjName: "目标库.orders", SrcDsType: "MYSQL", DstDsType: "STARROCKS"},
			},
		},
	}
}

func assertTableBlocksAligned(t *testing.T, output string) {
	t.Helper()

	lines := strings.Split(output, "\n")
	var block [][]int
	for _, line := range lines {
		positions := pipeDisplayPositions(line)
		if len(positions) == 0 {
			if len(block) > 0 {
				assertPipeBlock(t, block, output)
				block = nil
			}
			continue
		}
		block = append(block, positions)
	}
	if len(block) > 0 {
		assertPipeBlock(t, block, output)
	}
}

func assertPipeBlock(t *testing.T, block [][]int, output string) {
	t.Helper()

	if len(block) < 2 {
		return
	}
	baseline := block[0]
	for _, positions := range block[1:] {
		if len(positions) != len(baseline) {
			t.Fatalf("pipe count mismatch: got %v want %v in %q", positions, baseline, output)
		}
		for i := range positions {
			if positions[i] != baseline[i] {
				t.Fatalf("pipe position mismatch: got %v want %v in %q", positions, baseline, output)
			}
		}
	}
}

func assertLabelColumnAligned(t *testing.T, output string) {
	t.Helper()

	lines := strings.Split(output, "\n")
	var baseline int
	found := false
	for _, line := range lines {
		index := strings.Index(line, " : ")
		if index <= 0 {
			continue
		}
		position := runewidth.StringWidth(line[:index])
		if !found {
			baseline = position
			found = true
			continue
		}
		if position != baseline {
			t.Fatalf("label column mismatch: got %d want %d in %q", position, baseline, output)
		}
	}
	if !found {
		t.Fatalf("no aligned label lines found in %q", output)
	}
}

func pipeDisplayPositions(line string) []int {
	positions := make([]int, 0, strings.Count(line, "|"))
	width := 0
	for _, r := range line {
		if r == '|' {
			positions = append(positions, width)
		}
		width += runewidth.RuneWidth(r)
	}
	return positions
}
