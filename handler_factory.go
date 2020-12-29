package http

import (
	goHttp "net/http"

	"github.com/containerssh/log"
)

// NewServerHandler creates a new simplified HTTP handler that decodes JSON requests and encodes JSON responses.
func NewServerHandler(
	requestHandler RequestHandler,
	logger log.Logger,
) goHttp.Handler {
	if requestHandler == nil {
		panic("BUG: no requestHandler provided to http.NewServerHandler")
	}
	if logger == nil {
		panic("BUG: no logger provided to http.NewServerHandler")
	}
	return &handler{
		requestHandler: requestHandler,
		logger:         logger,
	}
}
