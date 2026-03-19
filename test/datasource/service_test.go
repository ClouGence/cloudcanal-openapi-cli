package datasource_test

import (
	"cloudcanal-openapi-cli/internal/config"
	"cloudcanal-openapi-cli/internal/datasource"
	"cloudcanal-openapi-cli/internal/openapi"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestServiceListsDataSources(t *testing.T) {
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}
		_, _ = w.Write([]byte(`{"code":"1","data":[{"id":7,"instanceId":"cc-mysql-1","dataSourceType":"MYSQL","hostType":"RDS","deployType":"ALIYUN","lifeCycleState":"ACTIVE","instanceDesc":"mysql source","consoleJobId":123}]}`))
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
	if string(sources[0].ConsoleJobID) != "123" {
		t.Fatalf("consoleJobId = %q, want 123", sources[0].ConsoleJobID)
	}
	if gotBody["type"] != "MYSQL" || gotBody["deployType"] != "ALIYUN" || gotBody["hostType"] != "RDS" || gotBody["lifeCycleState"] != "ACTIVE" {
		t.Fatalf("request body = %#v, want filters", gotBody)
	}
}

func TestServiceGetsDataSourceByID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code":"1","data":[{"id":9,"instanceId":"cc-sr-1","dataSourceType":"STARROCKS","consoleJobId":"456"}]}`))
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
	if string(source.ConsoleJobID) != "456" {
		t.Fatalf("consoleJobId = %q, want 456", source.ConsoleJobID)
	}
}

func TestServiceAddsDataSource(t *testing.T) {
	securityFile := writeTempFile(t, "security.pem", "security-content")
	secretFile := writeTempFile(t, "secret.key", "secret-content")

	var gotBody datasource.ApiDsAddData
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if got := r.URL.Query().Get("AccessKeyId"); got != "test-ak" {
			t.Fatalf("AccessKeyId = %q, want test-ak", got)
		}
		if got := r.URL.Query().Get("SignatureMethod"); got != "HmacSHA1" {
			t.Fatalf("SignatureMethod = %q, want HmacSHA1", got)
		}
		if got := r.URL.Query().Get("SignatureNonce"); got == "" {
			t.Fatalf("SignatureNonce is empty")
		}
		if got := r.URL.Query().Get("Signature"); got == "" {
			t.Fatalf("Signature is empty")
		}

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			t.Fatalf("ParseMultipartForm() error = %v", err)
		}
		if err := json.Unmarshal([]byte(r.FormValue("dataSourceAddData")), &gotBody); err != nil {
			t.Fatalf("unmarshal dataSourceAddData error = %v", err)
		}

		file, header, err := r.FormFile("securityFile")
		if err != nil {
			t.Fatalf("securityFile missing: %v", err)
		}
		defer file.Close()
		securityContent, _ := io.ReadAll(file)
		if header.Filename != filepath.Base(securityFile) || string(securityContent) != "security-content" {
			t.Fatalf("securityFile = %q/%q, want %q/security-content", header.Filename, string(securityContent), filepath.Base(securityFile))
		}

		file, header, err = r.FormFile("secretFile")
		if err != nil {
			t.Fatalf("secretFile missing: %v", err)
		}
		defer file.Close()
		secretContent, _ := io.ReadAll(file)
		if header.Filename != filepath.Base(secretFile) || string(secretContent) != "secret-content" {
			t.Fatalf("secretFile = %q/%q, want %q/secret-content", header.Filename, string(secretContent), filepath.Base(secretFile))
		}

		_, _ = w.Write([]byte(`{"code":"1","data":"ds-123"}`))
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
	data, err := service.Add(datasource.AddOptions{
		DataSourceAddData: datasource.ApiDsAddData{
			DeployType:        "ALIYUN",
			Region:            "cn-hangzhou",
			Type:              "MYSQL",
			Host:              "127.0.0.1:3306",
			PrivateHost:       "127.0.0.1:3306",
			HostType:          "PRIVATE",
			InstanceDesc:      "mysql source",
			InstanceID:        "mysql-1",
			AutoCreateAccount: true,
			Account:           "root",
			Password:          "secret",
			SecurityType:      "USER_PASSWD",
			ClusterIDs:        []int64{11, 12},
			LifeCycleState:    "ACTIVE",
			Version:           "8.0",
			Driver:            "com.mysql.cj.jdbc.Driver",
			ConnectType:       "DIRECT",
			DsKvConfigs: []datasource.KvBaseConfig{
				{ConfigName: "k1", ConfigValue: "v1"},
			},
		},
		SecurityFilePath: securityFile,
		SecretFilePath:   secretFile,
	})
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	if data != "ds-123" {
		t.Fatalf("data = %q, want ds-123", data)
	}
	if gotBody.Type != "MYSQL" || gotBody.Host != "127.0.0.1:3306" || !gotBody.AutoCreateAccount {
		t.Fatalf("multipart body = %#v, want datasource add payload", gotBody)
	}
	if len(gotBody.ClusterIDs) != 2 || gotBody.ClusterIDs[0] != 11 || gotBody.DsKvConfigs[0].ConfigName != "k1" {
		t.Fatalf("multipart body = %#v, want cluster and kv config fields", gotBody)
	}
}

func TestServiceRejectsAddFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code":"0","msg":"create failed"}`))
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
	if _, err := service.Add(datasource.AddOptions{DataSourceAddData: datasource.ApiDsAddData{Type: "MYSQL"}}); err == nil || err.Error() != "create failed" {
		t.Fatalf("Add() error = %v, want create failed", err)
	}
}

func TestServiceDeletesDataSource(t *testing.T) {
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

	service := datasource.NewService(client)
	if err := service.Delete(99); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if gotBody["dataSourceId"] != float64(99) {
		t.Fatalf("request body = %#v, want dataSourceId 99", gotBody)
	}
}

func TestServiceRejectsDeleteFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code":"0","msg":"delete failed"}`))
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
	if err := service.Delete(99); err == nil || err.Error() != "delete failed" {
		t.Fatalf("Delete() error = %v, want delete failed", err)
	}
}

func writeTempFile(t *testing.T, name string, content string) string {
	t.Helper()
	file, err := os.CreateTemp(t.TempDir(), name)
	if err != nil {
		t.Fatalf("CreateTemp() error = %v", err)
	}
	if _, err := file.WriteString(content); err != nil {
		t.Fatalf("WriteString() error = %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	return file.Name()
}
