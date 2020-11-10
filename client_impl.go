package http

import (
	"bytes"
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

func (c *client) Post(path string, requestBody interface{}, responseBody interface{}) (int, error) {
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
	buffer := &bytes.Buffer{}
	err := json.NewEncoder(buffer).Encode(requestBody)
	if err != nil {
		//This is a bug
		return 0, err
	}
	req, err := http.NewRequest(
		method,
		fmt.Sprintf("%s%s", c.config.Url, path),
		buffer,
	)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}

	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(responseBody); err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}
