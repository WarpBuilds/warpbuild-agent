# warpbuild-agent

Application for runner lifecycle management

## Generating the API client from the OpenAPI specification in backend-core

Use Redocly to generate the API client.

```bash
./scripts/generate-from-openapi-redocly.sh pkg/warpbuild --go 1.0.0
```

The filters in redocly.yaml are used to generate the API client.
