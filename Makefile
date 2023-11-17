PHONY: generate-sdk

generate-sdk:
	@echo "Generating SDK..."
	@sh ./scripts/generate-from-openapi.sh ${PWD}/pkg/warpbuild --release-go
	@echo "SDK generated successfully"

build-daemon:
	@echo "Building daemon..."
	go build \
		-o ${PWD}/bin/warpbuild-daemon \
		${PWD}/cmd/daemon/main.go

release:
	@echo "Releasing..."
	@sh ./scripts/release.sh
	@echo "Released successfully"