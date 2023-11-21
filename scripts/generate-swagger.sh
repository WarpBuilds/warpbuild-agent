#!/bin/bash

# this script needs swag v1.8.6, newer versions throw errors
# go install github.com/swaggo/swag/cmd/swag@v1.8.6

swag init \
	--markdownFiles docs/md/ \
	--propertyStrategy pascalcase \
	--parseInternal \
	--parseDependency \
	--parseDepth 5 \
	--exclude bin/ \
	--exclude sdk/