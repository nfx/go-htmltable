default: vendor

fmt:
	go fmt ./...

vendor:
	go mod vendor

test:
	go test ./... -coverprofile=./vendor/coverage.txt -timeout=30s

coverage: test
	go tool cover -html=./vendor/coverage.txt

.PHONY: build fmt coverage test vendor
