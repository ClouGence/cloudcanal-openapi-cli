package consolejob

import (
	"cloudcanal-openapi-cli/internal/config"
	"cloudcanal-openapi-cli/internal/openapi"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServiceGetsConsoleJob(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code":"1","data":{"id":21,"label":"WORKER_INSTALL","taskState":"RUNNING","jobToken":"abc","workerName":"worker-1","resourceType":"WORKER","resourceId":5,"taskVOList":[{"id":31,"taskState":"RUNNING","stepName":"Install","host":"10.0.0.5","executeOrder":1,"cancelable":true}]}}`))
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

	service := NewService(client)
	job, err := service.Get(21)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if job.ID != 21 || job.TaskState != "RUNNING" || len(job.TaskVOList) != 1 {
		t.Fatalf("job = %#v, want detailed console job", job)
	}
}
