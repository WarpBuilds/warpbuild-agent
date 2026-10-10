PHONY: generate-sdk generate-proto

generate-sdk:
	@echo "Generating SDK..."
	@sh ./scripts/generate-from-openapi.sh ${PWD}/pkg/warpbuild --release-go
	@echo "SDK generated successfully"

generate-proto:
	@echo "Generating sandbox protos..."
	@cd spec && buf lint && buf generate
	@go test ./pkg/sandboxspec/ -run TestWireContractUnchanged
	@echo "Protos generated; wire contract unchanged"

build-agentd:
	@echo "Building agentd..."
	go build \
		-o ${PWD}/bin/warpbuild-agentd \
		${PWD}/cmd/agentd/main.go

release:
	@echo "Releasing..."
	@sh ./scripts/release.sh
	@echo "Released successfully"
