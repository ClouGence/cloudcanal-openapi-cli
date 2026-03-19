package jobconfig

import "cloudcanal-openapi-cli/internal/openapi"

const listSpecsPath = "/cloudcanal/console/api/v1/openapi/constant/listspecs"

type Operations interface {
	ListSpecs(options ListSpecsOptions) ([]Spec, error)
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

type listSpecsRequest struct {
	DataJobType   string `json:"dataJobType,omitempty"`
	InitialSync   *bool  `json:"initialSync,omitempty"`
	ShortTermSync *bool  `json:"shortTermSync,omitempty"`
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
