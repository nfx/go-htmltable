fmt:
	go fmt ./...

test:
	go test  -coverpkg=./... -coverprofile=coverage.out -timeout=10s  ./...

coverage: test
	go tool cover -html=coverage.out

.PHONY: build fmt coverage test
