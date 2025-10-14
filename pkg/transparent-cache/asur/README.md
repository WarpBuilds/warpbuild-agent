# Azure to R2 Proxy

A proxy server that translates Azure Blob Storage API calls to R2 (Cloudflare) S3-compatible API.

## Setup

### 1. Environment Variables

Create a `.env` file in the `azproxy` directory with your credentials:

```bash
# R2 Credentials
R2_ACCESS_KEY_ID=your_r2_access_key_id
R2_SECRET_ACCESS_KEY=your_r2_secret_access_key
R2_ENDPOINT=https://your_account_id.r2.cloudflarestorage.com

# S3 Credentials (optional, used as fallback)
S3_ACCESS_KEY_ID=your_s3_access_key_id
S3_SECRET_ACCESS_KEY=your_s3_secret_access_key
```

You can also export these as environment variables directly:

```bash
export R2_ACCESS_KEY_ID=your_r2_access_key_id
export R2_SECRET_ACCESS_KEY=your_r2_secret_access_key
export R2_ENDPOINT=https://your_account_id.r2.cloudflarestorage.com
```

### 2. Configuration

The proxy supports several configuration options via environment variables:

- `AZPROXY_UPLOAD_METHOD`: Choose between `http` (default) or `s3` for upload method
- `AZPROXY_UPLOAD_CONCURRENCY`: Number of concurrent upload workers (default: 10)
- `AZPROXY_DEBUG`: Set to `true` for verbose logging
- `AZPROXY_ENABLE_HTTP2`: Set to `true` to enable HTTP/2 (default: disabled)
- `AZPROXY_MAX_CONNS_PER_HOST`: Maximum connections per host (default: 100)
- `AZPROXY_DISABLE_KEEPALIVE`: Set to `true` to disable HTTP keep-alive

### 3. Running the Proxy

```bash
# Build
go build -o azproxy .

# Run with environment variables
./azproxy

# Or run directly with go
go run main.go uploader.go
```

The proxy will listen on port 10000 by default.

### 4. Debug Endpoints

When `AZPROXY_DEBUG=true`, the following endpoints are available:

- `/_debug/health` - Health check and S3 connectivity status
- `/_debug/stats` - Performance statistics and active uploads

## Upload Methods

### HTTP Method (Default)

Uses direct HTTP requests with AWS v4 signing to upload to R2. This method:

- Implements the S3 multipart upload protocol directly
- May provide better performance for R2
- Supports concurrent part uploads

### S3 SDK Method

Uses the AWS SDK for Go v2 to upload. Set `AZPROXY_UPLOAD_METHOD=s3` to use this method.

## Example .env file

```bash
# R2 Credentials
R2_ACCESS_KEY_ID=your_r2_access_key_id_here
R2_SECRET_ACCESS_KEY=your_r2_secret_access_key_here
R2_ENDPOINT=https://your_account_id.r2.cloudflarestorage.com
R2_BUCKET=your_bucket_name  # Optional, used by test scripts
```
