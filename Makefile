.PHONY: help unit_test integration_test e2e_test test lint coverage_report

help:
	cat Makefile

unit_test:
	go clean -testcache && go test -v ./internal/...

integration_test:
	go clean -testcache && go test -v .

test: unit_test integration_test

mock_generate:
	mockgen -source=internal/storage/types.go -destination=internal/mocks/storage.go -package=mock Storage Transaction

swag_generate:
	swag init --dir cmd/gophermart,internal

swag_format:
	swag fmt

lint:
	go fmt ./...
	find . -name '*.go' -exec goimports -w {} +
	find . -name '*.go' -exec golines -w {} -m 120 \;
	golangci-lint run ./...

run_accrual_service:
	./cmd/accrual/accrual_linux_amd64

run_gophermart:
	go run ./cmd/gophermart

run_db:
	docker compose -f docker-compose.db.yml up -d

coverage_report:
	go test -coverpkg=./... -count=1 -coverprofile=.coverage.out ./...
	go tool cover -html .coverage.out -o .coverage.html