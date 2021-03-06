# Message/error codes

| Code | Explanation |
| `HTTP_CLIENT_CONNECTION_FAILED` | Sending a HTTP request failed for reasons outside of ContainerSSH |
| `HTTP_CLIENT_DECODE_FAILED` | ContainerSSH failed to decode the JSON response after sending a request. Check if the HTTP server is misbehaving. |
| `HTTP_CLIENT_ENCODE_FAILED` | ContainerSSH failed to encode the payload for sending a HTTP request. This is a bug, please file an issue. |
| `HTTP_CLIENT_REDIRECTS_DISABLED` | ContainerSSH is refusing to follow a HTTP redirect received because the `allowRedirects` option is disabled. |
| `HTTP_CLIENT_REDIRECT` | ContainerSSH has received a HTTP redirect from the server. |
| `HTTP_CLIENT_RESPONSE` | ContainerSSH has received a response to a HTTP request sent to a server. |
| `HTTP_CLIENT_REQUEST` | ContainerSSH is sending a HTTP request to a server. |
