package schema

import "cloudcanal-openapi-cli/internal/openapi"

const listTransObjsByMetaPath = "/cloudcanal/console/api/v1/openapi/schema/listTransObjsByMeta"

type Operations interface {
	ListTransObjsByMeta(options ListTransObjsByMetaOptions) ([]ApiTransferObjIndexDO, error)
}

type Service struct {
	client *openapi.Client
}

func NewService(client *openapi.Client) *Service {
	return &Service{client: client}
}

type ListTransObjsByMetaOptions struct {
	SrcDb       string
	SrcSchema   string
	SrcTransObj string
	DstDb       string
	DstSchema   string
	DstTranObj  string
}

type listTransObjsByMetaRequest struct {
	SrcDb       string `json:"srcDb,omitempty"`
	SrcSchema   string `json:"srcSchema,omitempty"`
	SrcTransObj string `json:"srcTransObj,omitempty"`
	DstDb       string `json:"dstDb,omitempty"`
	DstSchema   string `json:"dstSchema,omitempty"`
	DstTranObj  string `json:"dstTranObj,omitempty"`
}

type ApiTransferObjIndexDO struct {
	ID                     int64  `json:"id"`
	SrcTransferObjName     string `json:"srcTransferObjName"`
	SrcDsInstanceID        int64  `json:"srcDsInstanceId"`
	SrcDsInstanceName      string `json:"srcDsInstanceName"`
	SrcFullTransferObjName string `json:"srcFullTransferObjName"`
	SrcDsType              string `json:"srcDsType"`
	SrcDb                  string `json:"srcDb"`
	SrcSchema              string `json:"srcSchema"`
	FilterExpr             string `json:"filterExpr"`
	SpecifiedPks           string `json:"specifiedPks"`
	DstTransferObjName     string `json:"dstTransferObjName"`
	DstDsInstanceID        int64  `json:"dstDsInstanceId"`
	DstDsInstanceName      string `json:"dstDsInstanceName"`
	DstFullTransferObjName string `json:"dstFullTransferObjName"`
	DstDsType              string `json:"dstDsType"`
	DstDb                  string `json:"dstDb"`
	DstSchema              string `json:"dstSchema"`
	DataJobID              int64  `json:"dataJobId"`
	DataJobName            string `json:"dataJobName"`
	DataJobDesc            string `json:"dataJobDesc"`
}

type listTransObjsByMetaResponse struct {
	openapi.Response
	Data []ApiTransferObjIndexDO `json:"data"`
}

func (s *Service) ListTransObjsByMeta(options ListTransObjsByMetaOptions) ([]ApiTransferObjIndexDO, error) {
	var out listTransObjsByMetaResponse
	if err := s.client.PostJSONWithOptions(listTransObjsByMetaPath, listTransObjsByMetaRequest{
		SrcDb:       options.SrcDb,
		SrcSchema:   options.SrcSchema,
		SrcTransObj: options.SrcTransObj,
		DstDb:       options.DstDb,
		DstSchema:   options.DstSchema,
		DstTranObj:  options.DstTranObj,
	}, &out, openapi.RequestOptions{Retryable: true}); err != nil {
		return nil, err
	}
	if err := openapi.EnsureSuccess(out.Response, "failed to list transfer objects by meta"); err != nil {
		return nil, err
	}
	if out.Data == nil {
		return []ApiTransferObjIndexDO{}, nil
	}
	return out.Data, nil
}
