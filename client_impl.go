package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/containerssh/log"
)

type client struct {
	config     ClientConfiguration
	logger     log.Logger
	tlsConfig  *tls.Config
	httpClient *http.Client
}

func (c *client) Post(
	ctx context.Context,
	path string,
	requestBody interface{},
	responseBody interface{},
) (
	int,
	error,
) {
	return c.request(
		ctx,
		http.MethodPost,
		path,
		requestBody,
		responseBody,
	)
}

func (c *client) request(
	ctx context.Context,
	method string,
	path string,
	requestBody interface{},
	responseBody interface{},
) (int, error) {
	buffer := &bytes.Buffer{}
	err := json.NewEncoder(buffer).Encode(requestBody)
	if err != nil {
		//This is a bug
		return 0, ClientError{
			Reason:  FailureReasonEncodeFailed,
			Cause:   err,
			Message: "failed to encode request body",
		}
	}
	req, err := http.NewRequestWithContext(
		ctx,
		method,
		fmt.Sprintf("%s%s", c.config.URL, path),
		buffer,
	)
	if err != nil {
		return 0, &ClientError{
			Reason:  FailureReasonEncodeFailed,
			Cause:   err,
			Message: "failed to encode request body",
		}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, ClientError{
			Reason:  FailureReasonConnectionFailed,
			Cause:   err,
			Message: "failed on HTTP request",
		}
	}
	defer func() { _ = resp.Body.Close() }()

	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(responseBody); err != nil {
		return resp.StatusCode, ClientError{
			Reason:  FailureReasonDecodeFailed,
			Cause:   err,
			Message: "failed to decode response",
		}
	}
	return resp.StatusCode, nil
}
