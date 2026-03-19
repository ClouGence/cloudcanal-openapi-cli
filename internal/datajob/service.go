package datajob

import (
	"cloudcanal-openapi-cli/internal/openapi"
)

const (
	listPath           = "/cloudcanal/console/api/v1/openapi/datajob/list"
	queryPath          = "/cloudcanal/console/api/v1/openapi/datajob/queryjob"
	schemaPath         = "/cloudcanal/console/api/v1/openapi/datajob/queryjobschemabyid"
	createPath         = "/cloudcanal/console/api/v1/openapi/datajob/create"
	startPath          = "/cloudcanal/console/api/v1/openapi/datajob/start"
	stopPath           = "/cloudcanal/console/api/v1/openapi/datajob/stop"
	deletePath         = "/cloudcanal/console/api/v1/openapi/datajob/delete"
	replayPath         = "/cloudcanal/console/api/v1/openapi/datajob/replay"
	attachIncrePath    = "/cloudcanal/console/api/v1/openapi/datajob/attachincretask"
	detachIncrePath    = "/cloudcanal/console/api/v1/openapi/datajob/detachincretask"
	updateIncrePosPath = "/cloudcanal/console/api/v1/openapi/datajob/updateincrepos"
)

type Operations interface {
	ListJobs(options ListOptions) ([]Job, error)
	GetJob(jobID int64) (Job, error)
	GetJobSchema(jobID int64) (JobSchema, error)
	CreateJob(request CreateJobRequest) (CreateJobResult, error)
	StartJob(jobID int64) error
	StopJob(jobID int64) error
	DeleteJob(jobID int64) error
	ReplayJob(jobID int64, options ReplayOptions) error
	AttachIncreJob(jobID int64) error
	DetachIncreJob(jobID int64) error
	UpdateIncrePos(request UpdateIncrePosRequest) (UpdateIncrePosResult, error)
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

type createJobResponse struct {
	openapi.Response
	Data string `json:"data"`
}

type createJobRequest struct {
	ClusterID              *int64 `json:"clusterId,omitempty"`
	SrcDsID                *int64 `json:"srcDsId,omitempty"`
	DstDsID                *int64 `json:"dstDsId,omitempty"`
	SrcHostType            string `json:"srcHostType,omitempty"`
	DstHostType            string `json:"dstHostType,omitempty"`
	SchemaWhiteListLevel   string `json:"schemaWhiteListLevel,omitempty"`
	SrcSchema              string `json:"srcSchema,omitempty"`
	DstSchema              string `json:"dstSchema,omitempty"`
	MappingDef             string `json:"mappingDef,omitempty"`
	SrcCaseSensitiveType   string `json:"srcCaseSensitiveType,omitempty"`
	DstCaseSensitiveType   string `json:"dstCaseSensitiveType,omitempty"`
	SrcDsCharset           string `json:"srcDsCharset,omitempty"`
	TarDsCharset           string `json:"tarDsCharset,omitempty"`
	KeyConflictStrategy    string `json:"keyConflictStrategy,omitempty"`
	JobType                string `json:"jobType,omitempty"`
	DataJobDesc            string `json:"dataJobDesc,omitempty"`
	StructMigration        *bool  `json:"structMigration,omitempty"`
	InitialSync            *bool  `json:"initialSync,omitempty"`
	ShortTermSync          *bool  `json:"shortTermSync,omitempty"`
	ShortTermNum           *int   `json:"shortTermNum,omitempty"`
	FilterDDL              *bool  `json:"filterDDL,omitempty"`
	SpecID                 *int64 `json:"specId,omitempty"`
	AutoStart              *bool  `json:"autoStart,omitempty"`
	CheckOnce              *bool  `json:"checkOnce,omitempty"`
	CheckPeriod            *bool  `json:"checkPeriod,omitempty"`
	CheckPeriodCronExpr    string `json:"checkPeriodCronExpr,omitempty"`
	FullPeriod             *bool  `json:"fullPeriod,omitempty"`
	FullPeriodCronExpr     string `json:"fullPeriodCronExpr,omitempty"`
	DstMqDefaultTopic      string `json:"dstMqDefaultTopic,omitempty"`
	DstMqDefaultTopicParts *int   `json:"dstMqDefaultTopicPartitions,omitempty"`
	DstMqDdlTopic          string `json:"dstMqDdlTopic,omitempty"`
	DstMqDdlTopicParts     *int   `json:"dstMqDdlTopicPartitions,omitempty"`
	SrcSchemaLessFormat    string `json:"srcSchemaLessFormat,omitempty"`
	DstSchemaLessFormat    string `json:"dstSchemaLessFormat,omitempty"`
	OriginDecodeMsgFormat  string `json:"originDecodeMsgFormat,omitempty"`
	DstCkTableEngine       string `json:"dstCkTableEngine,omitempty"`
	DstSrOrDorisTableModel string `json:"dstSrOrDorisTableModel,omitempty"`
	KafkaConsumerGroupID   string `json:"kafkaConsumerGroupId,omitempty"`
	KuduNumReplicas        *int   `json:"kuduNumReplicas,omitempty"`
	SrcRocketMqGroupID     string `json:"srcRocketMqGroupId,omitempty"`
	SrcRabbitMqVhost       string `json:"srcRabbitMqVhost,omitempty"`
	SrcRabbitExchange      string `json:"srcRabbitExchange,omitempty"`
	DstRabbitMqVhost       string `json:"dstRabbitMqVhost,omitempty"`
	DstRabbitExchange      string `json:"dstRabbitExchange,omitempty"`
	ObTenant               string `json:"obTenant,omitempty"`
	DbHeartbeatEnable      *bool  `json:"dbHeartbeatEnable,omitempty"`
}

type CreateJobRequest struct {
	ClusterID              int64  `json:"clusterId"`
	SrcDsID                int64  `json:"srcDsId"`
	DstDsID                int64  `json:"dstDsId"`
	SrcHostType            string `json:"srcHostType,omitempty"`
	DstHostType            string `json:"dstHostType,omitempty"`
	SchemaWhiteListLevel   string `json:"schemaWhiteListLevel,omitempty"`
	SrcSchema              string `json:"srcSchema,omitempty"`
	DstSchema              string `json:"dstSchema,omitempty"`
	MappingDef             string `json:"mappingDef,omitempty"`
	SrcCaseSensitiveType   string `json:"srcCaseSensitiveType,omitempty"`
	DstCaseSensitiveType   string `json:"dstCaseSensitiveType,omitempty"`
	SrcDsCharset           string `json:"srcDsCharset,omitempty"`
	TarDsCharset           string `json:"tarDsCharset,omitempty"`
	KeyConflictStrategy    string `json:"keyConflictStrategy,omitempty"`
	JobType                string `json:"jobType,omitempty"`
	DataJobDesc            string `json:"dataJobDesc,omitempty"`
	StructMigration        *bool  `json:"structMigration,omitempty"`
	InitialSync            *bool  `json:"initialSync,omitempty"`
	ShortTermSync          *bool  `json:"shortTermSync,omitempty"`
	ShortTermNum           *int   `json:"shortTermNum,omitempty"`
	FilterDDL              *bool  `json:"filterDDL,omitempty"`
	SpecID                 *int64 `json:"specId,omitempty"`
	AutoStart              *bool  `json:"autoStart,omitempty"`
	CheckOnce              *bool  `json:"checkOnce,omitempty"`
	CheckPeriod            *bool  `json:"checkPeriod,omitempty"`
	CheckPeriodCronExpr    string `json:"checkPeriodCronExpr,omitempty"`
	FullPeriod             *bool  `json:"fullPeriod,omitempty"`
	FullPeriodCronExpr     string `json:"fullPeriodCronExpr,omitempty"`
	DstMqDefaultTopic      string `json:"dstMqDefaultTopic,omitempty"`
	DstMqDefaultTopicParts *int   `json:"dstMqDefaultTopicPartitions,omitempty"`
	DstMqDdlTopic          string `json:"dstMqDdlTopic,omitempty"`
	DstMqDdlTopicParts     *int   `json:"dstMqDdlTopicPartitions,omitempty"`
	SrcSchemaLessFormat    string `json:"srcSchemaLessFormat,omitempty"`
	DstSchemaLessFormat    string `json:"dstSchemaLessFormat,omitempty"`
	OriginDecodeMsgFormat  string `json:"originDecodeMsgFormat,omitempty"`
	DstCkTableEngine       string `json:"dstCkTableEngine,omitempty"`
	DstSrOrDorisTableModel string `json:"dstSrOrDorisTableModel,omitempty"`
	KafkaConsumerGroupID   string `json:"kafkaConsumerGroupId,omitempty"`
	KuduNumReplicas        *int   `json:"kuduNumReplicas,omitempty"`
	SrcRocketMqGroupID     string `json:"srcRocketMqGroupId,omitempty"`
	SrcRabbitMqVhost       string `json:"srcRabbitMqVhost,omitempty"`
	SrcRabbitExchange      string `json:"srcRabbitExchange,omitempty"`
	DstRabbitMqVhost       string `json:"dstRabbitMqVhost,omitempty"`
	DstRabbitExchange      string `json:"dstRabbitExchange,omitempty"`
	ObTenant               string `json:"obTenant,omitempty"`
	DbHeartbeatEnable      *bool  `json:"dbHeartbeatEnable,omitempty"`
}

type CreateJobResult struct {
	JobID string `json:"jobId"`
	Data  string `json:"data"`
}

type updateIncrePosRequest struct {
	TaskID            int64  `json:"taskId"`
	PosType           string `json:"posType"`
	JournalFile       string `json:"journalFile,omitempty"`
	FilePosition      *int64 `json:"filePosition,omitempty"`
	GtidPosition      string `json:"gtidPosition,omitempty"`
	PositionTimestamp *int64 `json:"positionTimestamp,omitempty"`
	ServerID          *int64 `json:"serverId,omitempty"`
	Lsn               string `json:"lsn,omitempty"`
	Scn               *int64 `json:"scn,omitempty"`
	ScnIndex          *int64 `json:"scnIndex,omitempty"`
	CommonPosStr      string `json:"commonPosStr,omitempty"`
	DataID            *int64 `json:"dataId,omitempty"`
	TransactionID     *int64 `json:"transactionId,omitempty"`
}

type UpdateIncrePosRequest struct {
	TaskID            int64  `json:"taskId"`
	PosType           string `json:"posType"`
	JournalFile       string `json:"journalFile,omitempty"`
	FilePosition      *int64 `json:"filePosition,omitempty"`
	GtidPosition      string `json:"gtidPosition,omitempty"`
	PositionTimestamp *int64 `json:"positionTimestamp,omitempty"`
	ServerID          *int64 `json:"serverId,omitempty"`
	Lsn               string `json:"lsn,omitempty"`
	Scn               *int64 `json:"scn,omitempty"`
	ScnIndex          *int64 `json:"scnIndex,omitempty"`
	CommonPosStr      string `json:"commonPosStr,omitempty"`
	DataID            *int64 `json:"dataId,omitempty"`
	TransactionID     *int64 `json:"transactionId,omitempty"`
}

type updateIncrePosResponse struct {
	openapi.Response
	Data string `json:"data"`
}

type UpdateIncrePosResult struct {
	Data string `json:"data"`
}

type actionRequest struct {
	JobID int64 `json:"jobId"`
}

type attachDetachRequest struct {
	JobID int64 `json:"jobId"`
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
	if err := s.client.PostJSONWithOptions(listPath, newListJobsRequest(options), &out, openapi.RequestOptions{Retryable: true}); err != nil {
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
	if err := s.client.PostJSONWithOptions(queryPath, jobActionRequest{JobID: jobID}, &out, openapi.RequestOptions{Retryable: true}); err != nil {
		return Job{}, err
	}
	if err := openapi.EnsureSuccess(out.Response, "failed to query job"); err != nil {
		return Job{}, err
	}
	return out.Data, nil
}

func (s *Service) GetJobSchema(jobID int64) (JobSchema, error) {
	var out queryJobSchemaResponse
	if err := s.client.PostJSONWithOptions(schemaPath, jobActionRequest{JobID: jobID}, &out, openapi.RequestOptions{Retryable: true}); err != nil {
		return JobSchema{}, err
	}
	if err := openapi.EnsureSuccess(out.Response, "failed to query job schema"); err != nil {
		return JobSchema{}, err
	}
	return out.Data, nil
}

func (s *Service) CreateJob(request CreateJobRequest) (CreateJobResult, error) {
	var out createJobResponse
	if err := s.client.PostJSONWithOptions(createPath, newCreateJobRequest(request), &out, openapi.RequestOptions{Retryable: true}); err != nil {
		return CreateJobResult{}, err
	}
	if err := openapi.EnsureSuccess(out.Response, "failed to create job"); err != nil {
		return CreateJobResult{}, err
	}
	return CreateJobResult{JobID: out.Data, Data: out.Data}, nil
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

func (s *Service) AttachIncreJob(jobID int64) error {
	var out openapi.Response
	if err := s.client.PostJSON(attachIncrePath, attachDetachRequest{JobID: jobID}, &out); err != nil {
		return err
	}
	return openapi.EnsureSuccess(out, "failed to attach incre task")
}

func (s *Service) DetachIncreJob(jobID int64) error {
	var out openapi.Response
	if err := s.client.PostJSON(detachIncrePath, attachDetachRequest{JobID: jobID}, &out); err != nil {
		return err
	}
	return openapi.EnsureSuccess(out, "failed to detach incre task")
}

func (s *Service) UpdateIncrePos(request UpdateIncrePosRequest) (UpdateIncrePosResult, error) {
	var out updateIncrePosResponse
	if err := s.client.PostJSONWithOptions(updateIncrePosPath, newUpdateIncrePosRequest(request), &out, openapi.RequestOptions{Retryable: true}); err != nil {
		return UpdateIncrePosResult{}, err
	}
	if err := openapi.EnsureSuccess(out.Response, "failed to update incre pos"); err != nil {
		return UpdateIncrePosResult{}, err
	}
	return UpdateIncrePosResult{Data: out.Data}, nil
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

func newCreateJobRequest(request CreateJobRequest) createJobRequest {
	req := createJobRequest{
		ClusterID:              ptrInt64(request.ClusterID),
		SrcDsID:                ptrInt64(request.SrcDsID),
		DstDsID:                ptrInt64(request.DstDsID),
		SrcHostType:            request.SrcHostType,
		DstHostType:            request.DstHostType,
		SchemaWhiteListLevel:   request.SchemaWhiteListLevel,
		SrcSchema:              request.SrcSchema,
		DstSchema:              request.DstSchema,
		MappingDef:             request.MappingDef,
		SrcCaseSensitiveType:   request.SrcCaseSensitiveType,
		DstCaseSensitiveType:   request.DstCaseSensitiveType,
		SrcDsCharset:           request.SrcDsCharset,
		TarDsCharset:           request.TarDsCharset,
		KeyConflictStrategy:    request.KeyConflictStrategy,
		JobType:                request.JobType,
		DataJobDesc:            request.DataJobDesc,
		StructMigration:        request.StructMigration,
		InitialSync:            request.InitialSync,
		ShortTermSync:          request.ShortTermSync,
		ShortTermNum:           request.ShortTermNum,
		FilterDDL:              request.FilterDDL,
		SpecID:                 request.SpecID,
		AutoStart:              request.AutoStart,
		CheckOnce:              request.CheckOnce,
		CheckPeriod:            request.CheckPeriod,
		CheckPeriodCronExpr:    request.CheckPeriodCronExpr,
		FullPeriod:             request.FullPeriod,
		FullPeriodCronExpr:     request.FullPeriodCronExpr,
		DstMqDefaultTopic:      request.DstMqDefaultTopic,
		DstMqDefaultTopicParts: request.DstMqDefaultTopicParts,
		DstMqDdlTopic:          request.DstMqDdlTopic,
		DstMqDdlTopicParts:     request.DstMqDdlTopicParts,
		SrcSchemaLessFormat:    request.SrcSchemaLessFormat,
		DstSchemaLessFormat:    request.DstSchemaLessFormat,
		OriginDecodeMsgFormat:  request.OriginDecodeMsgFormat,
		DstCkTableEngine:       request.DstCkTableEngine,
		DstSrOrDorisTableModel: request.DstSrOrDorisTableModel,
		KafkaConsumerGroupID:   request.KafkaConsumerGroupID,
		KuduNumReplicas:        request.KuduNumReplicas,
		SrcRocketMqGroupID:     request.SrcRocketMqGroupID,
		SrcRabbitMqVhost:       request.SrcRabbitMqVhost,
		SrcRabbitExchange:      request.SrcRabbitExchange,
		DstRabbitMqVhost:       request.DstRabbitMqVhost,
		DstRabbitExchange:      request.DstRabbitExchange,
		ObTenant:               request.ObTenant,
		DbHeartbeatEnable:      request.DbHeartbeatEnable,
	}
	return req
}

func newUpdateIncrePosRequest(request UpdateIncrePosRequest) updateIncrePosRequest {
	req := updateIncrePosRequest{
		TaskID:            request.TaskID,
		PosType:           request.PosType,
		JournalFile:       request.JournalFile,
		FilePosition:      request.FilePosition,
		GtidPosition:      request.GtidPosition,
		PositionTimestamp: request.PositionTimestamp,
		ServerID:          request.ServerID,
		Lsn:               request.Lsn,
		Scn:               request.Scn,
		ScnIndex:          request.ScnIndex,
		CommonPosStr:      request.CommonPosStr,
		DataID:            request.DataID,
		TransactionID:     request.TransactionID,
	}
	return req
}

func ptrInt64(value int64) *int64 {
	return &value
}
