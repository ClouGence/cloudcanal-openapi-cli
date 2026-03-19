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
