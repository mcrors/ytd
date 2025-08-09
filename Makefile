# Define the app name and build directory
APP_NAME = ytd
BUILD_DIR = bin
MAIN_FILE = cmd/ytd/main.go

# Default target to run the app
.PHONY: run build docker fmt lint

run:
	go run $(MAIN_FILE)

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)

docker:
	docker build -t $(APP_NAME):latest $(MAIN_FILE)

fmt:
	go fmt ./...

lint:
	golangci-lint run
