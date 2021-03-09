package http

import (
	"github.com/containerssh/service"
)

// Server is an interface that specifies the minimum requirements for the server.
type Server interface {
	service.Service
}
