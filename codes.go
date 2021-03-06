package http

const (
	// EFailureEncodeFailed indicates that JSON encoding the request failed. This is usually a bug.
	EFailureEncodeFailed = "HTTP_CLIENT_ENCODE_FAILED"
	// EFailureConnectionFailed indicates a connection failure on the network level.
	EFailureConnectionFailed = "HTTP_CLIENT_CONNECTION_FAILED"
	// EFailureDecodeFailed indicates that decoding the JSON response has failed. The status code is set for this
	// code.
	EFailureDecodeFailed = "HTTP_CLIENT_DECODE_FAILED"
	// EClientRedirectsDisabled indicates that ContainerSSH is not following a HTTP redirect sent by the server.
	EClientRedirectsDisabled = "HTTP_CLIENT_REDIRECTS_DISABLED"

	// MClientRequest is a message indicating a HTTP request sent from ContainerSSH
	MClientRequest = "HTTP_CLIENT_REQUEST"

	// MClientRedirect indicates that the server has sent a HTTP redirect.
	MClientRedirect = "HTTP_CLIENT_REDIRECT"

	// MClientResponse is a message indicating receiving a HTTP response to a client request
	MClientResponse = "HTTP_CLIENT_RESPONSE"
)
