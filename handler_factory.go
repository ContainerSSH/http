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
	return &handler{
		requestHandler: requestHandler,
		logger:         logger,
	}
}
