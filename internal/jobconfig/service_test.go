package jobconfig

import (
	"cloudcanal-openapi-cli/internal/config"
	"cloudcanal-openapi-cli/internal/openapi"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServiceListsSpecs(t *testing.T) {
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}
		_, _ = w.Write([]byte(`{"code":"1","data":[{"id":1,"specKind":"SYNC","specKindCn":"同步","spec":"STANDARD","fullMemoryMb":2048,"increMemoryMb":1024,"checkMemoryMb":512}]}`))
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

	initialSync := true
	shortTermSync := false
	service := NewService(client)
	specs, err := service.ListSpecs(ListSpecsOptions{
		DataJobType:   "SYNC",
		InitialSync:   &initialSync,
		ShortTermSync: &shortTermSync,
	})
	if err != nil {
		t.Fatalf("ListSpecs() error = %v", err)
	}
	if len(specs) != 1 || specs[0].ID != 1 || specs[0].Spec != "STANDARD" {
		t.Fatalf("specs = %#v, want single spec", specs)
	}
	if gotBody["dataJobType"] != "SYNC" || gotBody["initialSync"] != true || gotBody["shortTermSync"] != false {
		t.Fatalf("request body = %#v, want filter body", gotBody)
	}
}
