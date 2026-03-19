package datajob

import (
	"cloudcanal-openapi-cli/internal/openapi"
)

const (
	listPath   = "/cloudcanal/console/api/v1/openapi/datajob/list"
	queryPath  = "/cloudcanal/console/api/v1/openapi/datajob/queryjob"
	schemaPath = "/cloudcanal/console/api/v1/openapi/datajob/queryjobschemabyid"
	startPath  = "/cloudcanal/console/api/v1/openapi/datajob/start"
	stopPath   = "/cloudcanal/console/api/v1/openapi/datajob/stop"
	deletePath = "/cloudcanal/console/api/v1/openapi/datajob/delete"
	replayPath = "/cloudcanal/console/api/v1/openapi/datajob/replay"
)

type Operations interface {
	ListJobs(options ListOptions) ([]Job, error)
	GetJob(jobID int64) (Job, error)
	GetJobSchema(jobID int64) (JobSchema, error)
	StartJob(jobID int64) error
	StopJob(jobID int64) error
	DeleteJob(jobID int64) error
	ReplayJob(jobID int64, options ReplayOptions) error
}

type Service struct {
	client *openapi.Client
}

func NewService(client *openapi.Client) *Service {
	return &Service{client: client}
}

type ListOptions struct {
	DataJobName      string
	DataJobType      string
	Desc             string
	SourceInstanceID int64
	TargetInstanceID int64
}

type listJobsRequest struct {
	DataJobName      string `json:"dataJobName,omitempty"`
	DataJobType      string `json:"dataJobType,omitempty"`
	Desc             string `json:"desc,omitempty"`
	SourceInstanceID *int64 `json:"sourceInstanceId,omitempty"`
	TargetInstanceID *int64 `json:"targetInstanceId,omitempty"`
}

type jobActionRequest struct {
	JobID int64 `json:"jobId"`
}

type replayJobRequest struct {
	JobID          int64 `json:"jobId"`
	AutoStart      *bool `json:"autoStart,omitempty"`
	ResetToCreated *bool `json:"resetToCreated,omitempty"`
}

type ReplayOptions struct {
	AutoStart      bool
	ResetToCreated bool
}

type Source struct {
	InstanceDesc   string `json:"instanceDesc"`
	InstanceID     string `json:"instanceId"`
	DataSourceType string `json:"dataSourceType"`
	HostType       string `json:"hostType"`
	DeployType     string `json:"deployType"`
	Region         string `json:"region"`
	LifeCycleState string `json:"lifeCycleState"`
}

type Task struct {
	DataTaskID     int64  `json:"dataTaskId"`
	DataTaskType   string `json:"dataTaskType"`
	DataTaskName   string `json:"dataTaskName"`
	DataTaskStatus string `json:"dataTaskStatus"`
	WorkerIP       string `json:"workerIp"`
}

type Job struct {
	DataJobID        int64   `json:"dataJobId"`
	DataJobName      string  `json:"dataJobName"`
	DataJobDesc      string  `json:"dataJobDesc"`
	UserName         string  `json:"userName"`
	DataJobType      string  `json:"dataJobType"`
	DataTaskState    string  `json:"dataTaskState"`
	CurrTaskStatus   string  `json:"currTaskStatus"`
	SourceDS         *Source `json:"sourceDsVO"`
	TargetDS         *Source `json:"targetDsVO"`
	SourceSchema     string  `json:"sourceSchema"`
	TargetSchema     string  `json:"targetSchema"`
	ConsoleJobID     int64   `json:"consoleJobId"`
	ConsoleTaskState string  `json:"consoleTaskState"`
	LifeCycleState   string  `json:"lifeCycleState"`
	HaveException    bool    `json:"haveException"`
	DataTasks        []Task  `json:"dataTasks"`
}

