buildAt=`date +%Y/%m/%d-%H:%M:%S`
gitHash=`(git log -1 --pretty=format:%H 2>/dev/null || echo "no-commit")`
version=`(git describe --tags --exact-match HEAD 2>/dev/null || git rev-parse --abbrev-ref HEAD 2>/dev/null | grep -v HEAD 2>/dev/null || echo "dev")`
# 修改为项目特定的变量路径
flags="-X '{{.ModuleName}}/pkg/utils.Version=${version}' -X '{{.ModuleName}}/pkg/utils.BuildAt=${buildAt}' -X '{{.ModuleName}}/pkg/utils.GitHash=${gitHash}'"
release_flags="-w -s ${flags}"

GOPATH:=$(shell go env GOPATH)

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: release
release:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags=${flags} -o bin/release/{{.ProjectName}} .
	@cp config.toml bin/release/

.PHONY: build
build:
	@go build -ldflags=${flags} -o bin/{{.ProjectName}} .

.PHONY: run
run: build
	@./bin/{{.ProjectName}}

.PHONY: test
test:
	@go test -v ./tests/... -cover

.PHONY: info
info:
	@echo "Build Information:"
	@echo "=================="
	@echo "Build Time: $(buildAt)"
	@echo "Git Hash: $(gitHash)"
	@echo "Version: $(version)"

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: tools
tools:
	go install github.com/air-verse/air@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go get -u go.ipao.vip/atom
	go get -u google.golang.org/genproto
	go get -u github.com/gofiber/fiber/v3
.PHONY: init
init: tools
	@atomctl swag init
	@atomctl gen enum
	@atomctl gen route
	@atomctl gen service
	@buf generate
	@go mod tidy
	@go get -u
	@go mod tidy