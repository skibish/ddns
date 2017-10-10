# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

.PHONY: test build
.DEFAULT_GOAL := help

test: ## run tests
	go test -v -cover -race `go list ./... | grep -v /vendor/`

build: ## build binaries for distribution
	GOOS=linux GOARCH=386 go build -o ddns-Linux-x86_64 .
	GOOS=darwin GOARCH=386 go build -o ddns-Darwin-x86_64 .
	GOOS=linux GOARCH=arm GOARM=7 go build -o ddns-Linux-armv7l .

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
