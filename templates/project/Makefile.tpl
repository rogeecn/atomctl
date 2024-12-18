buildAt=`date +%Y/%m/%d-%H:%M:%S`
gitHash=`git rev-parse HEAD`
version=`git rev-parse --abbrev-ref HEAD | grep -v HEAD || git describe --exact-match HEAD || git rev-parse HEAD`  ## todo: use current release git tag
flags="-X 'atom/utils.Version=${version}' -X 'atom/utils.BuildAt=${buildAt}' -X 'atom/utils.GitHash=${gitHash}'"
release_flags="-w -s ${flags}"

GOPATH:=$(shell go env GOPATH)

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: dist
dist:
	@go build -ldflags=${flags} -o bin/debug/atom
	@cp config.toml bin/debug/

.PHONY: release
release:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/release/app .
	@cp config.toml bin/release/

.PHONY: test
test:
	@go test -v ./... -cover

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: proto
proto:
	@buf generate

.PHONY: fresh
fresh:
	@go run . migrate down
	@go run . migrate up

.PHONY: mup
mup:
	@go run . migrate up

.PHONY: mdown
mdown:
	@go run . migrate down

.PHONY: model
model:
	go test -run ^Test_GenModel
