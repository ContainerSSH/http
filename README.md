[![ContainerSSH - Launch Containers on Demand](https://containerssh.github.io/images/logo-for-embedding.svg)](https://containerssh.github.io/)

<!--suppress HtmlDeprecatedAttribute -->
<h1 align="center">ContainerSSH HTTP Library</h1>

[![Go Report Card](https://goreportcard.com/badge/github.com/containerssh/http?style=for-the-badge)](https://goreportcard.com/report/github.com/containerssh/http)
[![LGTM Alerts](https://img.shields.io/lgtm/alerts/github/ContainerSSH/http?style=for-the-badge)](https://lgtm.com/projects/g/ContainerSSH/http/)

This library provides a common layer for HTTP clients and servers in use by ContainerSSH.

<p align="center"><strong>Note: This is a developer documentation.</strong><br />The user documentation for ContainerSSH is located at <a href="https://containerssh.github.io">containerssh.github.io</a>.</p>

## Using this library

This library provides a much simplified API for both the HTTP client and server.

### Using the client

The client library takes a request object that [can be marshalled into JSON format](https://gobyexample.com/json) and sends it to the server. It then fills a response object with the response received from the server. In code:

```go
// Logger is from the github.com/containerssh/log package
logger := standard.New()
clientConfig := http.ClientConfiguration{
    Url:        "http://127.0.0.1:8080/",
    Timeout:    2 * time.Second,
    // You can add TLS configuration here
}
client, err := http.NewClient(clientConfig, logger)
if err != nil {
    // Handle validation error
}

request := yourRequestStruct{}
response := yourResponseStruct{}
responseStatus := uint16(0)

if err := client.Post("/relative/path/from/base/url", &request, &responseStatus, &response); err != nil {
    // Handle connection error
}

if responseStatus > 399 {
    // Handle error
}
```

The `logger` parameter is a logger from the [github.com/containerssh/log](https://github.com/containerssh/log) package.

### Using the server

The server consist of two parts: the HTTP server and the handler. The HTTP server can be used as follows:

```go
server, err := http.NewServer(
    http.ServerConfiguration{
        Listen:       "127.0.0.1:8080",
        // You can also add TLS configuration and certificates here
    },
    handler,
    logger,
)
if err != nil {
    // Handle configuration error
}
go func() {
    if err := server.Run(); err != nil {
        // Handle error
    }
}()
// Do something else, then shut down the server.
// You can pass a context for the shutdown deadline.
server.Shutdown(context.Background())
```

Like before, the `logger` parameter is a logger from the [github.com/containerssh/log](https://github.com/containerssh/log) package. The `handler` is a regular [go HTTP handler](https://golang.org/pkg/net/http/#Handler) that satisfies this interface:

```go
type Handler interface {
    ServeHTTP(http.ResponseWriter, *http.Request)
}
```

## Using a simplified handler

This package also provides a simplified handler that helps with encoding and decoding JSON messages. It can be created as follows:

```go
handler := http.NewServerHandler(yourController, logger)
```

The `yourController` variable then only needs to implement the following interface:

```go
type RequestHandler interface {
	OnRequest(request ServerRequest, response ServerResponse) error
}
```

For example:

```go
type MyRequest struct {
    Message string `json:"message"`
}

type MyResponse struct {
    Message string `json:"message"`
}

type myController struct {
}

func (c *myController) OnRequest(request http.ServerRequest, response http.ServerResponse) error {
    req := MyRequest{}
	if err := request.Decode(&req); err != nil {
		return err
	}
	if req.Message == "Hi" {
		response.SetBody(&MyResponse{
			Message: "Hello world!",
		})
	} else {
        response.SetStatus(400)
		response.SetBody(&MyResponse{
			Message: "Be nice and greet me!",
		})
	}
	return nil
}
```

In other words, the `ServerRequest` object gives you the ability to decode the request into a struct of your choice. The `ServerResponse`, conversely, encodes a struct into the the response body and provides the ability to enter a status code.

## Using multiple handlers

This is a very simple handler example. You can use utility like [gorilla/mux](https://github.com/gorilla/mux) as an intermediate handler between the simplified handler and the server itself.
