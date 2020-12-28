#!/usr/bin/env bash

JAEGER_ENDPOINT="http://localhost:16686/search"

JAEGER_IMAGE="jaegertracing/all-in-one:1.21"

set -eux

docker kill jaeger || true

docker rm jaeger || true

docker run -d --name jaeger \
  -p 16686:16686 \
  -p 14268:14268 \
  "${JAEGER_IMAGE}"

if which open; then
  open "${JAEGER_ENDPOINT}"
else
  echo "Browse to ${JAEGER_ENDPOINT}"
fi
