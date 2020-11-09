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

### Implementing the server

The server library is not provided