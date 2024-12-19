buildAt=`date +%Y/%m/%d-%H:%M:%S`
gitHash=`git rev-parse HEAD`
version=`git rev-parse --abbrev-ref HEAD | grep -v HEAD || git describe --exact-match HEAD || git rev-parse HEAD`  ## todo: use current release git tag
flags="-X 'atom/utils.Version=${version}' -X 'atom/utils.BuildAt=${buildAt}' -X 'atom/utils.GitHash=${gitHash}'"
release_flags="-w -s ${flags}"

GOPATH:=$(shell go env GOPATH)

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: release
release:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -ldflags=${flags} -o bin/release/{{.ProjectName}} .
	@cp config.toml bin/release/

.PHONY: test
test:
	@go test -v ./... -cover

.PHONY: lint
lint:
	@golangci-lint run
