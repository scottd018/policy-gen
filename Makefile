build:
	@go build -o bin/policy-gen internal/cmd/policygen/policygen.go

GOLANGCI_LINT_VERSION ?= v1.55.2
install-linter:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

lint:
	golangci-lint run

test-commit:
	scripts/commit-check-latest.sh

test-unit:
	go test -cover -coverprofile=./bin/coverage.out ./...

test-functional-aws:
	bin/policy-gen aws \
		--input-path=internal/pkg/aws/testinput \
		--output-path=internal/pkg/aws/testoutput \
		--documentation=internal/pkg/aws/testoutput/README.md \
		--force \
		--debug

test-e2e-start-aws:
	docker run \
		--rm -it \
		-p 4566:4566 \
		-e LS_LOG=trace \
		localstack/localstack

test-e2e-aws:
	aws iam create-policy \
		--policy-name test \
		--policy-document file://internal/pkg/aws/testoutput/test.json \
		--description "this is a test aws policy"

test-coverage-view: test-unit
	go tool cover -html=./bin/coverage.out
