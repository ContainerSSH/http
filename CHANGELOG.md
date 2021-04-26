# Changelog

## 1.1.0: Adding support for additional HTTP methods

This release adds support for the Delete, Put, and Patch methods.

## 1.0.2: Adding support for www-urlencoded request body

This release adds the ability to switch request bodies to the `www-urlencoded` encoding for the purposes of OAuth2 authentication.

## 1.0.1: Explicitly setting the `Accept` header

This release explicitly adds the "Accept" header on client requests.

## 1.0.0: Initial stable version

This release tags the first stable version for the ContainerSSH 0.4.0 release.

## 0.9.9: Added onReady hook, TLS settings

Added the onReady hook to allow implementing services to directly inject ready handlers. Also added configurable server TLS settings.

## 0.9.8: Message code cleanup

This release cleans up the message codes being emitted.

## 0.9.7: Bugfixing validation

This release fixes a validation bug introduced in the previous version where TLS parameters were validated even if the URL didn't point to a `https://` URL.

## 0.9.6: Configurable TLS support, unified logging

This release adds configurable TLS versions, ciphers, ECDH curves, as well as transitioning to the unified logging interface. 

## 0.9.5: Added config validation

This release adds a `Validate()` method to both the client and the server configuration. This will allow for a central validation of the entire configuration structure.

## 0.9.4: Rolled back context handling

This release removes the previously erroneously added context handling.

## 0.9.3: Context handling, better errors

This release includes two changes:

1. The `Post()` method now accepts a context variable as its first parameter for timeout handling.
2. The `Post()` method now exclusively returns a `http.ClientError`, which includes the reason for failure.

## 0.9.2: URL instead or Url

This release changes the `Url` variable for the client to `URL`. It also bumps the [log dependency](https://github.com/containerssh/log) to the latest release.

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