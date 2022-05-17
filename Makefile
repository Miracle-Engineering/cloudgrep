TARGETS = darwin/amd64 linux/amd64 linux/386 windows/amd64 windows/386
VERSION ?= dev
GITHUB_SHA ?= $(shell git rev-parse HEAD)
BUILD_TIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" | tr -d '\n')
GO_VERSION = $(shell go version | awk {'print $$3'})
LDFLAGS = -s -w
PKG = github.com/run-x/cloudgrep

DOCKER_RELEASE_TAG = "ghcr.io/run-x/cloudgrep:$(VERSION)"
DOCKER_LATEST_TAG = "ghcr.io/run-x/cloudgrep:main"

usage:
	@echo ""
	@echo "Task                 : Description"
	@echo "-----------------    : -------------------"
	@echo "make clean           : Remove all build files and reset assets"
	@echo "make build           : Generate build for current OS"
	@echo "make format      	: Format code"
	@echo "make frontend-build  : Build the frontend assets"
	@echo "make frontend-deploy : Deploy the frontend assets"
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

frontend-build:
	docker run -v "$(PWD)/fe":/usr/src/app -w /usr/src/app node:18 npm install
	docker run -v "$(PWD)/fe":/usr/src/app -w /usr/src/app node:18 npm run build

frontend-deploy:
	rm -rf ./static/css ./static/js ./static/*.ico ./static/*.html ./static/*.txt ./static/*.json ./static/*.png
	cp -r ./fe/build/static/css ./static
	cp -r ./fe/build/static/js ./static
	cp ./fe/build/*.ico ./static
	cp ./fe/build/*.html ./static
	cp ./fe/build/*.txt ./static
	cp ./fe/build/*.json ./static
	cp ./fe/build/*.png ./static

build: LDFLAGS += -X $(PKG)/pkg/api.GitCommit=$(GITHUB_SHA)
build: LDFLAGS += -X $(PKG)/pkg/api.BuildTime=$(BUILD_TIME)
build: LDFLAGS += -X $(PKG)/pkg/api.GoVersion=$(GO_VERSION)
build: LDFLAGS += -X $(PKG)/pkg/api.Version=$(VERSION)
build:
	go build -ldflags "$(LDFLAGS)"
	@echo "You can now execute ./cloudgrep"

release: LDFLAGS += -X $(PKG)/pkg/api.GitCommit=$(GITHUB_SHA)
release: LDFLAGS += -X $(PKG)/pkg/api.BuildTime=$(BUILD_TIME)
release: LDFLAGS += -X $(PKG)/pkg/api.GoVersion=$(GO_VERSION)
release: LDFLAGS += -X $(PKG)/pkg/api.Version=$(VERSION)
release:
	@echo "Building binaries..."
	@CGO_ENABLED=1 gox \
		-osarch "$(TARGETS)" \
		-ldflags "$(LDFLAGS)" \
		-output "./bin/cloudgrep_{{.OS}}_{{.Arch}}"

	@echo "Building Linux ARM64 binaries..."
	CC="/usr/bin/aarch64-linux-gnu-gcc" CGO_ENABLED=1 GOOS=linux GOARCH=arm64 GOARM=7 go build -ldflags "$(LDFLAGS)" -o "./bin/cloudgrep_linux_arm64_v7"

	@echo "\nPackaging binaries...\n"
	@./script/package.sh

release-darwin:
	@echo "Building Darwin ARM64 binaries..."
	CGO_LDFLAGS="-L/usr/lib" CGO_ENABLED=1 GOARCH=arm64 GOOS=darwin \
		go build -ldflags "-s -w -linkmode=external"  -o "./bin/cloudgrep_darwin_arm64"

	@echo "Building Darwin AMD64 binaries (require Mac OS)..."
	CGO_LDFLAGS="-L/usr/lib" CGO_ENABLED=1 GOARCH=amd64 GOOS=darwin \
		go build -ldflags "-s -w -linkmode=external"  -o "./bin/cloudgrep_darwin_amd64"
		
	@echo "\nPackaging binaries...\n"
	@./script/package.sh

release-linux-amd64: LDFLAGS += -X $(PKG)/pkg/api.GitCommit=$(GIT_COMMIT)
release-linux-amd64: LDFLAGS += -X $(PKG)/pkg/api.BuildTime=$(BUILD_TIME)
release-linux-amd64: LDFLAGS += -X $(PKG)/pkg/api.GoVersion=$(GO_VERSION)
release-linux-amd64: LDFLAGS += -X $(PKG)/pkg/api.Version=$(VERSION)
release-linux-amd64:
	@echo "Building Linux binaries..."
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o "./bin/cloudgrep_linux_amd64"


docker-build:
	docker build --no-cache -t $(DOCKER_RELEASE_TAG) .
	docker tag $(DOCKER_RELEASE_TAG) $(DOCKER_LATEST_TAG)
	docker images $(DOCKER_RELEASE_TAG)

docker-push:
	docker push $(DOCKER_RELEASE_TAG)
	docker push $(DOCKER_LATEST_TAG)

setup:
	go install github.com/mitchellh/gox@v1.0.1
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.2

clean:
	@rm -f ./cloudgrep
	@rm -rf ./bin/*
