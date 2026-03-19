package cluster

import "cloudcanal-openapi-cli/internal/openapi"

const listPath = "/cloudcanal/console/api/v1/openapi/cluster/listclusters"

type Operations interface {
	List(options ListOptions) ([]Cluster, error)
}

type Service struct {
	client *openapi.Client
}

func NewService(client *openapi.Client) *Service {
	return &Service{client: client}
}

type ListOptions struct {
	CloudOrIDCName string
	ClusterDesc    string
	ClusterName    string
	Region         string
}

type listRequest struct {
	CloudOrIDCName string `json:"cloudOrIdcName,omitempty"`
	ClusterDesc    string `json:"clusterDescLike,omitempty"`
	ClusterName    string `json:"clusterNameLike,omitempty"`
	Region         string `json:"region,omitempty"`
}

type Cluster struct {
	ID             int64  `json:"id"`
	ClusterName    string `json:"clusterName"`
	Region         string `json:"region"`
	CloudOrIDCName string `json:"cloudOrIdcName"`
	ClusterDesc    string `json:"clusterDesc"`
	WorkerCount    int    `json:"workerCount"`
	RunningCount   int    `json:"runningCount"`
	AbnormalCount  int    `json:"abnormalCount"`
	OwnerName      string `json:"ownerName"`
}

type listResponse struct {
	openapi.Response
	Data []Cluster `json:"data"`
}

func (s *Service) List(options ListOptions) ([]Cluster, error) {
	var out listResponse
	if err := s.client.PostJSON(listPath, listRequest{
		CloudOrIDCName: options.CloudOrIDCName,
		ClusterDesc:    options.ClusterDesc,
		ClusterName:    options.ClusterName,
		Region:         options.Region,
	}, &out); err != nil {
		return nil, err
	}
	if err := openapi.EnsureSuccess(out.Response, "failed to list clusters"); err != nil {
		return nil, err
	}
	if out.Data == nil {
		return []Cluster{}, nil
	}
	return out.Data, nil
}
