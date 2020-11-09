package http

import (
	"context"
)

// ServerConfiguration is a structure to configure the simple HTTP server by.
type ServerConfiguration struct {
	// Listen contains the IP and port to listen on.
	Listen string `json:"listen" yaml:"listen" default:"0.0.0.0:8080"`
	// CaCert contains either a file name to a CA certificate, or the certificate itself in PEM format to use as a
	// server CA.
	CaCert string `json:"cacert" yaml:"cacert"`
	// Key contains either a file name to a private key, or the private key itself in PEM format to use as a server key.
	Key string `json:"key" yaml:"key"`
	// Cert contains either a file to a certificate, or the certificate itself in PEM format to use as a server
	// certificate.
	Cert string `json:"cert" yaml:"cert"`
	// ClientCaCert contains either a file or a certificate in PEM format to verify the connecting clients by.
	ClientCaCert string `json:"clientcacert" yaml:"clientcacert"`
}

// Server is an interface that specifies the minimum requirements for the server.
type Server interface {
	Run() error
	Shutdown(shutdownContext context.Context)
}
