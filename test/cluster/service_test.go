package cluster_test

import (
	"cloudcanal-openapi-cli/internal/cluster"
	"cloudcanal-openapi-cli/internal/config"
	"cloudcanal-openapi-cli/internal/openapi"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServiceListsClusters(t *testing.T) {
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}
		_, _ = w.Write([]byte(`{"code":"1","data":[{"id":3,"clusterName":"prod-cluster","region":"cn-hangzhou","cloudOrIdcName":"ALIYUN","workerCount":5,"runningCount":4,"abnormalCount":1,"ownerName":"admin"}]}`))
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

	service := cluster.NewService(client)
	clusters, err := service.List(cluster.ListOptions{
		ClusterName:    "prod",
		ClusterDesc:    "main",
		CloudOrIDCName: "ALIYUN",
		Region:         "cn-hangzhou",
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(clusters) != 1 || clusters[0].ID != 3 || clusters[0].ClusterName != "prod-cluster" {
		t.Fatalf("clusters = %#v, want single cluster", clusters)
	}
	if gotBody["clusterNameLike"] != "prod" || gotBody["clusterDescLike"] != "main" || gotBody["cloudOrIdcName"] != "ALIYUN" || gotBody["region"] != "cn-hangzhou" {
		t.Fatalf("request body = %#v, want filters", gotBody)
	}
}
