package datasource

import (
	"bytes"
	"cloudcanal-openapi-cli/internal/openapi"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	listPath   = "/cloudcanal/console/api/v1/openapi/datasource/listds"
	addPath    = "/cloudcanal/console/api/v1/openapi/datasource/addds"
	deletePath = "/cloudcanal/console/api/v1/openapi/datasource/deleteds"
)

type Operations interface {
	List(options ListOptions) ([]DataSource, error)
	Get(dataSourceID int64) (DataSource, error)
	Add(options AddOptions) (string, error)
	Delete(dataSourceID int64) error
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

type AddOptions struct {
	DataSourceAddData ApiDsAddData `json:"dataSourceAddData"`
	SecurityFilePath  string       `json:"securityFilePath,omitempty"`
	SecretFilePath    string       `json:"secretFilePath,omitempty"`
}

type listRequest struct {
	DataSourceID   *int64 `json:"dataSourceId,omitempty"`
	DeployType     string `json:"deployType,omitempty"`
	HostType       string `json:"hostType,omitempty"`
	LifeCycleState string `json:"lifeCycleState,omitempty"`
	Type           string `json:"type,omitempty"`
}

type addResponse struct {
	openapi.Response
	Data string `json:"data"`
}

type deleteRequest struct {
	DataSourceID int64 `json:"dataSourceId"`
}

type ApiDsAddData struct {
	DeployType               string         `json:"deployType,omitempty"`
	Region                   string         `json:"region,omitempty"`
	Type                     string         `json:"type,omitempty"`
	DbName                   string         `json:"dbName,omitempty"`
	Host                     string         `json:"host,omitempty"`
	PrivateHost              string         `json:"privateHost,omitempty"`
	PublicHost               string         `json:"publicHost,omitempty"`
	HostType                 string         `json:"hostType,omitempty"`
	InstanceDesc             string         `json:"instanceDesc,omitempty"`
	InstanceID               string         `json:"instanceId,omitempty"`
	AutoCreateAccount        bool           `json:"autoCreateAccount,omitempty"`
	Account                  string         `json:"account,omitempty"`
	Password                 string         `json:"password,omitempty"`
	AccessKey                string         `json:"accessKey,omitempty"`
	SecretKey                string         `json:"secretKey,omitempty"`
	SecurityType             string         `json:"securityType,omitempty"`
	ExtraData                string         `json:"extraData,omitempty"`
	ClusterIDs               []int64        `json:"clusterIds,omitempty"`
	LifeCycleState           string         `json:"lifeCycleState,omitempty"`
	ClientTrustStorePassword string         `json:"clientTrustStorePassword,omitempty"`
	WhiteListAddType         string         `json:"whiteListAddType,omitempty"`
	ParentDsID               int64          `json:"parentDsId,omitempty"`
	Version                  string         `json:"version,omitempty"`
	Driver                   string         `json:"driver,omitempty"`
	ConnectType              string         `json:"connectType,omitempty"`
	DsKvConfigs              []KvBaseConfig `json:"dsKvConfigs,omitempty"`
}

type KvBaseConfig struct {
	ConfigName  string `json:"configName,omitempty"`
	ConfigValue string `json:"configValue,omitempty"`
}

type DataSource struct {
	ID               int64     `json:"id"`
	InstanceID       string    `json:"instanceId"`
	DeployType       string    `json:"deployType"`
	Region           string    `json:"region"`
	DataSourceType   string    `json:"dataSourceType"`
	HostType         string    `json:"hostType"`
	InstanceDesc     string    `json:"instanceDesc"`
	ConsoleJobID     Stringish `json:"consoleJobId"`
	ConsoleTaskState string    `json:"consoleTaskState"`
	AccountName      string    `json:"accountName"`
	LifeCycleState   string    `json:"lifeCycleState"`
	SecurityType     string    `json:"securityType"`
}

type Stringish string

func (s *Stringish) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*s = ""
		return nil
	}

	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		*s = Stringish(text)
		return nil
	}

	var number json.Number
	if err := json.Unmarshal(data, &number); err == nil {
		*s = Stringish(number.String())
		return nil
	}

	return fmt.Errorf("unsupported stringish value: %s", string(data))
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

func (s *Service) Add(options AddOptions) (string, error) {
	body, contentType, err := buildAddMultipartBody(options)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, s.client.SignedURL(addPath), body)
	if err != nil {
		return "", fmt.Errorf("failed to build add datasource request: %w", err)
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := s.client.HTTPClient().Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call OpenAPI: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read OpenAPI response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", &openapi.ServerError{StatusCode: resp.StatusCode, ResponseBody: string(responseBody)}
	}

	var out addResponse
	if err := json.Unmarshal(responseBody, &out); err != nil {
		return "", fmt.Errorf("failed to parse OpenAPI response: %w", err)
	}
	if err := openapi.EnsureSuccess(out.Response, "failed to add data source"); err != nil {
		return "", err
	}
	return out.Data, nil
}

func (s *Service) Delete(dataSourceID int64) error {
	var out openapi.Response
	if err := s.client.PostJSONWithOptions(deletePath, deleteRequest{DataSourceID: dataSourceID}, &out, openapi.RequestOptions{Retryable: true}); err != nil {
		return err
	}
	return openapi.EnsureSuccess(out, "failed to delete data source")
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

func buildAddMultipartBody(options AddOptions) (*bytes.Buffer, string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	dataJSON, err := json.Marshal(options.DataSourceAddData)
	if err != nil {
		return nil, "", fmt.Errorf("failed to encode dataSourceAddData: %w", err)
	}
	if err := writer.WriteField("dataSourceAddData", string(dataJSON)); err != nil {
		return nil, "", fmt.Errorf("failed to write dataSourceAddData: %w", err)
	}

	if err := addMultipartFile(writer, "securityFile", options.SecurityFilePath); err != nil {
		return nil, "", err
	}
	if err := addMultipartFile(writer, "secretFile", options.SecretFilePath); err != nil {
		return nil, "", err
	}
	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("failed to close multipart writer: %w", err)
	}
	return &body, writer.FormDataContentType(), nil
}

func addMultipartFile(writer *multipart.Writer, fieldName string, filePath string) error {
	if strings.TrimSpace(filePath) == "" {
		return nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", fieldName, err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("failed to create multipart file field %s: %w", fieldName, err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy %s: %w", fieldName, err)
	}
	return nil
}
