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

release:
	@echo "Releasing..."
	@sh ./scripts/release.sh
	@echo "Released successfully"
