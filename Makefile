.PHONY: build test run

build:
	go build ./cmd/...

test:
	docker compose up -d db
	go test ./...

run-api:
	go run ./cmd/api

run-api-mock:
	go run ./cmd/api -mock

run-sync:
	go run ./cmd/sync
