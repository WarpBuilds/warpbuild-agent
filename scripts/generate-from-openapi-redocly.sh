#!/bin/bash
set -e

remove_trailing_slash() {
  local path="$1"
  echo "${path%/}"
}

# Usage function
usage() {
  echo "Usage: $0 <output_path> [--typescript|--go] [version]"
  echo ""
  echo "Arguments:"
  echo "  output_path   Path where the generated SDK will be saved"
  echo "  --typescript  Generate TypeScript SDK (default)"
  echo "  --go         Generate Go SDK"
  echo "  version      Version number for the SDK (optional)"
  echo ""
  echo "Example:"
  echo "  $0 src/warpbuild --typescript 1.0.0"
  echo "  $0 pkg/warpbuild/src --go 1.0.0"
  exit 1
}

# Check arguments
if [[ -z $1 ]]; then
  usage
fi

api_path=$(remove_trailing_slash "$1")

# Determine SDK type (default to typescript)
SDK_TYPE="typescript"
VERSION=""

if [[ -n $2 ]]; then
  if [[ "$2" == "--typescript" ]]; then
    SDK_TYPE="typescript"
    VERSION="$3"
  elif [[ "$2" == "--go" ]]; then
    SDK_TYPE="go"
    VERSION="$3"
  else
    VERSION="$2"
  fi
fi

# Generate filtered swagger file
echo "Generating filtered swagger file"
# Check if redocly is installed
if ! command -v redocly &> /dev/null
then
    echo "redocly could not be found"
    echo "Visit 'https://redocly.com/docs/cli/installation' to install redocly"
    exit 1
fi

redocly bundle filter -o docs/swagger-filtered.yaml --remove-unused-components

swagger_path='docs/swagger-filtered.yaml'
echo "Checking if swagger file exists"
if [ ! -f "$swagger_path" ]; then
  echo "Swagger file does not exist"
  exit 1
fi

echo "Removing existing sdk at $api_path"
rm -rf "$api_path"

if [[ "$SDK_TYPE" == "typescript" ]]; then
  echo "Generating TypeScript SDK using typescript-fetch"
  
  docker run \
    --rm \
    -v ${PWD}:/local \
    openapitools/openapi-generator-cli:v5.2.1 generate \
    -i /local/$swagger_path \
    -g typescript-fetch \
    -o /local/$api_path \
    --additional-properties=supportsES6=true \
    --additional-properties=paramNaming=snake_case \
    --additional-properties=enumPropertyNaming=snake_case \
    --additional-properties=modelPropertyNaming=snake_case \
    --additional-properties=typescriptThreePlus=true \
    --additional-properties=useSingleRequestParameter=true \
    --additional-properties=withSeparateModelsAndApi=true,modelPackage=model,apiPackage=api,npmName=@argonautdev/warpbuild-js-sdk,npmVersion="${VERSION:-1.0.0}",legacyDiscriminatorBehavior=false,disallowAdditionalPropertiesIfNotPresent=false \
    --enable-post-process-file \
    --skip-validate-spec

  echo "Cleaning up TypeScript SDK files"
  rm -rf "$api_path/schemas"
  rmdir "$api_path/typescript-fetch" 2>/dev/null || true
  rm -f "$api_path/package.json" 2>/dev/null || true
  rm -f "$api_path/tsconfig.json" 2>/dev/null || true
  rm -f "$api_path/git_push.sh" 2>/dev/null || true
  rm -f "$api_path/README.md" 2>/dev/null || true
  rm -f "$api_path/.gitignore" 2>/dev/null || true
  rm -f "$api_path/.npmignore" 2>/dev/null || true
  
  echo "TypeScript SDK generated successfully at $api_path"

elif [[ "$SDK_TYPE" == "go" ]]; then
  echo "Generating Go SDK"
  
  docker run \
    --rm \
    -v ${PWD}:/local \
    openapitools/openapi-generator-cli:v7.0.1 generate \
    -i /local/$swagger_path \
    -g go \
    -o /local/$api_path \
    --additional-properties=disallowAdditionalPropertiesIfNotPresent=false \
    --additional-properties=generateInterfaces=true \
    --additional-properties=packageName=warpbuild \
    --additional-properties=isGoSubmodule=true \
    --additional-properties=packageVersion="${VERSION:-1.0.0}" \
    --skip-validate-spec

  echo "Cleaning up Go SDK files"
  rm -f "$api_path/go.mod" 2>/dev/null || true
  rm -f "$api_path/go.sum" 2>/dev/null || true
  rm -f "$api_path/git_push.sh" 2>/dev/null || true
  rm -f "$api_path/.travis.yml" 2>/dev/null || true
  rm -rf "$api_path/test" 2>/dev/null || true
  
  echo "Go SDK generated successfully at $api_path"
fi

# Clean up filtered swagger file
echo "Cleaning up filtered swagger file"
rm -f "$swagger_path"
echo "Done!"
