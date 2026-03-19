package openapi_test

import (
	"cloudcanal-openapi-cli/internal/config"
	"cloudcanal-openapi-cli/internal/openapi"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClientPostsJSONWithSignedParams(t *testing.T) {
	var gotMethod string
	var gotBody map[string]any
	var gotQuery map[string]string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotQuery = map[string]string{
			"AccessKeyId":     r.URL.Query().Get("AccessKeyId"),
			"SignatureMethod": r.URL.Query().Get("SignatureMethod"),
			"SignatureNonce":  r.URL.Query().Get("SignatureNonce"),
			"Signature":       r.URL.Query().Get("Signature"),
		}
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}
		_, _ = w.Write([]byte(`{"code":"1","msg":"ok"}`))
	}))
	defer server.Close()

	client, err := openapi.NewClient(config.AppConfig{
		APIBaseURL: server.URL,
		AccessKey:  "test-ak",
		SecretKey:  "test-sk",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	var response map[string]string
	if err := client.PostJSON("/cloudcanal/console/api/v1/openapi/datajob/start", map[string]any{"jobId": 123}, &response); err != nil {
		t.Fatalf("PostJSON() error = %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Fatalf("method = %q, want POST", gotMethod)
	}
	if gotBody["jobId"] != float64(123) {
		t.Fatalf("body jobId = %#v, want 123", gotBody["jobId"])
	}
	if gotQuery["AccessKeyId"] != "test-ak" || gotQuery["SignatureMethod"] != "HmacSHA1" {
		t.Fatalf("unexpected query params: %#v", gotQuery)
	}
	if gotQuery["SignatureNonce"] == "" || gotQuery["Signature"] == "" {
		t.Fatalf("missing signature query params: %#v", gotQuery)
	}
}

func TestClientReturnsServerErrorForNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer server.Close()

	client, err := openapi.NewClient(config.AppConfig{
		APIBaseURL: server.URL,
		AccessKey:  "test-ak",
		SecretKey:  "test-sk",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.PostJSON("/cloudcanal/console/api/v1/openapi/datajob/list", map[string]any{}, &map[string]any{})
	if err == nil {
		t.Fatal("PostJSON() error = nil, want non-nil")
	}
	serverErr, ok := err.(*openapi.ServerError)
	if !ok {
		t.Fatalf("error type = %T, want *ServerError", err)
	}
	if serverErr.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", serverErr.StatusCode, http.StatusInternalServerError)
	}
}

func TestClientReturnsErrorWhenConnectionFails(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen() error = %v", err)
	}
	addr := listener.Addr().String()
	_ = listener.Close()

	httpClient := &http.Client{Timeout: 200 * time.Millisecond}
	client, err := openapi.NewClientWithHTTP(config.AppConfig{
		APIBaseURL: "http://" + addr,
		AccessKey:  "test-ak",
		SecretKey:  "test-sk",
	}, httpClient)
	if err != nil {
		t.Fatalf("NewClientWithHTTP() error = %v", err)
	}

	if err := client.PostJSON("/cloudcanal/console/api/v1/openapi/datajob/list", map[string]any{}, &map[string]any{}); err == nil {
		t.Fatal("PostJSON() error = nil, want non-nil")
	}
}

func TestProbeAuthenticationAcceptsPermissionDeniedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code":"0","msg":"permission denied for datajob list"}`))
	}))
	defer server.Close()

	client, err := openapi.NewClient(config.AppConfig{
		APIBaseURL: server.URL,
		AccessKey:  "test-ak",
		SecretKey:  "test-sk",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if err := client.ProbeAuthentication(); err != nil {
		t.Fatalf("ProbeAuthentication() error = %v, want nil", err)
	}
}

func TestProbeAuthenticationRejectsInvalidSignatureStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(497)
		_, _ = w.Write([]byte(`{"code":"0","msg":"invalid signature"}`))
	}))
	defer server.Close()

	client, err := openapi.NewClient(config.AppConfig{
		APIBaseURL: server.URL,
		AccessKey:  "test-ak",
		SecretKey:  "test-sk",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.ProbeAuthentication()
	if err == nil {
		t.Fatal("ProbeAuthentication() error = nil, want non-nil")
	}
	if got := err.Error(); got != "invalid signature" {
		t.Fatalf("error = %q, want invalid signature", got)
	}
}

func TestProbeAuthenticationRejectsUnexpectedApplicationFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code":"0","msg":"backend not ready"}`))
	}))
	defer server.Close()

	client, err := openapi.NewClient(config.AppConfig{
		APIBaseURL: server.URL,
		AccessKey:  "test-ak",
		SecretKey:  "test-sk",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.ProbeAuthentication()
	if err == nil {
		t.Fatal("ProbeAuthentication() error = nil, want non-nil")
	}
	if got := err.Error(); got != "OpenAPI probe failed: backend not ready" {
		t.Fatalf("error = %q, want probe failure message", got)
	}
}

func TestClientRetriesRetryableRequestsOnServerFailure(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			http.Error(w, "temporary failure", http.StatusBadGateway)
			return
		}
		_, _ = w.Write([]byte(`{"code":"1","msg":"ok"}`))
	}))
	defer server.Close()

	client, err := openapi.NewClient(config.AppConfig{
		APIBaseURL:                 server.URL,
		AccessKey:                  "test-ak",
		SecretKey:                  "test-sk",
		HTTPReadMaxRetries:         2,
		HTTPReadRetryBackoffMillis: 1,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	var response map[string]string
	if err := client.PostJSONWithOptions("/cloudcanal/console/api/v1/openapi/datajob/list", map[string]any{}, &response, openapi.RequestOptions{Retryable: true}); err != nil {
		t.Fatalf("PostJSONWithOptions() error = %v", err)
	}
	if attempts != 3 {
		t.Fatalf("attempts = %d, want 3", attempts)
	}
}

func TestClientDoesNotRetryNonRetryableRequests(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		http.Error(w, "temporary failure", http.StatusBadGateway)
	}))
	defer server.Close()

	client, err := openapi.NewClient(config.AppConfig{
		APIBaseURL:                 server.URL,
		AccessKey:                  "test-ak",
		SecretKey:                  "test-sk",
		HTTPReadMaxRetries:         2,
		HTTPReadRetryBackoffMillis: 1,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.PostJSON("/cloudcanal/console/api/v1/openapi/datajob/start", map[string]any{"jobId": 1}, &map[string]any{})
	if err == nil {
		t.Fatal("PostJSON() error = nil, want non-nil")
	}
	if attempts != 1 {
		t.Fatalf("attempts = %d, want 1", attempts)
	}
}
