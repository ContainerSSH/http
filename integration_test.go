package http_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/containerssh/log/standard"
	"github.com/stretchr/testify/assert"

	"github.com/containerssh/http"
)

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	Error bool `json:"error"`
	Message string `json:"message"`
}

type handler struct {
}

func (s *handler) OnRequest(request http.ServerRequest, response http.ServerResponse) error {
	req := Request{}
	if err := request.Decode(&req); err != nil {
		return err
	}
	if req.Message == "Hi" {
		response.SetBody(&Response{
			Error:   false,
			Message: "Hello world!",
		})
	} else {
		response.SetBody(&Response{
			Error:   true,
			Message: "Be nice and greet me!",
		})
	}
	return nil
}

func TestUnencrypted(t *testing.T) {
	logger := standard.New()
	clientConfig := http.ClientConfiguration{
		Url:        "http://127.0.0.1:8080/",
		Timeout:    2 * time.Second,
	}

	client, err := http.NewClient(clientConfig, logger)
	if err != nil {
		assert.Fail(t, "failed to create client", err)
		return
	}

	server, err := http.NewServer(
		http.ServerConfiguration{
			Listen:       "127.0.0.1:8080",
		},
		http.NewServerHandler(&handler{}, logger),
		logger,
	)
	if err != nil {
		assert.Fail(t, "failed to create server", err)
		return
	}

	response := Response{}
	errorChannel := make(chan error, 2)
	clientDone := make(chan bool, 1)
	responseStatus := 0
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		if responseStatus, err = client.Post(
			"",
			&Request{Message: "Hi"},
			&response,
		); err != nil {
			errorChannel <- err
		}
		clientDone <- true
	}()
	go func() {
		defer wg.Done()
		if err := server.Run(); err != nil {
			errorChannel <- err
		}
	}()
	<-clientDone
	server.Shutdown(context.Background())
	wg.Wait()
	finished := false
	for {
		select {
		case err := <-errorChannel:
			assert.Fail(t, "error while executing HTTP query", err)
		default:
			finished = true
		}
		if finished {
			break
		}
	}
	assert.Equal(t, 200, responseStatus)
	assert.Equal(t, false, response.Error)
	assert.Equal(t, "Hello world!", response.Message)
}
