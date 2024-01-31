build:
	@go build -o bin/policygen internal/cmd/policygen/policygen.go

GOLANGCI_LINT_VERSION ?= v1.52.2
install-linter:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

lint:
	golangci-lint run

test-unit:
	go test -cover -coverprofile=./bin/coverage.out ./...

test-coverage-view: test
	go tool cover -html=./bin/coverage.out
