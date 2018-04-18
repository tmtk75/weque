.DEFAULT_GOAL := help

VERSION := $(shell git describe --tags --abbrev=0)
VERSION_LONG := $(shell git describe --tags)
VAR_VERSION := github.com/tmtk75/weque/cmd.Version

LDFLAGS := -ldflags "-X $(VAR_VERSION)=$(VERSION) \
	-X $(VAR_VERSION)Long=$(VERSION_LONG)"

SRCS := $(shell find . -type f -name '*.go')

.PHONY: build
build: weque  ## Build here

.PHONY: run
run:
	SECRET_TOKEN=abc123 go run ./cmd/weque/main.go serve

.PHONY: run-tls
run-tls:
	SECRET_TOKEN=abc123 go run ./cmd/weque/main.go serve --tls.enabled --tls.port=:1443

weque: $(SRCS)
	go build $(LDFLAGS) -o weque ./cmd/weque

gh-req:
	curl -i \
		-H"content-type: application/json" \
		-H"X-Hub-Signature: sha1=c699905923f6a533824e8fb13a0b344d52146e20" \
		-H"X-Github-Event: push" \
		-H"X-Github-Delivery: local-test" \
		localhost:9981/ \
		-d @github/payload.json

gh-req2:
	curl -i \
		-H"content-type: application/x-www-form-urlencoded" \
		-H"X-Hub-Signature: sha1=af9c4634ebadf38f19f14c713f2ab9c0328934ad" \
		-H"X-Github-Event: push" \
		-H"X-Github-Delivery: local-test2" \
		localhost:9981/ \
		-d @github/payload.txt

bb-req:
	curl -i \
		-H "X-Request-UUID: aabbcc" \
		-H "X-Hook-UUID: xxyyzz" \
		localhost:9981/?secret=abc123 \
		-d @bitbucket/payload.json

dr-req:
	curl -i \
		-H "X-weque-secret: abc123" \
		localhost:9981/registry \
		-d @registry/payload.json

.PHONY: install
install:  ## Install in GOPATH
	go install $(LDFLAGS) ./cmd/weque

.PHONY: clean
clean:  ## Clean
	rm -f weque


.PHONY: registry
registry:
	docker run -p 5000:5000 --rm --name registry \
		-v `pwd`/registry/config.yml:/etc/docker/registry/config.yml \
		registry:2

cert.pem key.pem:
	go run $$GOROOT/src/crypto/tls/generate_cert.go \
		--rsa-bits 1024 --host 127.0.0.1,::1,localhost \
		--ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h

.PHONY: tcpflow
tcpflow:
	tcpflow -i lo0 -C 'port 3000'

.PHONY: build-release archive
build-release: build/weque_linux_amd64
archive: build/weque_linux_amd64.gz
release: upload-archives

upload-archives: build/weque_linux_amd64.gz
	echo ghr -u tmtk75 $(VERSION) ./build

build/weque_linux_amd64.gz: build-release
	gzip -f -k build/weque_linux_amd64

build/weque_linux_amd64:
	GOARCH=amd64 GOOS=linux  go build -o build/weque_linux_amd64 ./cmd/weque/main.go

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-17s\033[0m %s\n", $$1, $$2}'
