.PHONY: test cover lint fmt vet tidy

test:
	go test ./...

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

lint: vet
	@echo "checking gofmt..."
	@test -z "$$(gofmt -l .)" || (echo "unformatted files:"; gofmt -l .; exit 1)

fmt:
	gofmt -w .

vet:
	go vet ./...

tidy:
	go mod tidy
