BUILD := build
GO ?= go
GOFMT ?= go fmt
GOLINT ?= golint

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

.PHONY: test
test:
	$(GO) test -v ./...

.PHONY: build
build:
	$(GO) build -o ./$(BUILD)/apiserver cmd/apiserver/main.go
	$(GO) build -o ./$(BUILD)/scheduler cmd/scheduler/main.go
	$(GO) build -o ./$(BUILD)/sender cmd/sender/main.go
	$(GO) build -o ./$(BUILD)/client cmd/client/main.go
