package worker_test

import (
	"cloudcanal-openapi-cli/internal/config"
	"cloudcanal-openapi-cli/internal/openapi"
	"cloudcanal-openapi-cli/internal/worker"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServiceListsWorkers(t *testing.T) {
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}
		_, _ = w.Write([]byte(`{"code":"1","data":[{"id":5,"clusterId":2,"workerName":"worker-1","workerState":"RUNNING","workerType":"FULL","privateIp":"10.0.0.5","healthLevel":"GREEN","workerLoad":0.8}]}`))
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

	service := worker.NewService(client)
	workers, err := service.List(worker.ListOptions{ClusterID: 2, SourceInstanceID: 1001, TargetInstanceID: 1002})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(workers) != 1 || workers[0].ID != 5 || workers[0].WorkerName != "worker-1" {
		t.Fatalf("workers = %#v, want single worker", workers)
	}
	if gotBody["clusterId"] != float64(2) || gotBody["sourceInstanceId"] != float64(1001) || gotBody["targetInstanceId"] != float64(1002) {
		t.Fatalf("request body = %#v, want list filters", gotBody)
	}
}

func TestServiceStartsAndStopsWorker(t *testing.T) {
	var requests []map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}
		requests = append(requests, body)
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

	service := worker.NewService(client)
	if err := service.Start(6); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if err := service.Stop(6); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}
	if len(requests) != 2 || requests[0]["workerId"] != float64(6) || requests[1]["workerId"] != float64(6) {
		t.Fatalf("requests = %#v, want workerId 6 for both calls", requests)
	}
}

func TestServiceSupportsWorkerActions(t *testing.T) {
	type recorded struct {
		path string
		body map[string]any
	}

	var requests []recorded
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}
		requests = append(requests, recorded{path: r.URL.Path, body: body})
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

	service := worker.NewService(client)
	if err := service.Delete(7); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if err := service.ModifyMemOverSold(8, 120); err != nil {
		t.Fatalf("ModifyMemOverSold() error = %v", err)
	}
	if err := service.UpdateWorkerAlert(9, true, false, true, false); err != nil {
		t.Fatalf("UpdateWorkerAlert() error = %v", err)
	}

	if len(requests) != 3 {
		t.Fatalf("requests = %#v, want 3 calls", requests)
	}
	if requests[0].path != "/cloudcanal/console/api/v1/openapi/worker/deleteWorker" || requests[0].body["workerId"] != float64(7) {
		t.Fatalf("delete request = %#v, want worker delete payload", requests[0])
	}
	if requests[1].path != "/cloudcanal/console/api/v1/openapi/worker/modifyMemOverSoldPercent" || requests[1].body["workerId"] != float64(8) || requests[1].body["memOverSoldPercent"] != float64(120) {
		t.Fatalf("modify oversold request = %#v, want worker oversold payload", requests[1])
	}
	if requests[2].path != "/cloudcanal/console/api/v1/openapi/worker/updateWorkerAlertConfig" {
		t.Fatalf("alert request path = %q, want updateWorkerAlertConfig", requests[2].path)
	}
	if requests[2].body["workerId"] != float64(9) || requests[2].body["phone"] != true || requests[2].body["email"] != false || requests[2].body["im"] != true || requests[2].body["sms"] != false {
		t.Fatalf("alert request body = %#v, want alert config payload", requests[2].body)
	}
}

func TestServicePropagatesWorkerBusinessFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code":"0","msg":"worker not found"}`))
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

	service := worker.NewService(client)
	if err := service.Delete(10); err == nil || err.Error() != "worker not found" {
		t.Fatalf("Delete() error = %v, want worker not found", err)
	}
}

func TestServiceDecodesExtendedWorkerFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code":"1","data":[{"id":5,"clusterId":2,"workerName":"worker-1","workerState":"RUNNING","workerType":"FULL","privateIp":"10.0.0.5","publicIp":"1.1.1.1","cloudOrIdcName":"aliyun","region":"cn-hangzhou","totalTaskMemMb":1024,"memOverSoldPercent":120,"physicMemMb":8192,"physicCoreNum":8,"logicalCoreNum":16,"physicDiskGb":500,"cpuUseRatio":0.3,"memUseRatio":0.5,"taskHeapSizeMb":256,"freeMemMb":4096,"freeDiskGb":200,"workerLoad":0.8,"deployStatus":"ONLINE","consoleJobId":99,"consoleTaskState":"DONE"}]}`))
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

	service := worker.NewService(client)
	workers, err := service.List(worker.ListOptions{ClusterID: 2})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(workers) != 1 || workers[0].MemOverSoldPercent != 120 || workers[0].PublicIP != "1.1.1.1" || workers[0].DeployStatus != "ONLINE" {
		t.Fatalf("workers = %#v, want extended worker fields", workers)
	}
}
