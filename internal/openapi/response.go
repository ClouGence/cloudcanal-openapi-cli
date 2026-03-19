package openapi

import "errors"

type Response struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func EnsureSuccess(resp Response, fallback string) error {
	if resp.Code == "1" {
		return nil
	}
	if resp.Msg != "" {
		return errors.New(resp.Msg)
	}
	return errors.New(fallback)
}
