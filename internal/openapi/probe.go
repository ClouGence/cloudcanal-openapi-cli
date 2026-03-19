package openapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	probePath              = "/cloudcanal/console/api/v1/openapi/datajob/list"
	httpStatusNoAccessKey  = 498
	httpStatusBadSignature = 497
)

// ProbeAuthentication verifies that the configured OpenAPI endpoint is reachable
// and the supplied AK/SK can pass signature validation. It intentionally treats
// application-level permission errors as success so initialization is not blocked
// for users who can access other resources but cannot list jobs.
func (c *Client) ProbeAuthentication() error {
	var out Response
	if err := c.PostJSONWithOptions(probePath, map[string]any{}, &out, RequestOptions{Retryable: true}); err != nil {
		return classifyProbeError(err)
	}

	switch {
	case out.Code == "1":
		return nil
	case isAuthenticationFailureMessage(out.Msg):
		return errors.New(strings.TrimSpace(out.Msg))
	case isPermissionErrorMessage(out.Msg):
		return nil
	default:
		return fmt.Errorf("OpenAPI probe failed: %s", probeMessage(out))
	}
}

func classifyProbeError(err error) error {
	var serverErr *ServerError
	if !errors.As(err, &serverErr) {
		return err
	}

	switch serverErr.StatusCode {
	case httpStatusBadSignature, httpStatusNoAccessKey:
		message := probeMessageFromBody(serverErr.ResponseBody)
		if message == "" {
			message = strings.TrimSpace(serverErr.ResponseBody)
		}
		if message == "" {
			message = serverErr.Error()
		}
		return errors.New(message)
	default:
		return err
	}
}

func probeMessage(resp Response) string {
	switch {
	case strings.TrimSpace(resp.Msg) != "":
		return strings.TrimSpace(resp.Msg)
	case strings.TrimSpace(resp.Code) != "":
		return "code=" + strings.TrimSpace(resp.Code)
	default:
		return "unknown response"
	}
}

func probeMessageFromBody(body string) string {
	var resp Response
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		return ""
	}
	return strings.TrimSpace(resp.Msg)
}

func isAuthenticationFailureMessage(message string) bool {
	lower := strings.ToLower(strings.TrimSpace(message))
	if lower == "" {
		return false
	}

	patterns := []string{
		"invalid signature",
		"signature is invalid",
		"signature mismatch",
		"accesskey",
		"access key",
		"ak/sk",
		"认证失败",
		"签名无效",
		"签名错误",
		"accesskeyid对应的用户不存在",
		"aksk",
	}
	for _, pattern := range patterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}
	return false
}

func isPermissionErrorMessage(message string) bool {
	lower := strings.ToLower(strings.TrimSpace(message))
	if lower == "" {
		return false
	}

	patterns := []string{
		"permission",
		"forbidden",
		"access denied",
		"not authorized",
		"no auth",
		"no permission",
		"无权限",
		"没有权限",
		"权限不足",
		"禁止访问",
	}
	for _, pattern := range patterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}
	return false
}
