PHONY: generate-sdk

generate-sdk:
	@echo "Generating SDK..."
	@sh ./scripts/generate-from-openapi.sh ${PWD}/pkg/warpbuild --release-go
	@echo "SDK generated successfully"