package datajob_test

import (
	"cloudcanal-openapi-cli/internal/config"
	"cloudcanal-openapi-cli/internal/datajob"
	"cloudcanal-openapi-cli/internal/openapi"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServiceListsJobs(t *testing.T) {
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}
		_, _ = w.Write([]byte(`{"code":"1","data":[{"dataJobId":11,"dataJobName":"sync-1","dataJobType":"SYNC","dataTaskState":"RUNNING"}]}`))
	}))
	defer server.Close()

	client, err := openapi.NewClient(config.AppConfig{
		APIBaseURL: server.URL,
		AccessKey:  "test-ak",
		SecretKey:  "test-sk",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	service := datajob.NewService(client)
	jobs, err := service.ListJobs(datajob.ListOptions{
		DataJobName:      "sync",
		DataJobType:      "SYNC",
		Desc:             "nightly",
		SourceInstanceID: 101,
		TargetInstanceID: 202,
	})
	if err != nil {
		t.Fatalf("ListJobs() error = %v", err)
	}
	if len(jobs) != 1 || jobs[0].DataJobID != 11 || jobs[0].DataJobName != "sync-1" {
		t.Fatalf("jobs = %#v, want single sync-1 job", jobs)
	}
	if gotBody["dataJobName"] != "sync" || gotBody["dataJobType"] != "SYNC" || gotBody["desc"] != "nightly" {
		t.Fatalf("request body = %#v, want list filters", gotBody)
	}
	if gotBody["sourceInstanceId"] != float64(101) || gotBody["targetInstanceId"] != float64(202) {
		t.Fatalf("request body = %#v, want instance filters", gotBody)
	}
}

func TestServiceGetsJobDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code":"1","data":{"dataJobId":11,"dataJobName":"sync-1","dataJobDesc":"nightly sync","dataJobType":"SYNC","dataTaskState":"RUNNING","currTaskStatus":"FULL_RUNNING","consoleJobId":21,"lifeCycleState":"ACTIVE","sourceDsVO":{"instanceDesc":"src-db","dataSourceType":"MYSQL"},"targetDsVO":{"instanceDesc":"dst-db","dataSourceType":"STARROCKS"},"dataTasks":[{"dataTaskId":101,"dataTaskName":"full-task","dataTaskStatus":"RUNNING"}]}}`))
	}))
	defer server.Close()

	client, err := openapi.NewClient(config.AppConfig{
		APIBaseURL: server.URL,
		AccessKey:  "test-ak",
		SecretKey:  "test-sk",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	service := datajob.NewService(client)
	job, err := service.GetJob(11)
	if err != nil {
		t.Fatalf("GetJob() error = %v", err)
	}
	if job.DataJobID != 11 || job.DataJobDesc != "nightly sync" || len(job.DataTasks) != 1 {
		t.Fatalf("job = %#v, want detailed job with one task", job)
	}
}

func TestServiceGetsJobSchema(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code":"1","data":{"sourceSchema":"src","targetSchema":"dst","mappingConfig":"{\"rules\":1}","defaultTopic":"topic-a","defaultTopicPartition":8,"schemaWhiteListLevel":"TABLE","srcSchemaLessFormat":"JSON","dstSchemaLessFormat":"AVRO"}}`))
	}))
	defer server.Close()

	client, err := openapi.NewClient(config.AppConfig{
		APIBaseURL: server.URL,
		AccessKey:  "test-ak",
		SecretKey:  "test-sk",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	service := datajob.NewService(client)
	schema, err := service.GetJobSchema(11)
	if err != nil {
		t.Fatalf("GetJobSchema() error = %v", err)
	}
	if schema.SourceSchema != "src" || schema.TargetSchema != "dst" || schema.DefaultTopicPartition != 8 {
		t.Fatalf("schema = %#v, want detailed schema", schema)
	}
}

func TestServiceReplayJobSendsFlags(t *testing.T) {
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}
		_, _ = w.Write([]byte(`{"code":"1","msg":"ok"}`))
	}))
	defer server.Close()

	client, err := openapi.NewClient(config.AppConfig{
		APIBaseURL: server.URL,
		AccessKey:  "test-ak",
		SecretKey:  "test-sk",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	service := datajob.NewService(client)
	if err := service.ReplayJob(12, datajob.ReplayOptions{AutoStart: true, ResetToCreated: true}); err != nil {
		t.Fatalf("ReplayJob() error = %v", err)
	}
	if gotBody["jobId"] != float64(12) || gotBody["autoStart"] != true || gotBody["resetToCreated"] != true {
		t.Fatalf("request body = %#v, want replay flags", gotBody)
	}
}

func TestServiceRejectsBusinessFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code":"0","msg":"invalid credentials"}`))
	}))
	defer server.Close()

	client, err := openapi.NewClient(config.AppConfig{
		APIBaseURL: server.URL,
		AccessKey:  "test-ak",
		SecretKey:  "test-sk",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	service := datajob.NewService(client)
	if err := service.StartJob(10); err == nil || err.Error() != "invalid credentials" {
		t.Fatalf("StartJob() error = %v, want invalid credentials", err)
	}
}

func TestServiceCreatesJobAndSendsFullPayload(t *testing.T) {
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}
		_, _ = w.Write([]byte(`{"code":"1","data":"9876"}`))
	}))
	defer server.Close()

	client, err := openapi.NewClient(config.AppConfig{
		APIBaseURL: server.URL,
		AccessKey:  "test-ak",
		SecretKey:  "test-sk",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	structMigration := true
	initialSync := true
	shortTermSync := false
	shortTermNum := 3
	filterDDL := true
	autoStart := true
	checkOnce := false
	checkPeriod := true
	fullPeriod := false
	specID := int64(18)
	dstTopicPartitions := 12
	dstDdlTopicPartitions := 2
	kuduNumReplicas := 3
	dbHeartbeatEnable := true

	service := datajob.NewService(client)
	result, err := service.CreateJob(datajob.CreateJobRequest{
		ClusterID:              1,
		SrcDsID:                391,
		DstDsID:                16,
		SrcHostType:            "PRIVATE",
		DstHostType:            "PUBLIC",
		SchemaWhiteListLevel:   "TABLE",
		SrcSchema:              `{"schema":"src"}`,
		DstSchema:              `{"schema":"dst"}`,
		MappingDef:             `[{"method":"DB_DB"}]`,
		SrcCaseSensitiveType:   "LOWER",
		DstCaseSensitiveType:   "UPPER",
		SrcDsCharset:           "utf8",
		TarDsCharset:           "utf8",
		KeyConflictStrategy:    "IGNORE",
		JobType:                "SYNC",
		DataJobDesc:            "api created",
		StructMigration:        &structMigration,
		InitialSync:            &initialSync,
		ShortTermSync:          &shortTermSync,
		ShortTermNum:           &shortTermNum,
		FilterDDL:              &filterDDL,
		SpecID:                 &specID,
		AutoStart:              &autoStart,
		CheckOnce:              &checkOnce,
		CheckPeriod:            &checkPeriod,
		CheckPeriodCronExpr:    "0 0 * * * ?",
		FullPeriod:             &fullPeriod,
		FullPeriodCronExpr:     "0 0 1 * * ?",
		DstMqDefaultTopic:      "topic-a",
		DstMqDefaultTopicParts: &dstTopicPartitions,
		DstMqDdlTopic:          "ddl-topic",
		DstMqDdlTopicParts:     &dstDdlTopicPartitions,
		SrcSchemaLessFormat:    "JSON",
		DstSchemaLessFormat:    "AVRO",
		OriginDecodeMsgFormat:  "DEFAULT",
		DstCkTableEngine:       "MergeTree",
		DstSrOrDorisTableModel: "DUPLICATE",
		KafkaConsumerGroupID:   "cc-group",
		KuduNumReplicas:        &kuduNumReplicas,
		SrcRocketMqGroupID:     "rocket-src",
		SrcRabbitMqVhost:       "/src",
		SrcRabbitExchange:      "src-ex",
		DstRabbitMqVhost:       "/dst",
		DstRabbitExchange:      "dst-ex",
		ObTenant:               "tenant-a",
		DbHeartbeatEnable:      &dbHeartbeatEnable,
	})
	if err != nil {
		t.Fatalf("CreateJob() error = %v", err)
	}
	if result.Data != "9876" || result.JobID != "9876" {
		t.Fatalf("result = %#v, want job id 9876", result)
	}
	if gotBody["clusterId"] != float64(1) || gotBody["srcDsId"] != float64(391) || gotBody["dstDsId"] != float64(16) {
		t.Fatalf("request body = %#v, want core create fields", gotBody)
	}
	if gotBody["dataJobDesc"] != "api created" || gotBody["jobType"] != "SYNC" || gotBody["autoStart"] != true {
		t.Fatalf("request body = %#v, want create flags", gotBody)
	}
	if gotBody["dstMqDefaultTopicPartitions"] != float64(12) || gotBody["dbHeartbeatEnable"] != true {
		t.Fatalf("request body = %#v, want advanced fields", gotBody)
	}
}

func TestServiceUpdatesIncrePosAndHandlesFailures(t *testing.T) {
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}
		_, _ = w.Write([]byte(`{"code":"1","data":"updated"}`))
	}))
	defer server.Close()

	client, err := openapi.NewClient(config.AppConfig{
		APIBaseURL: server.URL,
		AccessKey:  "test-ak",
		SecretKey:  "test-sk",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	filePosition := int64(794891)
	positionTimestamp := int64(1710000000)
	serverID := int64(88)
	scn := int64(99)
	scnIndex := int64(5)
	dataID := int64(1234)
	transactionID := int64(5678)

	service := datajob.NewService(client)
	result, err := service.UpdateIncrePos(datajob.UpdateIncrePosRequest{
		TaskID:            320,
		PosType:           "MYSQL_LOG_FILE_POS",
		JournalFile:       "binlog.000491",
		FilePosition:      &filePosition,
		GtidPosition:      "gtid-set",
		PositionTimestamp: &positionTimestamp,
		ServerID:          &serverID,
		Lsn:               "lsn-1",
		Scn:               &scn,
		ScnIndex:          &scnIndex,
		CommonPosStr:      "common-pos",
		DataID:            &dataID,
		TransactionID:     &transactionID,
	})
	if err != nil {
		t.Fatalf("UpdateIncrePos() error = %v", err)
	}
	if result.Data != "updated" {
		t.Fatalf("result = %#v, want updated", result)
	}
	if gotBody["taskId"] != float64(320) || gotBody["posType"] != "MYSQL_LOG_FILE_POS" || gotBody["journalFile"] != "binlog.000491" {
		t.Fatalf("request body = %#v, want position fields", gotBody)
	}
	if gotBody["filePosition"] != float64(794891) || gotBody["serverId"] != float64(88) || gotBody["transactionId"] != float64(5678) {
		t.Fatalf("request body = %#v, want numeric fields", gotBody)
	}

	failureServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code":"0","msg":"update failed"}`))
	}))
	defer failureServer.Close()

	failureClient, err := openapi.NewClient(config.AppConfig{
		APIBaseURL: failureServer.URL,
		AccessKey:  "test-ak",
		SecretKey:  "test-sk",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	failureService := datajob.NewService(failureClient)
	if _, err := failureService.UpdateIncrePos(datajob.UpdateIncrePosRequest{TaskID: 1, PosType: "MYSQL_LOG_FILE_POS"}); err == nil || err.Error() != "update failed" {
		t.Fatalf("UpdateIncrePos() error = %v, want update failed", err)
	}
}

func TestServiceAttachesAndDetachesIncrementTask(t *testing.T) {
	var bodies []map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}
		bodies = append(bodies, body)
		_, _ = w.Write([]byte(`{"code":"1","msg":"ok"}`))
	}))
	defer server.Close()

	client, err := openapi.NewClient(config.AppConfig{
		APIBaseURL: server.URL,
		AccessKey:  "test-ak",
		SecretKey:  "test-sk",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	service := datajob.NewService(client)
	if err := service.AttachIncreJob(11); err != nil {
		t.Fatalf("AttachIncreJob() error = %v", err)
	}
	if err := service.DetachIncreJob(11); err != nil {
		t.Fatalf("DetachIncreJob() error = %v", err)
	}
	if len(bodies) != 2 || bodies[0]["jobId"] != float64(11) || bodies[1]["jobId"] != float64(11) {
		t.Fatalf("bodies = %#v, want jobId 11 for both calls", bodies)
	}
}
