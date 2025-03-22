version: v2
inputs:
  - directory: proto
managed:
  enabled: true
  override:
    - file_option: go_package_prefix
      value: {{.ModuleName}}/pkg/proto

plugins:
  - local: protoc-gen-go
    out: pkg/proto
    opt: paths=source_relative
  #- local: protoc-gen-grpc-gateway
  #  out: pkg/proto
  #  opt:
  #    - paths=source_relative
  #    - generate_unbound_methods=true
  - local: protoc-gen-go-grpc
    out: pkg/proto
    opt: paths=source_relative
  # - local: protoc-gen-openapiv2
  #   out: docs/proto
