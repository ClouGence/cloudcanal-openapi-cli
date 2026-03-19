package schema_test

import (
	"cloudcanal-openapi-cli/internal/config"
	"cloudcanal-openapi-cli/internal/openapi"
	"cloudcanal-openapi-cli/internal/schema"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServiceListsTransferObjectsByMeta(t *testing.T) {
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}
		_, _ = w.Write([]byte(`{"code":"1","data":[{"id":1,"srcTransferObjName":"src_table","srcDsInstanceId":11,"srcDsInstanceName":"src-ds","srcFullTransferObjName":"src_db.src_table","srcDsType":"MYSQL","srcDb":"src_db","srcSchema":"public","filterExpr":"id > 0","specifiedPks":"id","dstTransferObjName":"dst_table","dstDsInstanceId":22,"dstDsInstanceName":"dst-ds","dstFullTransferObjName":"dst_db.dst_table","dstDsType":"STARROCKS","dstDb":"dst_db","dstSchema":"public","dataJobId":33,"dataJobName":"job-1","dataJobDesc":"sync job"}]}`))
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

	service := schema.NewService(client)
	items, err := service.ListTransObjsByMeta(schema.ListTransObjsByMetaOptions{
		SrcDb:       "src_db",
		SrcSchema:   "public",
		SrcTransObj: "src_table",
		DstDb:       "dst_db",
		DstSchema:   "public",
		DstTranObj:  "dst_table",
	})
	if err != nil {
		t.Fatalf("ListTransObjsByMeta() error = %v", err)
	}
	if len(items) != 1 || items[0].ID != 1 || items[0].DataJobName != "job-1" {
		t.Fatalf("items = %#v, want single transfer object", items)
	}
	if gotBody["srcDb"] != "src_db" || gotBody["dstTranObj"] != "dst_table" {
		t.Fatalf("request body = %#v, want meta filters", gotBody)
	}
}

func TestServiceRejectsTransferObjectsByMetaFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code":"0","msg":"meta query failed"}`))
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

	service := schema.NewService(client)
	if _, err := service.ListTransObjsByMeta(schema.ListTransObjsByMetaOptions{SrcDb: "src_db"}); err == nil || err.Error() != "meta query failed" {
		t.Fatalf("ListTransObjsByMeta() error = %v, want meta query failed", err)
	}
}
