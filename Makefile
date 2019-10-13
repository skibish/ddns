# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

.PHONY: test build
.DEFAULT_GOAL := help

GIT_COMMIT=$(shell git rev-list -1 HEAD)
GIT_VERSION=$(shell git describe --abbrev=0 --tags)

test: ## run tests
	go test -v -cover -race `go list ./... | grep -v /vendor/`

build: ## build binaries for distribution
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.buildVersion=$(GIT_VERSION) -X main.buildCommitHash=$(GIT_COMMIT)" -o ddns-Linux-x86_64 .
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.buildVersion=$(GIT_VERSION) -X main.buildCommitHash=$(GIT_COMMIT)" -o ddns-Darwin-x86_64 .
	GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-X main.buildVersion=$(GIT_VERSION) -X main.buildCommitHash=$(GIT_COMMIT)" -o ddns-Linux-armv7l .

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
