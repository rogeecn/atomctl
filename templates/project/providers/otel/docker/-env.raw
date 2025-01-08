# Dependent images
GRAFANA_IMAGE=docker.hub.ipao.vip/grafana/grafana:11.4.0
JAEGERTRACING_IMAGE=docker.hub.ipao.vip/jaegertracing/all-in-one:1.64.0
OPENSEARCH_IMAGE=docker.hub.ipao.vip/opensearchproject/opensearch:2.18.0
COLLECTOR_CONTRIB_IMAGE=docker-ghcr.hub.ipao.vip/open-telemetry/opentelemetry-collector-releases/opentelemetry-collector-contrib:0.116.1
PROMETHEUS_IMAGE=docker-quay.hub.ipao.vip/prometheus/prometheus:v3.0.1

# OpenTelemetry Collector
HOST_FILESYSTEM=/
DOCKER_SOCK=/var/run/docker.sock
OTEL_COLLECTOR_HOST=otel-collector
OTEL_COLLECTOR_PORT_GRPC=4317
OTEL_COLLECTOR_PORT_HTTP=4318
OTEL_COLLECTOR_CONFIG=./otel-collector/otelcol-config.yml
OTEL_COLLECTOR_CONFIG_EXTRAS=./otel-collector/otelcol-config-extras.yml
OTEL_EXPORTER_OTLP_ENDPOINT=http://${OTEL_COLLECTOR_HOST}:${OTEL_COLLECTOR_PORT_GRPC}
PUBLIC_OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=http://localhost:8080/otlp-http/v1/traces

# Grafana
GRAFANA_SERVICE_PORT=3000
GRAFANA_SERVICE_HOST=grafana

# Jaeger
JAEGER_SERVICE_PORT=16686
JAEGER_SERVICE_HOST=jaeger

# Prometheus
PROMETHEUS_SERVICE_PORT=9090
PROMETHEUS_SERVICE_HOST=prometheus
PROMETHEUS_ADDR=${PROMETHEUS_SERVICE_HOST}:${PROMETHEUS_SERVICE_PORT}
