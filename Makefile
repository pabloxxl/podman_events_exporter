VERSION=$(shell cat version)
COMMIT=$(shell git rev-parse --short HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
HOST=$(shell hostname)
TIMESTAMP=$(shell date '+%Y-%m-%dT%H:%M:%S')

NAME=podman_events_exporter
BINARY_NAME=${NAME}_${VERSION}
PACKAGE=github.com/pabloxxl/podman_events_exporter

all: build
run: build execute
 
build:
	mkdir -p bin
	CGO_ENABLED=0 go build -tags containers_image_openpgp -o bin/${BINARY_NAME} -ldflags="-X 'main.Version=${VERSION}' -X 'main.BuildCommit=${COMMIT}' \
	 -X 'main.BuildBranch=${BRANCH}' -X 'main.BuildHost=${HOST}' -X 'main.BuildTime=${TIMESTAMP}'" main.go
 
execute:
	bin/${BINARY_NAME}
 
clean:
	go clean
	rm -rf bin
