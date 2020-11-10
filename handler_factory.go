package http

import (
	goHttp "net/http"

	"github.com/containerssh/log"
)

func NewServerHandler(
	requestHandler RequestHandler,
	logger log.Logger,
) goHttp.Handler {
	return &handler{
		requestHandler: requestHandler,
		logger:         logger,
	}
}
