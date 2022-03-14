BIN_CALENDAR := "./bin/calendar"
BIN_CALENDAR_SCHEDULER := "./bin/calendar_scheduler"
BIN_CALENDAR_SENDER := "./bin/calendar_sender"
DOCKER_CALENDAR_IMG="calendar:develop"
DOCKER_CALENDAR_SCHEDULER_IMG="calendar_scheduler:develop"
DOCKER_CALENDAR_SENDER_IMG="calendar_sender:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

generate:
	rm -rf internal/server/grpc/pb
	mkdir -p internal/server/grpc/pb

	protoc \
		--proto_path=api/ \
		--go_out=internal/server/grpc/pb \
		--go-grpc_out=internal/server/grpc/pb \
		api/*.proto

build:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_CALENDAR_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
	go build -v -o $(BIN_CALENDAR_SENDER) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

run: build
	$(BIN_CALENDAR) --config ./configs/calendar_config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_CALENDAR_IMG) \
		-f build/calendar/Dockerfile .
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_CALENDAR_SCHEDULER_IMG) \
		-f build/calendar_scheduler/Dockerfile .
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_CALENDAR_SENDER_IMG) \
		-f build/calendar_sender/Dockerfile .

version: build
	$(BIN_CALENDAR) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

up:
	docker-compose -f deployments/docker-compose.yaml up --remove-orphans -d

down:
	docker-compose -f deployments/docker-compose.yaml down

integration-tests:
	docker-compose -f deployments/docker-compose.yaml up --remove-orphans -d
	go test -v -tags=integration -count=1 ./tests/integration/...
	docker-compose -f deployments/docker-compose.yaml down

.PHONY: build run build-img run-img version test lint
