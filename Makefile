default: vendor

fmt:
	go fmt ./...

vendor:
	go mod vendor

test:
	go test ./... -coverprofile=coverage.out -timeout=10s

coverage: test
	go tool cover -html=coverage.out

.PHONY: build fmt coverage test vendor
