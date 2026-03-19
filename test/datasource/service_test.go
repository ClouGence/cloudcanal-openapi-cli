package datasource_test

import (
	"cloudcanal-openapi-cli/internal/config"
	"cloudcanal-openapi-cli/internal/datasource"
	"cloudcanal-openapi-cli/internal/openapi"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServiceListsDataSources(t *testing.T) {
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}
		_, _ = w.Write([]byte(`{"code":"1","data":[{"id":7,"instanceId":"cc-mysql-1","dataSourceType":"MYSQL","hostType":"RDS","deployType":"ALIYUN","lifeCycleState":"ACTIVE","instanceDesc":"mysql source"}]}`))
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

	service := datasource.NewService(client)
	sources, err := service.List(datasource.ListOptions{Type: "MYSQL", DeployType: "ALIYUN", HostType: "RDS", LifeCycleState: "ACTIVE"})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(sources) != 1 || sources[0].ID != 7 || sources[0].InstanceID != "cc-mysql-1" {
		t.Fatalf("sources = %#v, want single datasource", sources)
	}
	if gotBody["type"] != "MYSQL" || gotBody["deployType"] != "ALIYUN" || gotBody["hostType"] != "RDS" || gotBody["lifeCycleState"] != "ACTIVE" {
		t.Fatalf("request body = %#v, want filters", gotBody)
	}
}

func TestServiceGetsDataSourceByID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code":"1","data":[{"id":9,"instanceId":"cc-sr-1","dataSourceType":"STARROCKS"}]}`))
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

	service := datasource.NewService(client)
	source, err := service.Get(9)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if source.ID != 9 || source.DataSourceType != "STARROCKS" {
		t.Fatalf("source = %#v, want id 9", source)
	}
}
