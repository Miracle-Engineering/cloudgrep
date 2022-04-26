TARGETS = darwin/amd64 darwin/arm64 linux/amd64 linux/386 windows/amd64 windows/386
GIT_COMMIT = $(shell git rev-parse HEAD)
BUILD_TIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" | tr -d '\n')
GO_VERSION = $(shell go version | awk {'print $$3'})
LDFLAGS = -s -w
PKG = github.com/run-x/cloudgrep

usage:
	@echo ""
	@echo "Task                 : Description"
	@echo "-----------------    : -------------------"
	@echo "make clean           : Remove all build files and reset assets"
	@echo "make build           : Generate build for current OS"
	@echo "make format      	: Format code"
	@echo "make release         : Generate binaries for all supported OSes"
	@echo "make run           	: Run using local code"
	@echo "make setup           : Install all necessary dependencies"
	@echo "make test            : Execute test suite"
	@echo "make version         : Show version"
	@echo ""

format:
	go fmt github.com/run-x/...

lint:
	golangci-lint run ./...

test:
	go test -race -cover ./pkg/...

run:
	go run main.go

version:
	@go run main.go --version

build:
	go build
	@echo "You can now execute ./cloudgrep"

release: LDFLAGS += -X $(PKG)/pkg/command.GitCommit=$(GIT_COMMIT)
release: LDFLAGS += -X $(PKG)/pkg/command.BuildTime=$(BUILD_TIME)
release: LDFLAGS += -X $(PKG)/pkg/command.GoVersion=$(GO_VERSION)
release: LDFLAGS += -X $(PKG)/pkg/command.Version=$(VERSION)
release:
	@echo "Building binaries..."
	@gox \
		-osarch "$(TARGETS)" \
		-ldflags "$(LDFLAGS)" \
		-output "./bin/cloudgrep_{{.OS}}_{{.Arch}}"

	@echo "Building ARM binaries..."
	GOOS=linux GOARCH=arm GOARM=5 go build -ldflags "$(LDFLAGS)" -o "./bin/cloudgrep_linux_arm_v5"

	@echo "Building ARM64 binaries..."
	GOOS=linux GOARCH=arm64 GOARM=7 go build -ldflags "$(LDFLAGS)" -o "./bin/cloudgrep_linux_arm64_v7"

	@echo "\nPackaging binaries...\n"
	@./script/package.sh

setup:
	go install github.com/mitchellh/gox@v1.0.1
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.2

clean:
	@rm -f ./cloudgrep
	@rm -rf ./bin/*
