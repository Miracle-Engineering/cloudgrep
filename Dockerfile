##
## Build
##

# TODO once this repo is public we can simply wget the binary from the latest release instead of building it
FROM golang:1.18-buster AS build

ARG VERSION
ARG GIT_COMMIT

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY ./pkg ./pkg
COPY ./static ./static
COPY Makefile ./

RUN make release-linux-amd64
RUN mv ./bin/cloudgrep_linux_amd64 /cloudgrep

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /cloudgrep /cloudgrep
COPY ./static ./static

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/cloudgrep",  "--bind=0.0.0.0", "--listen=8080"]
