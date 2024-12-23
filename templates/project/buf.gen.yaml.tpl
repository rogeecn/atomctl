version: v2
managed:
  enabled: true
  override:
    - file_option: go_package_prefix
      value: {{.ModuleName}}/pkg/proto

plugins:
  - remote: buf.build/protocolbuffers/go
    out: gen/go
    opt: paths=source_relative
  - remote: buf.build/grpc/go
    out: gen/go
    opt: paths=source_relative,require_unimplemented_servers=false
  - remote: buf.build/grpc-ecosystem/gateway
    out: gen/go
    opt: paths=source_relative
  # languages
  - remote: buf.build/protocolbuffers/js
    out: gen/js
  - remote: buf.build/protocolbuffers/php
    out: gen/php
  # generate openapi documentation for api
  - remote: buf.build/grpc-ecosystem/openapiv2
    out: gen/openapiv2
    opt: allow_merge=true,merge_file_name=services

inputs:
  - directory: proto
