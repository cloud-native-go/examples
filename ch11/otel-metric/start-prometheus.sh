#!/usr/bin/env bash

PROMETHEUS_ENDPOINT="http://localhost:9090"

PROMETHEUS_IMAGE="prom/prometheus:v2.23.0"

set -eux

docker kill prometheus || true

docker rm prometheus || true

docker run -d --name prometheus \
  -p 9090:9090 \
  -v "${PWD}/prometheus.yml:/etc/prometheus/prometheus.yml" \
  "${PROMETHEUS_IMAGE}"

if which open; then
  open "${PROMETHEUS_ENDPOINT}"
else
  echo "Browse to ${PROMETHEUS_ENDPOINT}"
fi
