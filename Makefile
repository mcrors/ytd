# Define the app name and build directory
APP_NAME = ytd
BUILD_DIR = bin
MAIN_PKG = ./cmd/ytd

.PHONY: run build docker fmt lint vet unit-test integration-test e2e-test

run:
	go run $(MAIN_PKG)

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PKG)

docker:
	docker build -t $(APP_NAME):latest .

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run

unit-test:
	go test ./... -v -cover -short=true

integration-test:
	go test ./... -v -cover

e2e-test:
	docker compose -f docker-compose.test.yaml up --build --exit-code-from test; \
	EXIT=$$?; \
	docker compose -f docker-compose.test.yaml down -v; \
	exit $$EXIT
