-include .env
export

test:
	@echo "Running tests"
	@go test -v ./...
