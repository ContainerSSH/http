# Changelog

## 0.9.1: Service Integration (November 23, 2020)

This release changes the API of the HTTP server to integrate with the [Service library](https://github.com/containerssh/service). The public interface now requires using the `Lifecycle` object to start and stop the server. The `Lifecycle` also allows adding hooks for lifecycle events.

```go
server, err := http.NewServer(
    "service name",
    http.ServerConfiguration{
        //...
    },
    handler,
    logger,
)
if err != nil {
    // Handle configuration error
}
// Lifecycle from the github.com/containerssh/service package
lifecycle := service.NewLifecycle(server)
//Add an event hook
lifecycle.OnRunning(...)
go func() {
    if err := lifecycle.Run(); err != nil {
        // Handle error
    }
}()
// Do something else, then shut down the server.
// You can pass a context for the shutdown deadline.
lifecycle.Shutdown(context.Background())
```

## 0.9.0: Initial Release (November 11, 2020)

This is the initial release of the library.