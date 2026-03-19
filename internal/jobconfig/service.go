package jobconfig

import (
	"cloudcanal-openapi-cli/internal/openapi"
	"encoding/json"
)

const (
	listSpecsPath        = "/cloudcanal/console/api/v1/openapi/constant/listspecs"
	transformJobTypePath = "/cloudcanal/console/api/v1/openapi/constant/transformjobtype"
)

type Operations interface {
	ListSpecs(options ListSpecsOptions) ([]Spec, error)
	TransformJobType(options TransformJobTypeOptions) (TransformJobTypeResponse, error)
}

type Service struct {
	client *openapi.Client
}

func NewService(client *openapi.Client) *Service {
	return &Service{client: client}
}

type ListSpecsOptions struct {
	DataJobType   string
	InitialSync   *bool
	ShortTermSync *bool
}

type TransformJobTypeOptions struct {
	SourceType string
	TargetType string
}

type listSpecsRequest struct {
	DataJobType   string `json:"dataJobType,omitempty"`
	InitialSync   *bool  `json:"initialSync,omitempty"`
	ShortTermSync *bool  `json:"shortTermSync,omitempty"`
}

type transformJobTypeRequest struct {
	SourceType string `json:"sourceType,omitempty"`
	TargetType string `json:"targetType,omitempty"`
}

type Spec struct {
	ID            int64  `json:"id"`
	SpecKind      string `json:"specKind"`
	SpecKindCN    string `json:"specKindCn"`
	Spec          string `json:"spec"`
	Description   string `json:"description"`
	FullMemoryMB  int    `json:"fullMemoryMb"`
	IncreMemoryMB int    `json:"increMemoryMb"`
	CheckMemoryMB int    `json:"checkMemoryMb"`
}

type listSpecsResponse struct {
	openapi.Response
	Data []Spec `json:"data"`
}

type TransformJobTypeResponse struct {
	openapi.Response
	Data json.RawMessage `json:"data"`
}

func (s *Service) ListSpecs(options ListSpecsOptions) ([]Spec, error) {
	var out listSpecsResponse
	if err := s.client.PostJSONWithOptions(listSpecsPath, listSpecsRequest{
		DataJobType:   options.DataJobType,
		InitialSync:   options.InitialSync,
		ShortTermSync: options.ShortTermSync,
	}, &out, openapi.RequestOptions{Retryable: true}); err != nil {
		return nil, err
	}
	if err := openapi.EnsureSuccess(out.Response, "failed to list job specs"); err != nil {
		return nil, err
	}
	if out.Data == nil {
		return []Spec{}, nil
	}
	return out.Data, nil
}

func (s *Service) TransformJobType(options TransformJobTypeOptions) (TransformJobTypeResponse, error) {
	var out TransformJobTypeResponse
	if err := s.client.PostJSONWithOptions(transformJobTypePath, transformJobTypeRequest{
		SourceType: options.SourceType,
		TargetType: options.TargetType,
	}, &out, openapi.RequestOptions{Retryable: true}); err != nil {
		return TransformJobTypeResponse{}, err
	}
	if err := openapi.EnsureSuccess(out.Response, "failed to transform job type"); err != nil {
		return TransformJobTypeResponse{}, err
	}
	return out, nil
}
