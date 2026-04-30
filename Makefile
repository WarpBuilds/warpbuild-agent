PHONY: generate-sdk

generate-sdk:
	@echo "Generating SDK..."
	@sh ./scripts/generate-from-openapi.sh ${PWD}/pkg/warpbuild --release-go
	@echo "SDK generated successfully"

build-agentd:
	@echo "Building agentd..."
	go build \
		-o ${PWD}/bin/warpbuild-agentd \
		${PWD}/cmd/agentd/main.go

build-pty-broker:
	@echo "Building pty-broker..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" \
		-o ${PWD}/bin/warpbuild-pty-broker \
		${PWD}/cmd/pty-broker/main.go

release:
	@echo "Releasing..."
	@sh ./scripts/release.sh
	@echo "Released successfully"
