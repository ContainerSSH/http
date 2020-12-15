package http

import (
	"fmt"
	"time"
)

// Client is a simplified HTTP interface that ensures that a struct is transported to a remote endpoint
// properly encoded, and the response is decoded into the response struct.
type Client interface {
	// Post queries the configured endpoint with the path, sending the requestBody and providing the
	// response in the responseBody structure. It returns the HTTP status code and any potential errors.
	//
	// The returned error is always one of ClientError
	Post(
		path string,
		requestBody interface{},
		responseBody interface{},
	) (statusCode int, err error)
}

// ClientConfiguration is the configuration structure for HTTP clients
type ClientConfiguration struct {
	// URL is the base URL for requests.
	URL string `json:"url" yaml:"url" comment:"Base URL of the server to connect."`
	// CaCert is either the CA certificate to expect on the server in PEM format
	//         or the name of a file containing the PEM.
	CaCert string `json:"cacert" yaml:"cacert" comment:"CA certificate in PEM format to use for host verification. Note: due to a bug in Go on Windows this has to be explicitly provided."`
	// Timeout is the time the client should wait for a response.
	Timeout time.Duration `json:"timeout" yaml:"timeout" comment:"HTTP call timeout." default:"2s"`
	// ClientCert is a PEM containing an x509 certificate to present to the server or a file name containing the PEM.
	ClientCert string `json:"cert" yaml:"cert" comment:"Client certificate file in PEM format."`
	// ClientKey is a PEM containing a private key to use to connect the server or a file name containing the PEM.
	ClientKey string `json:"key" yaml:"key" comment:"Client key file in PEM format."`
}

// FailureReason describes the Reason why the request failed.
type FailureReason string

const (
	// FailureReasonEncodeFailed indicates that JSON encoding the request failed. This is usually a bug.
	FailureReasonEncodeFailed FailureReason = "encode_failed"
	// FailureReasonConnectionFailed indicates a connection failure.
	FailureReasonConnectionFailed FailureReason = "connection_failed"
	// FailureReasonDecodeFailed indicates that decoding the JSON response has failed. The status code is set for this
	// code.
	FailureReasonDecodeFailed FailureReason = "decode_failed"
)

// ClientError is the the description of the failure of the client request.
type ClientError struct {
	// Reason is one of FailureReason describing the cause of the failure.
	Reason FailureReason `json:"reason" yaml:"reason"`
	// Cause is the original error that is responsible for the error
	Cause error `json:"cause" yaml:"cause"`
	// Message is the message that can be printed into a log.
	Message string `json:"message" yaml:"message"`
}

// Unwrap returns the original error.
func (c ClientError) Unwrap() error {
	return c.Cause
}

// Error returns the error string.
func (c ClientError) Error() string {
	return c.Message
}

// String returns a printable string
func (c ClientError) String() string {
	return fmt.Sprintf("%s: %s (%v)", c.Reason, c.Message, c.Cause)
}
