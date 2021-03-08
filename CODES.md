# Message / error codes

| Code | Explanation |
|------|-------------|
| `HTTP_CLIENT_CONNECTION_FAILED` | This message indicates a connection failure on the network level. |
| `HTTP_CLIENT_DECODE_FAILED` | This message indicates that decoding the JSON response has failed. The status code is set for this code. |
| `HTTP_CLIENT_ENCODE_FAILED` | This message indicates that JSON encoding the request failed. This is usually a bug. |
| `HTTP_CLIENT_REDIRECT` | This message indicates that the server responded with a HTTP redirect. |
| `HTTP_CLIENT_REDIRECTS_DISABLED` | This message indicates that ContainerSSH is not following a HTTP redirect sent by the server. Use the allowRedirects option to allow following HTTP redirects. |
| `HTTP_CLIENT_REQUEST` | This message indicates that a HTTP request is being sent from ContainerSSH |
| `HTTP_CLIENT_RESPONSE` | This message indicates that ContainerSSH received a HTTP response from a server. |

