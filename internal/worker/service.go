package worker

import "cloudcanal-openapi-cli/internal/openapi"

const (
	listPath              = "/cloudcanal/console/api/v1/openapi/worker/listworkers"
	startPath             = "/cloudcanal/console/api/v1/openapi/worker/startWorker"
	stopPath              = "/cloudcanal/console/api/v1/openapi/worker/stopWorker"
	deletePath            = "/cloudcanal/console/api/v1/openapi/worker/deleteWorker"
	modifyMemOverSoldPath = "/cloudcanal/console/api/v1/openapi/worker/modifyMemOverSoldPercent"
	updateAlertPath       = "/cloudcanal/console/api/v1/openapi/worker/updateWorkerAlertConfig"
)

type Operations interface {
	List(options ListOptions) ([]Worker, error)
	Start(workerID int64) error
	Stop(workerID int64) error
	Delete(workerID int64) error
	ModifyMemOverSold(workerID int64, memOverSoldPercent int) error
	UpdateWorkerAlert(workerID int64, phone, email, im, sms bool) error
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

type workerActionRequest struct {
	WorkerID int64 `json:"workerId"`
}

type modifyMemOverSoldRequest struct {
	WorkerID           int64 `json:"workerId"`
	MemOverSoldPercent int   `json:"memOverSoldPercent"`
}

type updateWorkerAlertRequest struct {
	WorkerID int64 `json:"workerId"`
	Phone    bool  `json:"phone"`
	Email    bool  `json:"email"`
	Im       bool  `json:"im"`
	Sms      bool  `json:"sms"`
}

type Worker struct {
	ID                    int64   `json:"id"`
	ClusterID             int64   `json:"clusterId"`
	PrivateIP             string  `json:"privateIp"`
	PublicIP              string  `json:"publicIp"`
	CloudOrIDCName        string  `json:"cloudOrIdcName"`
	Region                string  `json:"region"`
	TotalTaskMemMB        int64   `json:"totalTaskMemMb"`
	MemOverSoldPercent    int     `json:"memOverSoldPercent"`
	PhysicMemMB           int64   `json:"physicMemMb"`
	PhysicCoreNum         int     `json:"physicCoreNum"`
	LogicalCoreNum        int     `json:"logicalCoreNum"`
	PhysicDiskGB          int64   `json:"physicDiskGb"`
	WorkerType            string  `json:"workerType"`
	WorkerState           string  `json:"workerState"`
	CPUUseRatio           float64 `json:"cpuUseRatio"`
	MemUseRatio           float64 `json:"memUseRatio"`
	HealthLevel           string  `json:"healthLevel"`
	TaskHeapSizeMB        int64   `json:"taskHeapSizeMb"`
	FreeMemMB             int64   `json:"freeMemMb"`
	FreeDiskGB            int64   `json:"freeDiskGb"`
	WorkerLoad            float64 `json:"workerLoad"`
	WorkerName            string  `json:"workerName"`
	WorkerSeqNumber       string  `json:"workerSeqNumber"`
	WorkerDesc            string  `json:"workerDesc"`
	InstallConsoleJobID   int64   `json:"installConsoleJobId"`
	UninstallConsoleJobID int64   `json:"uninstallConsoleJobId"`
	DeployStatus          string  `json:"deployStatus"`
	ConsoleJobID          int64   `json:"consoleJobId"`
	ConsoleTaskState      string  `json:"consoleTaskState"`
}

type listResponse struct {
	openapi.Response
	Data []Worker `json:"data"`
}

func (s *Service) List(options ListOptions) ([]Worker, error) {
	var out listResponse
	if err := s.client.PostJSONWithOptions(listPath, newListRequest(options), &out, openapi.RequestOptions{Retryable: true}); err != nil {
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

func (s *Service) Delete(workerID int64) error {
	return s.doAction(deletePath, workerID, "failed to delete worker")
}

func (s *Service) ModifyMemOverSold(workerID int64, memOverSoldPercent int) error {
	var out openapi.Response
	req := modifyMemOverSoldRequest{WorkerID: workerID, MemOverSoldPercent: memOverSoldPercent}
	if err := s.client.PostJSON(modifyMemOverSoldPath, req, &out); err != nil {
		return err
	}
	return openapi.EnsureSuccess(out, "failed to modify worker mem oversold percent")
}

func (s *Service) UpdateWorkerAlert(workerID int64, phone, email, im, sms bool) error {
	var out openapi.Response
	req := updateWorkerAlertRequest{
		WorkerID: workerID,
		Phone:    phone,
		Email:    email,
		Im:       im,
		Sms:      sms,
	}
	if err := s.client.PostJSON(updateAlertPath, req, &out); err != nil {
		return err
	}
	return openapi.EnsureSuccess(out, "failed to update worker alert config")
}

func (s *Service) doAction(path string, workerID int64, fallback string) error {
	var out openapi.Response
	if err := s.client.PostJSON(path, workerActionRequest{WorkerID: workerID}, &out); err != nil {
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
