.PHONY: test coverage lint vet

COMMIT := $(shell git rev-parse HEAD)
VERSION?=latest
DATE := $(shell date -I)

CONTAINER_RUNTIME?=podman

GOARCH?=amd64
GOOS?=linux

dist:
	mkdir -p dist/
build-image:
	$(CONTAINER_RUNTIME) build --pull \
        --build-arg EXPORTER_VERSION=$(VERSION) \
        --build-arg EXPORTER_COMMIT=$(COMMIT) \
        --build-arg EXPORTER_DATE=$(DATE) \
        -t ghcr.io/netways/icinga2-exporter:latest .
build: dist
	GOARCH=$(GOARCH) GOOS=$(GOOS) CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" -o dist/
lint:
	go fmt $(go list ./... | grep -v /vendor/)
vet:
	go vet $(go list ./... | grep -v /vendor/)
test:
	go test -v -race ./...
coverage:
	go test -v -cover -coverprofile=coverage.out ./... &&\
	go tool cover -html=coverage.out -o coverage.html
container:
	podman build --pull -t icinga2-exporter:latest .
clean:
	rm -f dist/*
