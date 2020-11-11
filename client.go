package http

import (
	"time"
)

// Client is a simplified HTTP interface that ensures that a struct is transported to a remote endpoint
// properly encoded, and the response is decoded into the response struct.
type Client interface {
	// Post queries the configured endpoint with the path, sending the requestBody and providing the
	// response in the responseBody structure. It returns the HTTP status code and any potential errors.
	Post(path string, requestBody interface{}, responseBody interface{}) (statusCode int, err error)
}

// ClientConfiguration is the configuration structure for HTTP clients
type ClientConfiguration struct {
	// URL is the base URL for requests.
	Url string `json:"url" yaml:"url" comment:"URL of the authentication server."`
	// CaCerts is either the CA certificate to expect on the server in PEM format
	//         or the name of a file containing the PEM.
	CaCert string `json:"cacert" yaml:"cacert" comment:"CA certificate in PEM format to use for host verification. Note: due to a bug in Go on Windows this has to be explicitly provided."`
	// Timeout is the time the client should wait for a response.
	Timeout time.Duration `json:"timeout" yaml:"timeout" comment:"Timeout in seconds" default:"2s"`
	// ClientCert is a PEM containing an x509 certificate to present to the server or a file name containing the PEM.
	ClientCert string `json:"cert" yaml:"cert" comment:"Client certificate file in PEM format."`
	// ClientKey is a PEM containing a private key to use to connect the server or a file name containing the PEM.
	ClientKey string `json:"key" yaml:"key" comment:"Client key file in PEM format."`
}
