#!/usr/bin/env bash

JAEGER_ENDPOINT="http://localhost:16686/search"

set -eux

docker kill jaeger || true

docker rm jaeger || true

docker run -d --name jaeger \
  -p 16686:16686 \
  -p 14268:14268 \
  jaegertracing/all-in-one:1.21

if which open; then
  open "${JAEGER_ENDPOINT}"
else
  echo "Browse to ${JAEGER_ENDPOINT}"
fi
