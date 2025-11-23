APP_NAME=pr-assign-service
CMD_PATH=./cmd/pr-assign-service

.PHONY: build run docker-build docker-up docker-down

build:
	go build -o bin/$(APP_NAME) $(CMD_PATH)

run:
	go run $(CMD_PATH)

test:
	go test ./internal/api/handlers_test -v

docker-build:
	docker build -t $(APP_NAME) .

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down -v

lint:
	golangci-lint run ./...