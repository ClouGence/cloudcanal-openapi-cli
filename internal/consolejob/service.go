package consolejob

import "cloudcanal-openapi-cli/internal/openapi"

const queryPath = "/cloudcanal/console/api/v1/openapi/consolejob/queryconsolejob"

type Operations interface {
	Get(consoleJobID int64) (Job, error)
}

type Service struct {
	client *openapi.Client
}

func NewService(client *openapi.Client) *Service {
	return &Service{client: client}
}

type queryRequest struct {
	ConsoleJobID int64 `json:"consoleJobId"`
}

type Job struct {
	ID             int64  `json:"id"`
	JobToken       string `json:"jobToken"`
	Label          string `json:"label"`
	DataJobName    string `json:"dataJobName"`
	DataJobDesc    string `json:"dataJobDesc"`
	WorkerName     string `json:"workerName"`
	WorkerDesc     string `json:"workerDesc"`
	DsInstanceID   string `json:"dsInstanceId"`
	DatasourceDesc string `json:"datasourceDesc"`
	TaskState      string `json:"taskState"`
	Launcher       string `json:"launcher"`
	ResourceType   string `json:"resourceType"`
	ResourceID     int64  `json:"resourceId"`
	TaskVOList     []Task `json:"taskVOList"`
}

type Task struct {
	ID               int64  `json:"id"`
	JobID            int64  `json:"jobId"`
	TaskState        string `json:"taskState"`
	HandlerBeanName  string `json:"handlerBeanName"`
	HandlerClassName string `json:"handlerClassName"`
	Host             string `json:"host"`
	ExecuteOrder     int    `json:"executeOrder"`
	Message          string `json:"message"`
	Cancelable       bool   `json:"cancelable"`
	StepName         string `json:"stepName"`
}

type queryResponse struct {
	openapi.Response
	Data Job `json:"data"`
}

func (s *Service) Get(consoleJobID int64) (Job, error) {
	var out queryResponse
	if err := s.client.PostJSONWithOptions(queryPath, queryRequest{ConsoleJobID: consoleJobID}, &out, openapi.RequestOptions{Retryable: true}); err != nil {
		return Job{}, err
	}
	if err := openapi.EnsureSuccess(out.Response, "failed to query console job"); err != nil {
		return Job{}, err
	}
	return out.Data, nil
}
