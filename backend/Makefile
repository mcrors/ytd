# Define the app name and build directory
APP_NAME = ytd-api
BUILD_DIR = build

# Default target to run the app
.PHONY: run build docker fmt lint

run:
	go run main.go

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) .

docker:
	docker build -t $(APP_NAME):latest .

fmt:
	go fmt ./...

lint:
	golangci-lint run
