LINUX_TARGETS = linux/amd64 linux/386
VERSION ?= dev
GITHUB_SHA ?= $(shell git rev-parse HEAD)
BUILD_TIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" | tr -d '\n')
GO_VERSION = $(shell go version | awk {'print $$3'})
LDFLAGS = -s -w
PKG = github.com/run-x/cloudgrep

DOCKER_RELEASE_TAG = "ghcr.io/run-x/cloudgrep:$(VERSION)"
DOCKER_LATEST_TAG = "ghcr.io/run-x/cloudgrep:main"

THIS_FILE := $(lastword $(MAKEFILE_LIST))

usage:
	@echo ""
	@echo "Task                 : Description"
	@echo "-----------------    : -------------------"
	@echo "make clean           : Remove all build files and reset assets"
	@echo "make build           : Generate build for current OS"
	@echo "make format      	: Format code"
	@echo "make frontend-build  : Build the frontend assets"
	@echo "make frontend-deploy : Deploy the frontend assets"
	@echo "make release-linux   : Generate binaries for Linux"
	@echo "make release-windows : Generate binaries for Windows"
	@echo "make release-darwin  : Generate binaries for macOS"
	@echo "make run           	: Run using local code"
	@echo "make run-demo       	: Run the demo"
	@echo "make setup           : Install all necessary dependencies"
	@echo "make test            : Execute test suite"
	@echo "make load-test       : Execute load test suite"
	@echo "make version         : Show version"
	@echo ""

format:
	go fmt github.com/run-x/...

lint:
	golangci-lint run ./...

test:
	go test -race ./hack/... ./pkg/... ./cmd/... -coverprofile=coverage.out -covermode=atomic

load-test:
	go clean -testcache && go test ./loadtest/...

pre-commit:
	@$(MAKE) -f $(THIS_FILE) format
	@$(MAKE) -f $(THIS_FILE) lint
	@$(MAKE) -f $(THIS_FILE) test

run:
	go run -race main.go

run-demo:
	go run -race main.go  --config demo/demo.yaml

version:
	@go run -race main.go --version

frontend-build:
	docker run -v "$(PWD)/fe":/usr/src/app -w /usr/src/app node:18 npm install
	docker run -v "$(PWD)/fe":/usr/src/app -w /usr/src/app node:18 npm run build

frontend-deploy:
	rm -rf ./static/css ./static/js ./static/*.ico ./static/*.html ./static/*.txt ./static/*.json ./static/*.png
	cp -r ./fe/build/static/css ./static
	cp -r ./fe/build/static/js ./static
	cp ./fe/build/*.ico ./static
	cp ./fe/build/*.html ./static
	cp ./fe/build/*.js ./static/js
	cp ./fe/build/*.txt ./static
	cp ./fe/build/*.json ./static
	cp ./fe/build/*.png ./static

build: LDFLAGS += -X $(PKG)/pkg/version.GitCommit=$(GITHUB_SHA)
build: LDFLAGS += -X $(PKG)/pkg/version.BuildTime=$(BUILD_TIME)
build: LDFLAGS += -X $(PKG)/pkg/version.GoVersion=$(GO_VERSION)
build: LDFLAGS += -X $(PKG)/pkg/version.Version=$(VERSION)
build:
	go build -race -ldflags "$(LDFLAGS)"
	@echo "You can now execute ./cloudgrep"

release-windows: LDFLAGS += -X $(PKG)/pkg/version.GitCommit=$(GITHUB_SHA)
release-windows: LDFLAGS += -X $(PKG)/pkg/version.BuildTime=$(BUILD_TIME)
release-windows: LDFLAGS += -X $(PKG)/pkg/version.GoVersion=$(GO_VERSION)
release-windows: LDFLAGS += -X $(PKG)/pkg/version.Version=$(VERSION)
release-windows:
	@echo "Building binary for 386."
	@CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ CC_FOR_TARGET=i686-w64-mingw32-gcc GOOS=windows GOARCH=386 go build \
		-ldflags "$(LDFLAGS) -extld=i686-w64-mingw32-gcc" \
		-o "./bin/cloudgrep_windows_386"

	@echo "Building binary for amd64."
	@CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ CC_FOR_TARGET=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build \
		-ldflags "$(LDFLAGS) -extld=x86_64-w64-mingw32-gcc" \
		-o "./bin/cloudgrep_windows_amd64"

	@echo "\nPackaging binaries...\n"
	@./script/package.sh

release-linux: LDFLAGS += -X $(PKG)/pkg/version.GitCommit=$(GITHUB_SHA)
release-linux: LDFLAGS += -X $(PKG)/pkg/version.BuildTime=$(BUILD_TIME)
release-linux: LDFLAGS += -X $(PKG)/pkg/version.GoVersion=$(GO_VERSION)
release-linux: LDFLAGS += -X $(PKG)/pkg/version.Version=$(VERSION)
release-linux:
	@echo "Building binaries..."
	@CGO_ENABLED=1 gox \
		-osarch "$(LINUX_TARGETS)" \
		-ldflags "$(LDFLAGS)" \
		-output "./bin/cloudgrep_{{.OS}}_{{.Arch}}"

	@echo "Building Linux ARM64 binaries..."
	CC="/usr/bin/aarch64-linux-gnu-gcc" CGO_ENABLED=1 GOOS=linux GOARCH=arm64 GOARM=7 go build -ldflags "$(LDFLAGS)" -o "./bin/cloudgrep_linux_arm64_v7"

	@echo "\nPackaging binaries...\n"
	@./script/package.sh

release-darwin: LDFLAGS += -X $(PKG)/pkg/version.GitCommit=$(GITHUB_SHA)
release-darwin: LDFLAGS += -X $(PKG)/pkg/version.BuildTime=$(BUILD_TIME)
release-darwin: LDFLAGS += -X $(PKG)/pkg/version.GoVersion=$(GO_VERSION)
release-darwin: LDFLAGS += -X $(PKG)/pkg/version.Version=$(VERSION)
release-darwin: LDFLAGS += "-linkmode=external"
release-darwin:
	@echo "Building Darwin ARM64 binary"
	CGO_LDFLAGS="-L/usr/lib" CGO_ENABLED=1 GOARCH=arm64 GOOS=darwin \
		go build -ldflags "$(LDFLAGS)" -o "./bin/cloudgrep_darwin_arm64"

	@echo "Building Darwin AMD64 binary..."
	CGO_LDFLAGS="-L/usr/lib" CGO_ENABLED=1 GOARCH=amd64 GOOS=darwin \
		go build -ldflags "$(LDFLAGS)"  -o "./bin/cloudgrep_darwin_amd64"

	@echo "Signing Darwin ARM64 binary..."
	gon gon_arm64.hcl

	@echo "Signing Darwin AMD64 binary..."
	gon gon_amd64.hcl

docker-build:
	docker build --no-cache -t $(DOCKER_RELEASE_TAG) .
	docker tag $(DOCKER_RELEASE_TAG) $(DOCKER_LATEST_TAG)
	docker images $(DOCKER_RELEASE_TAG)

docker-push:
	docker push $(DOCKER_RELEASE_TAG)
	docker push $(DOCKER_LATEST_TAG)

setup:
	go install github.com/mitchellh/gox@v1.0.1
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2

clean:
	@rm -f ./cloudgrep
	@rm -rf ./bin/*

awsgen:
	CGO_ENABLED=1 go run -race ./hack/awsgen --config pkg/provider/aws/config/config.yaml --output-dir pkg/provider/aws

