package datasource

import (
	"cloudcanal-openapi-cli/internal/openapi"
	"fmt"
)

const listPath = "/cloudcanal/console/api/v1/openapi/datasource/listds"

type Operations interface {
	List(options ListOptions) ([]DataSource, error)
	Get(dataSourceID int64) (DataSource, error)
}

type Service struct {
	client *openapi.Client
}

func NewService(client *openapi.Client) *Service {
	return &Service{client: client}
}

type ListOptions struct {
	DataSourceID   int64
	DeployType     string
	HostType       string
	LifeCycleState string
	Type           string
}

type listRequest struct {
	DataSourceID   *int64 `json:"dataSourceId,omitempty"`
	DeployType     string `json:"deployType,omitempty"`
	HostType       string `json:"hostType,omitempty"`
	LifeCycleState string `json:"lifeCycleState,omitempty"`
	Type           string `json:"type,omitempty"`
}

type DataSource struct {
	ID               int64  `json:"id"`
	InstanceID       string `json:"instanceId"`
	DeployType       string `json:"deployType"`
	Region           string `json:"region"`
	DataSourceType   string `json:"dataSourceType"`
	HostType         string `json:"hostType"`
	InstanceDesc     string `json:"instanceDesc"`
	ConsoleJobID     string `json:"consoleJobId"`
	ConsoleTaskState string `json:"consoleTaskState"`
	AccountName      string `json:"accountName"`
	LifeCycleState   string `json:"lifeCycleState"`
	SecurityType     string `json:"securityType"`
}

type listResponse struct {
	openapi.Response
	Data []DataSource `json:"data"`
}

func (s *Service) List(options ListOptions) ([]DataSource, error) {
	var out listResponse
	if err := s.client.PostJSONWithOptions(listPath, newListRequest(options), &out, openapi.RequestOptions{Retryable: true}); err != nil {
		return nil, err
	}
	if err := openapi.EnsureSuccess(out.Response, "failed to list data sources"); err != nil {
		return nil, err
	}
	if out.Data == nil {
		return []DataSource{}, nil
	}
	return out.Data, nil
}

func (s *Service) Get(dataSourceID int64) (DataSource, error) {
	sources, err := s.List(ListOptions{DataSourceID: dataSourceID})
	if err != nil {
		return DataSource{}, err
	}
	if len(sources) == 0 {
		return DataSource{}, fmt.Errorf("dataSourceId %d not found", dataSourceID)
	}
	return sources[0], nil
}

func newListRequest(options ListOptions) listRequest {
	req := listRequest{
		DeployType:     options.DeployType,
		HostType:       options.HostType,
		LifeCycleState: options.LifeCycleState,
		Type:           options.Type,
	}
	if options.DataSourceID > 0 {
		req.DataSourceID = &options.DataSourceID
	}
	return req
}
