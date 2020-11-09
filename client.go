package http

// Client is a simplified HTTP interface that ensures that a struct is transported to a remote endpoint
// properly encoded, and the response is decoded into the response struct.
type Client interface {
	// Post queries the configured endpoint with the path, sending the requestBody and providing the
	// response in the responseBody structure. It returns the HTTP status code and any potential errors.
	Post(path string, requestBody interface{}, responseBody interface{}) (int, error)
}