type JobSchema struct {
	SourceSchema          string `json:"sourceSchema"`
	TargetSchema          string `json:"targetSchema"`
	MappingConfig         string `json:"mappingConfig"`
	DefaultTopic          string `json:"defaultTopic"`
	DefaultTopicPartition int    `json:"defaultTopicPartition"`
	SchemaWhiteListLevel  string `json:"schemaWhiteListLevel"`
	SrcSchemaLessFormat   string `json:"srcSchemaLessFormat"`
	DstSchemaLessFormat   string `json:"dstSchemaLessFormat"`
}

type listJobsResponse struct {
	openapi.Response
	Data []Job `json:"data"`
}

type queryJobResponse struct {
	openapi.Response
	Data Job `json:"data"`
}

type queryJobSchemaResponse struct {
	openapi.Response
	Data JobSchema `json:"data"`
}

func (s *Service) ListJobs(options ListOptions) ([]Job, error) {
	var out listJobsResponse
	if err := s.client.PostJSON(listPath, newListJobsRequest(options), &out); err != nil {
		return nil, err
	}
	if err := openapi.EnsureSuccess(out.Response, "failed to list jobs"); err != nil {
		return nil, err
	}
	if out.Data == nil {
		return []Job{}, nil
	}
	return out.Data, nil
}

func (s *Service) GetJob(jobID int64) (Job, error) {
	var out queryJobResponse
	if err := s.client.PostJSON(queryPath, jobActionRequest{JobID: jobID}, &out); err != nil {
		return Job{}, err
	}
	if err := openapi.EnsureSuccess(out.Response, "failed to query job"); err != nil {
		return Job{}, err
	}
	return out.Data, nil
}

func (s *Service) GetJobSchema(jobID int64) (JobSchema, error) {
	var out queryJobSchemaResponse
	if err := s.client.PostJSON(schemaPath, jobActionRequest{JobID: jobID}, &out); err != nil {
		return JobSchema{}, err
	}
	if err := openapi.EnsureSuccess(out.Response, "failed to query job schema"); err != nil {
		return JobSchema{}, err
	}
	return out.Data, nil
}

func (s *Service) StartJob(jobID int64) error {
	var out openapi.Response
	if err := s.client.PostJSON(startPath, jobActionRequest{JobID: jobID}, &out); err != nil {
		return err
	}
	return openapi.EnsureSuccess(out, "failed to start job")
}

func (s *Service) StopJob(jobID int64) error {
	var out openapi.Response
	if err := s.client.PostJSON(stopPath, jobActionRequest{JobID: jobID}, &out); err != nil {
		return err
	}
	return openapi.EnsureSuccess(out, "failed to stop job")
}

func (s *Service) DeleteJob(jobID int64) error {
	var out openapi.Response
	if err := s.client.PostJSON(deletePath, jobActionRequest{JobID: jobID}, &out); err != nil {
		return err
	}
	return openapi.EnsureSuccess(out, "failed to delete job")
}

func (s *Service) ReplayJob(jobID int64, options ReplayOptions) error {
	var out openapi.Response
	if err := s.client.PostJSON(replayPath, newReplayJobRequest(jobID, options), &out); err != nil {
		return err
	}
	return openapi.EnsureSuccess(out, "failed to replay job")
}

func newListJobsRequest(options ListOptions) listJobsRequest {
	req := listJobsRequest{
		DataJobName: options.DataJobName,
		DataJobType: options.DataJobType,
		Desc:        options.Desc,
	}
	if options.SourceInstanceID > 0 {
		req.SourceInstanceID = ptrInt64(options.SourceInstanceID)
	}
	if options.TargetInstanceID > 0 {
		req.TargetInstanceID = ptrInt64(options.TargetInstanceID)
	}
	return req
}

func newReplayJobRequest(jobID int64, options ReplayOptions) replayJobRequest {
	req := replayJobRequest{JobID: jobID}
	if options.AutoStart {
		value := true
		req.AutoStart = &value
	}
	if options.ResetToCreated {
		value := true
		req.ResetToCreated = &value
	}
	return req
}

func ptrInt64(value int64) *int64 {
	return &value
}
