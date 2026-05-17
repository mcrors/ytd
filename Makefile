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
	docker compose -f docker-compose.dev.yaml up -d --build
	@echo "Waiting for service..."
	@for i in $$(seq 1 30); do \
		curl -sf http://localhost:8080/healthz > /dev/null && break; \
		sleep 1; \
		[ $$i -eq 30 ] && echo "Service did not start in time" && docker compose -f docker-compose.dev.yaml down && exit 1; \
	done
	@echo "--- /healthz ---"
	@curl -sf http://localhost:8080/healthz
	@echo ""
	@echo "--- /readyz (503 expected in dev: no yt-dlp) ---"
	@curl -s http://localhost:8080/readyz
	@echo ""
	@echo "--- /api/directories ---"
	@curl -sf http://localhost:8080/api/directories
	@echo ""
	docker compose -f docker-compose.dev.yaml down
