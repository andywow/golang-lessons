BUILD := build
GO ?= go
GOFMT ?= go fmt
GOLINT ?= golangci-lint run

.PHONY: clean
clean:
	$(GO) clean -i ./...
	rm -rf $(BUILD)

.PHONY: fmt
fmt:
	$(GOFMT) ./...

.PHONY: lint
lint:
	$(GOLINT) ./...

.PHONY: unittest
unittest:
	$(GO) test -v ./...

.PHONY: build
build:
	$(GO) build -o ./$(BUILD)/apiserver cmd/apiserver/main.go
	$(GO) build -o ./$(BUILD)/scheduler cmd/scheduler/main.go
	$(GO) build -o ./$(BUILD)/sender cmd/sender/main.go
	$(GO) build -o ./$(BUILD)/client cmd/client/main.go

.PHONY: dockerbuild
dockerbuild:
	docker build -f deploy/apiserver/Dockerfile -t andywow/calendar-apiserver:latest .
	docker build -f deploy/scheduler/Dockerfile -t andywow/calendar-scheduler:latest .
	docker build -f deploy/sender/Dockerfile -t andywow/calendar-sender:latest .

.PHONY: start
start:
	docker-compose up --build

.PHONY: stop
stop:
	docker-compose down

.PHONY: test
test:
	set -e
	docker-compose -f docker-compose.yml -f docker-compose.test.yml up --build -d
	docker-compose -f docker-compose.yml -f docker-compose.test.yml run tests go test
	rc=$$?
	docker-compose -f docker-compose.yml -f docker-compose.test.yml down
	exit $$rc

.PHONY: testdown
testdown:
	docker-compose -f docker-compose.yml -f docker-compose.test.yml down
