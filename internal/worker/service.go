package worker

import "cloudcanal-openapi-cli/internal/openapi"

const (
	listPath  = "/cloudcanal/console/api/v1/openapi/worker/listworkers"
	startPath = "/cloudcanal/console/api/v1/openapi/worker/startWorker"
	stopPath  = "/cloudcanal/console/api/v1/openapi/worker/stopWorker"
)

type Operations interface {
	List(options ListOptions) ([]Worker, error)
	Start(workerID int64) error
	Stop(workerID int64) error
}

type Service struct {
	client *openapi.Client
}

func NewService(client *openapi.Client) *Service {
	return &Service{client: client}
}

type ListOptions struct {
	ClusterID        int64
	SourceInstanceID int64
	TargetInstanceID int64
}

type listRequest struct {
	ClusterID        *int64 `json:"clusterId,omitempty"`
	SourceInstanceID *int64 `json:"sourceInstanceId,omitempty"`
	TargetInstanceID *int64 `json:"targetInstanceId,omitempty"`
}

type actionRequest struct {
	WorkerID int64 `json:"workerId"`
}

type Worker struct {
	ID               int64   `json:"id"`
	ClusterID        int64   `json:"clusterId"`
	PrivateIP        string  `json:"privateIp"`
	PublicIP         string  `json:"publicIp"`
	CloudOrIDCName   string  `json:"cloudOrIdcName"`
	Region           string  `json:"region"`
	WorkerType       string  `json:"workerType"`
	WorkerState      string  `json:"workerState"`
	HealthLevel      string  `json:"healthLevel"`
	WorkerLoad       float64 `json:"workerLoad"`
	WorkerName       string  `json:"workerName"`
	WorkerSeqNumber  string  `json:"workerSeqNumber"`
	WorkerDesc       string  `json:"workerDesc"`
	ConsoleJobID     int64   `json:"consoleJobId"`
	ConsoleTaskState string  `json:"consoleTaskState"`
}

type listResponse struct {
	openapi.Response
	Data []Worker `json:"data"`
}

func (s *Service) List(options ListOptions) ([]Worker, error) {
	var out listResponse
	if err := s.client.PostJSON(listPath, newListRequest(options), &out); err != nil {
		return nil, err
	}
	if err := openapi.EnsureSuccess(out.Response, "failed to list workers"); err != nil {
		return nil, err
	}
	if out.Data == nil {
		return []Worker{}, nil
	}
	return out.Data, nil
}

func (s *Service) Start(workerID int64) error {
	return s.doAction(startPath, workerID, "failed to start worker")
}

func (s *Service) Stop(workerID int64) error {
	return s.doAction(stopPath, workerID, "failed to stop worker")
}

func (s *Service) doAction(path string, workerID int64, fallback string) error {
	var out openapi.Response
	if err := s.client.PostJSON(path, actionRequest{WorkerID: workerID}, &out); err != nil {
		return err
	}
	return openapi.EnsureSuccess(out, fallback)
}

func newListRequest(options ListOptions) listRequest {
	req := listRequest{}
	if options.ClusterID > 0 {
		req.ClusterID = &options.ClusterID
	}
	if options.SourceInstanceID > 0 {
		req.SourceInstanceID = &options.SourceInstanceID
	}
	if options.TargetInstanceID > 0 {
		req.TargetInstanceID = &options.TargetInstanceID
	}
	return req
}
