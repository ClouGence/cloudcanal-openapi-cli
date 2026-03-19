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
