# TLS Proxy

A simple TLS proxy that supports SNI-based routing to redirect requests to different backend servers. This is a lightweight alternative to nginx for reverse proxy scenarios.

## Features

- SNI-based routing (routes requests based on the hostname)
- TLS/SSL termination with HTTP/2 support
- Header preservation and forwarding (Host, X-Real-IP, X-Forwarded-For, X-Forwarded-Proto)
- Streaming support without buffering
- Configurable timeouts
- JSON-based configuration

## Configuration

The proxy can be configured using a JSON file (default: `config.json`). Here's an example configuration:

```json
{
  "servers": [
    {
      "serverName": "warpbuild.blob.core.windows.net",
      "certFile": "/etc/nginx/ssl/warpbuild.crt",
      "keyFile": "/etc/nginx/ssl/warpbuild.key",
      "targetURL": "http://127.0.0.1:10000",
      "timeoutSeconds": 600
    },
    {
      "serverName": "results-receiver.actions.githubusercontent.com",
      "certFile": "/etc/nginx/ssl/results-receiver.crt",
      "keyFile": "/etc/nginx/ssl/results-receiver.key",
      "targetURL": "http://127.0.0.1:50051",
      "timeoutSeconds": 600
    }
  ],
  "listenAddr": ":443",
  "enableHTTP2": true,
  "tlsMinVersion": "1.2"
}
```

### Configuration Options

- `servers`: Array of server configurations
  - `serverName`: The hostname to match (SNI)
  - `certFile`: Path to the SSL certificate file
  - `keyFile`: Path to the SSL private key file
  - `targetURL`: Backend URL to proxy requests to
  - `timeoutSeconds`: Request timeout in seconds
- `listenAddr`: Address to listen on (default: `:443`)
- `enableHTTP2`: Enable HTTP/2 support (default: `true`)
- `tlsMinVersion`: Minimum TLS version ("1.0", "1.1", "1.2", "1.3", default: "1.2")

## Usage

### Build the proxy

```bash
go build -o tlsproxy main.go
```

### Run with default configuration

```bash
sudo ./tlsproxy
```

### Run with custom configuration file

```bash
sudo ./tlsproxy -config /path/to/config.json
```

Note: The proxy needs to run as root or with appropriate permissions to bind to port 443.

## How it works

1. The proxy listens on port 443 (or configured port) with TLS enabled
2. When a client connects, it uses SNI to determine which backend to route to
3. The appropriate SSL certificate is loaded based on the hostname
4. Requests are proxied to the configured backend with headers preserved
5. Responses are streamed back to the client without buffering

## Differences from nginx

This proxy provides a simpler alternative to nginx for basic reverse proxy scenarios:

- No request buffering (equivalent to `proxy_request_buffering off`)
- Automatic streaming support
- Preserves all required headers
- Supports unlimited request body size
- Configurable timeouts

## Development

To modify or extend the proxy:

1. Edit `main.go` to add new features
2. Update the `Config` struct for new configuration options
3. Rebuild with `go build`

## Testing

To test the proxy locally, you can:

1. Generate self-signed certificates for testing
2. Update your `/etc/hosts` file to point test domains to localhost
3. Run backend services on the configured ports
4. Start the proxy and make HTTPS requests to test routing
