package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/containerssh/log"
)

type client struct {
	config    ClientConfiguration
	logger    log.Logger
	tlsConfig *tls.Config
}

func (c *client) Post(
	path string,
	requestBody interface{},
	responseBody interface{},
) (
	int,
	error,
) {
	return c.request(
		http.MethodPost,
		path,
		requestBody,
		responseBody,
	)
}

func (c *client) request(
	method string,
	path string,
	requestBody interface{},
	responseBody interface{},
) (int, error) {
	logger := c.logger.WithLabel("method", method).WithLabel("path", path)

	httpClient := c.createHTTPClient(logger)

	req, err := c.createRequest(method, path, requestBody, logger)
	if err != nil {
		return 0, err
	}

	logger.Debug(log.NewMessage(MClientRequest, "HTTP %s request to %s%s", method, c.config.URL, path))

	resp, err := httpClient.Do(req)
	if err != nil {
		var typedError log.Message
		if errors.As(err, &typedError) {
			return 0, err
		}
		err = log.Wrap(err, EFailureConnectionFailed, "HTTP %s request to %s%s failed", method, c.config.URL, path)
		logger.Debug(err)
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	logger.Debug(log.NewMessage(
		MClientResponse,
		"HTTP response with status %d",
		resp.StatusCode,
	).Label("statusCode", resp.StatusCode))

	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(responseBody); err != nil {
		err = log.Wrap(err, EFailureDecodeFailed, "Failed to decode HTTP response")
		logger.Debug(err)
		return resp.StatusCode, err
	}
	return resp.StatusCode, nil
}

func (c *client) createRequest(method string, path string, requestBody interface{}, logger log.Logger) (
	*http.Request,
	error,
) {
	buffer := &bytes.Buffer{}
	err := json.NewEncoder(buffer).Encode(requestBody)
	if err != nil {
		//This is a bug
		err := log.Wrap(err, EFailureEncodeFailed, "BUG: HTTP request encoding failed")
		logger.Debug(err)
		return nil, err
	}
	req, err := http.NewRequest(
		method,
		fmt.Sprintf("%s%s", c.config.URL, path),
		buffer,
	)
	if err != nil {
		err := log.Wrap(err, EFailureEncodeFailed, "BUG: HTTP request encoding failed")
		logger.Debug(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *client) createHTTPClient(logger log.Logger) *http.Client {
	transport := &http.Transport{
		TLSClientConfig: c.tlsConfig,
	}

	httpClient := &http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !c.config.AllowRedirects {
				return log.NewMessage(
					EClientRedirectsDisabled,
					"Redirects disabled, server tried to redirect to %s", req.URL,
				).Label("redirect", req.URL)
			}
			logger.Debug(
				log.NewMessage(
					MClientRedirect, "HTTP redirect to %s", req.URL,
				).Label("redirect", req.URL),
			)
			return nil
		},
		Timeout: c.config.Timeout,
	}
	return httpClient
}
