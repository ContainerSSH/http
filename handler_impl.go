package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	goHttp "net/http"

	"github.com/containerssh/log"
)

type serverResponse struct {
	statusCode uint16
	body       interface{}
}

func (s *serverResponse) Error() string {
	return fmt.Sprintf("%v", s.body)
}

func (s *serverResponse) SetStatus(statusCode uint16) {
	s.statusCode = statusCode
}

func (s *serverResponse) SetBody(body interface{}) {
	s.body = body
}

type handler struct {
	requestHandler RequestHandler
	logger         log.Logger
}

var internalErrorResponse = serverResponse{
	500,
	map[string]string{"error": "Internal Server Error"},
}

var badRequestResponse = serverResponse{
	400,
	map[string]string{"error": "Bad Request"},
}

func (h *handler) ServeHTTP(goWriter goHttp.ResponseWriter, goRequest *goHttp.Request) {
	response := serverResponse{
		statusCode: 200,
		body:       nil,
	}
	if err := h.requestHandler.OnRequest(
		&internalRequest{
			request: goRequest,
			writer:  goWriter,
		},
		&response,
	); err != nil {
		if errors.Is(err, &badRequestResponse) {
			response = badRequestResponse
		} else {
			h.logger.Warningf("handler returned error response (%w)", err)
			response = internalErrorResponse
		}
	}
	bytes, err := json.Marshal(response.body)
	if err != nil {
		h.logger.Errorf("failed to marshal response %v (%w)", response, err)
		response = internalErrorResponse
		bytes, err = json.Marshal(internalErrorResponse.body)
		if err != nil {
			panic(fmt.Errorf("bug: failed to marshal internal server error JSON response (%w)", err))
		}
	}
	goWriter.WriteHeader(int(response.statusCode))
	goWriter.Header().Add("Content-Type", "application/json")
	if _, err := goWriter.Write(bytes); err != nil {
		h.logger.Infof("failed to write HTTP response")
	}
}

type internalRequest struct {
	writer  goHttp.ResponseWriter
	request *goHttp.Request
}

func (i *internalRequest) Decode(target interface{}) error {
	bytes, err := ioutil.ReadAll(i.request.Body)
	if err != nil {
		return &badRequestResponse
	}
	return json.Unmarshal(bytes, target)
}
